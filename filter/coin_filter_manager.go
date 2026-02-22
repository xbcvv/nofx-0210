package filter

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"nofx/market"
)

// FilterConfig holds hard filters for candidate coins
type FilterConfig struct {
	MinListingDays    int
	MinQuoteVolume24h float64 // in USDT
	MinOpenInterest   float64 // in USDT
}

var DefaultFilterConfig = FilterConfig{
	MinListingDays:    90,
	MinQuoteVolume24h: 50_000_000,
	MinOpenInterest:   10_000_000,
}

// CoinFilterManager manages the slow-loop candidate coin filtering
type CoinFilterManager struct {
	config    FilterConfig
	apiClient *market.APIClient

	tickerCache map[string]market.Ticker24hr
	oiCache     map[string]float64
	// listingCache stores the listing timestamp (Unix ms) for symbols
	listingCache map[string]int64
	cacheMux     sync.RWMutex

	lastUpdateTime time.Time
}

var GlobalCoinFilter *CoinFilterManager

// InitGlobalCoinFilter initializes the global singleton for coin filtering
func InitGlobalCoinFilter() {
	if GlobalCoinFilter != nil {
		return
	}

	GlobalCoinFilter = &CoinFilterManager{
		config:       DefaultFilterConfig,
		apiClient:    market.NewAPIClientWithTimeout(30 * time.Second),
		tickerCache:  make(map[string]market.Ticker24hr),
		oiCache:      make(map[string]float64),
		listingCache: make(map[string]int64),
	}

	// Load previously saved listing dates to avoid redundant API calls
	GlobalCoinFilter.loadListingCache()

	// Start the daemon for background updates
	go GlobalCoinFilter.daemonLoop()
}

// GetCleanCoins takes the raw list of candidate coins (e.g. from AI500)
// and returns the top N filtered pristine coins based on the strict rules.
func (m *CoinFilterManager) GetCleanCoins(rawSymbols []string, limit int) []string {
	m.cacheMux.RLock()
	defer m.cacheMux.RUnlock()

	var survived []string

	// Time threshold
	now := time.Now()
	// Calculate the time threshold for 90 days ago in ms (approx 90 * 24h)
	ninetyDaysMs := int64(m.config.MinListingDays * 24 * 60 * 60 * 1000)

	for _, symbol := range rawSymbols {
		symbol = strings.ToUpper(symbol)

		// 1. Liquidity (24H Quote Volume) Filter
		ticker, ok := m.tickerCache[symbol]
		if !ok {
			// If not in cache (meaning it wasn't fetched in the last background cycle), skip for safety
			continue
		}
		qVol, _ := strconv.ParseFloat(ticker.QuoteVolume, 64)
		if qVol < m.config.MinQuoteVolume24h {
			// Log occasionally or skip silently (to avoid spam)
			continue
		}

		// 2. Listing Time Filter
		listingTime, ok := m.listingCache[symbol]
		if !ok {
			// if we don't know the listing time yet, we skip or assume it's unsafe
			// The daemon should have fetched it.
			continue
		}
		if now.UnixNano()/1e6-listingTime < ninetyDaysMs {
			continue // Less than 90 days
		}

		// 3. Open Interest (OI) Filter
		oiUSDT := m.oiCache[symbol]
		if oiUSDT < m.config.MinOpenInterest {
			continue // Less than 10M OI
		}

		// 4. Funding rate filter could be added here (e.g. from a funding cache)

		survived = append(survived, symbol)
	}

	// Since rawSymbols is typically pre-sorted by Score/Heat (AI500 list preserves API order)
	// 'survived' retains that relative order. We just truncate to limit.
	if limit > 0 && len(survived) > limit {
		survived = survived[:limit]
	}

	return survived
}

// GetTopVolumeCoins returns the top N coins by 24h quote volume from the local cache
func (m *CoinFilterManager) GetTopVolumeCoins(limit int) []string {
	m.cacheMux.RLock()
	defer m.cacheMux.RUnlock()

	type coinVol struct {
		symbol string
		vol    float64
	}

	var vols []coinVol
	for symbol, ticker := range m.tickerCache {
		// Only consider USDT pairs
		if !strings.HasSuffix(symbol, "USDT") {
			continue
		}
		vol, _ := strconv.ParseFloat(ticker.QuoteVolume, 64)
		vols = append(vols, coinVol{symbol: symbol, vol: vol})
	}

	// Sort descending by volume
	sort.Slice(vols, func(i, j int) bool {
		return vols[i].vol > vols[j].vol
	})

	var result []string
	for i, cv := range vols {
		if limit > 0 && i >= limit {
			break
		}
		result = append(result, cv.symbol)
	}

	return result
}

// daemonLoop updates the cache every 30 minutes
func (m *CoinFilterManager) daemonLoop() {
	// First immediate run
	m.runUpdateCycle()

	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.runUpdateCycle()
	}
}

func (m *CoinFilterManager) runUpdateCycle() {
	log.Println("[CoinFilter] Starting background update cycle (slow loop)...")

	// 1. Fetch 24H Tickers for ALL symbols globally
	tickers, err := m.apiClient.GetAll24hrTickers()
	if err != nil {
		log.Printf("[CoinFilter] Failed to fetch 24hr tickers: %v", err)
		return
	}

	// 2. Fetch Open Interest for pairs that passed volume criteria to save API calls
	newTickerCache := make(map[string]market.Ticker24hr)
	newOICache := make(map[string]float64)
	newListingRequests := []string{}

	m.cacheMux.RLock()
	knownListings := len(m.listingCache) > 0
	m.cacheMux.RUnlock()

	for _, t := range tickers {
		newTickerCache[t.Symbol] = t

		// Only fetch OI heavily if it passes the Volume filter check
		qVol, _ := strconv.ParseFloat(t.QuoteVolume, 64)
		if qVol >= m.config.MinQuoteVolume24h {
			// Check if we need its listing date
			m.cacheMux.RLock()
			_, hasListing := m.listingCache[t.Symbol]
			m.cacheMux.RUnlock()

			if !hasListing && strings.HasSuffix(t.Symbol, "USDT") {
				newListingRequests = append(newListingRequests, t.Symbol)
			}

			// Fetch OI. Since we throttle inside fetchOI, doing it synchronously here inside the slow loop is safe.
			// Typically, few dozen coins pass 50M volume.
			lastPrice, _ := strconv.ParseFloat(t.LastPrice, 64)
			if lastPrice == 0 && t.PriceChange != "" {
				// Parse roughly
			}
			oiValue := m.fetchOIValue(t.Symbol, lastPrice)
			newOICache[t.Symbol] = oiValue
		}
	}

	// 3. Fetch missing listing dates
	if len(newListingRequests) > 0 {
		log.Printf("[CoinFilter] Fetching listing dates for %d new valid coins...", len(newListingRequests))
		m.fetchListingDates(newListingRequests)
	}

	// Update Caches safely
	m.cacheMux.Lock()
	m.tickerCache = newTickerCache
	m.oiCache = newOICache
	// listingCache is updated inside fetchListingDates and loaded from file initially
	if !knownListings && len(m.listingCache) > 0 {
		m.saveListingCache()
	}
	m.lastUpdateTime = time.Now()
	m.cacheMux.Unlock()

	// If we got new listings, save them
	if len(newListingRequests) > 0 {
		m.saveListingCache()
	}

	log.Printf("[CoinFilter] Cycle Complete. Tickers: %d, OI: %d, Listing: %d", len(newTickerCache), len(newOICache), len(m.listingCache))
}

// fetchOIValue retrieves OI and calculates USD value
func (m *CoinFilterManager) fetchOIValue(symbol string, lastApproxPrice float64) float64 {
	// Be polite to Binance API (Wait 50ms)
	time.Sleep(50 * time.Millisecond)

	url := fmt.Sprintf("https://fapi.binance.com/fapi/v1/openInterest?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0
	}

	var result struct {
		Symbol       string `json:"symbol"`
		OpenInterest string `json:"openInterest"` // As a string representing coin quantity
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0
	}

	oiCoins, _ := strconv.ParseFloat(result.OpenInterest, 64)

	// Since OI endpoint returns contracts/coins, calculate USDT value.
	// If the price is extremely low/0, we fetch the real price.
	if lastApproxPrice == 0 {
		lastApproxPrice, _ = m.apiClient.GetCurrentPrice(symbol)
	}

	return oiCoins * lastApproxPrice
}

// fetchListingDates requests the very first 1D kline. openTime == listing date
func (m *CoinFilterManager) fetchListingDates(symbols []string) {
	for _, symbol := range symbols {
		time.Sleep(100 * time.Millisecond) // Rate limiting
		
		// limit=1 & startTime=0 to get the absolute earliest kline
		url := fmt.Sprintf("https://fapi.binance.com/fapi/v1/klines?symbol=%s&interval=1d&limit=1&startTime=0", symbol)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("[CoinFilter] Error fetching listing time for %s: %v", symbol, err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		var klines []interface{}
		if err := json.Unmarshal(body, &klines); err != nil || len(klines) == 0 {
			continue
		}

		firstKline, ok := klines[0].([]interface{})
		if !ok || len(firstKline) == 0 {
			continue
		}

		openTimeFloat, ok := firstKline[0].(float64)
		if !ok {
			continue
		}

		openTimeMs := int64(openTimeFloat)
		m.cacheMux.Lock()
		m.listingCache[symbol] = openTimeMs
		m.cacheMux.Unlock()
	}
}

// saveListingCache persists the listing timeline to disk (we only need to fetch this once in a coin's lifetime!)
func (m *CoinFilterManager) saveListingCache() {
	m.cacheMux.RLock()
	defer m.cacheMux.RUnlock()

	data, err := json.MarshalIndent(m.listingCache, "", "  ")
	if err != nil {
		return
	}

	path := filepath.Join("data", "listing_cache.json")
	os.MkdirAll("data", 0755)
	os.WriteFile(path, data, 0644)
	log.Printf("[CoinFilter] Successfully synced %d listing dates to disk.", len(m.listingCache))
}

// loadListingCache loads the listing dates from disk at boot
func (m *CoinFilterManager) loadListingCache() {
	path := filepath.Join("data", "listing_cache.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var cache map[string]int64
	if err := json.Unmarshal(data, &cache); err == nil {
		m.cacheMux.Lock()
		m.listingCache = cache
		m.cacheMux.Unlock()
		log.Printf("[CoinFilter] Loaded %d listing dates from disk.", len(cache))
	}
}
