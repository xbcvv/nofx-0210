package telegram

import (
	"fmt"
	"net/http"
	"net/url"
	"nofx/config"
	"nofx/logger"
	"nofx/manager"
	"nofx/store"
	"nofx/trader"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	manager         *manager.TraderManager
	store           *store.Store
	bot             *tgbotapi.BotAPI
	adminID         int64
	currentTraderID string
	stateMu         sync.Mutex
	decisionMu      sync.Mutex
	decisionAfter   map[string]int64
}

// NewBot creates a Telegram bot instance if TELEGRAM_BOT_TOKEN is configured.
func NewBot(manager *manager.TraderManager, st *store.Store) (*Bot, error) {
	cfg := config.Get()
	token := cfg.TelegramBotToken
	if token == "" {
		logger.Info("📡 Telegram bot disabled (TELEGRAM_BOT_TOKEN not set)")
		return nil, nil
	}

	userIDStr := cfg.TelegramUserID
	if userIDStr == "" {
		return nil, fmt.Errorf("TELEGRAM_USER_ID not set")
	}
	adminID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid TELEGRAM_USER_ID: %w", err)
	}

	// Create custom HTTP client to support proxy if configured
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Check for proxy settings
	httpProxy := os.Getenv("HTTP_PROXY")
	httpsProxy := os.Getenv("HTTPS_PROXY")
	
	if httpProxy != "" || httpsProxy != "" {
		logger.Infof("📡 Telegram using proxy: HTTP=%s HTTPS=%s", httpProxy, httpsProxy)
		// Go's default transport automatically uses HTTP_PROXY/HTTPS_PROXY environment variables.
		// We log it just to be sure.
	} else {
		// If user wants to configure proxy via .env specific key like TELEGRAM_PROXY
		if proxyURLStr := os.Getenv("TELEGRAM_PROXY"); proxyURLStr != "" {
			proxyURL, err := url.Parse(proxyURLStr)
			if err != nil {
				logger.Warnf("⚠️ Invalid TELEGRAM_PROXY: %v", err)
			} else {
				logger.Infof("📡 Telegram using custom proxy: %s", proxyURLStr)
				httpClient.Transport = &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				}
			}
		}
	}

	api, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w (check your network/proxy)", err)
	}
	
	// Delete any existing webhook to ensure long polling works
	if _, err := api.Request(tgbotapi.DeleteWebhookConfig{DropPendingUpdates: true}); err != nil {
		logger.Warnf("⚠️ Telegram deleteWebhook failed: %v. This might be due to network issues.", err)
	}

	b := &Bot{
		manager:       manager,
		store:         st,
		bot:           api,
		adminID:       adminID,
		decisionAfter: make(map[string]int64),
	}
	b.initDecisionCursors()
	
	// Send startup message
	b.reply(adminID, "🚀 NOFX 系统已启动\nAI 交易系统正在运行中...")
	
	logger.Infof("📡 Telegram bot enabled: @%s (Admin ID: %d)", api.Self.UserName, adminID)
	return b, nil
}

// Start begins polling updates and decision notifications.
func (b *Bot) Start() {
	if b == nil || b.bot == nil {
		return
	}

	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60
	updates := b.bot.GetUpdatesChan(update)

	// Start monitoring for new AI decisions
	go b.watchDecisions()

	for upd := range updates {
		if upd.Message == nil {
			continue
		}
		go b.handleMessage(upd.Message)
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	if msg == nil {
		return
	}
	chatID := msg.Chat.ID
	userID := msg.From.ID

	// Log every message for debugging purposes (temporary, or debug level)
	logger.Infof("📩 Telegram msg received: ChatID=%d UserID=%d Text='%s'", chatID, userID, msg.Text)

	// Strict whitelist check - only allow admin
	if chatID != b.adminID {
		logger.Warnf("⛔ Ignored message from unauthorized ChatID: %d (Expected AdminID: %d)", chatID, b.adminID)
		return
	}

	command := strings.ToLower(strings.TrimSpace(msg.Command()))
	args := strings.Fields(strings.TrimSpace(msg.CommandArguments()))

	if command == "" {
		// Handle generic text messages or "ping"
		if msg.Text == "ping" {
			b.reply(chatID, "pong 🏓")
		}
		return
	}
	logger.Infof("📡 Telegram command received: cmd=%s args=%v", command, args)

	switch command {
	case "start":
		b.reply(chatID, "欢迎使用 NOFX 机器人。\n您已通过身份验证。\n输入 /help 查看指令列表。")
		return
	case "help":
		b.reply(chatID, b.helpText())
		return
	}

	switch command {
	case "status":
		b.handleStatus(chatID, args)
	case "balance":
		b.handleBalance(chatID, args)
	case "positions":
		b.handlePositions(chatID, args)
	case "orders":
		b.handleOrders(chatID, args)
	case "deferred":
		b.handleDeferred(chatID, args)
	case "decision":
		b.handleDecision(chatID, args)
	case "alerts":
		b.handleAlerts(chatID, args)
	case "traders":
		b.handleTraders(chatID)
	default:
		b.reply(chatID, "未知命令，请发送 /help 查看指令列表。")
	}
}

func (b *Bot) handleStatus(chatID int64, args []string) {
	at, errMsg := b.resolveTrader(args)
	if errMsg != "" {
		b.reply(chatID, errMsg)
		return
	}

	status := at.GetStatus()
	lines := []string{
		fmt.Sprintf("🤖 交易员：%s", at.GetName()),
		fmt.Sprintf("运行状态：%v", status["is_running"]),
		fmt.Sprintf("交易所：%v", status["exchange"]),
		fmt.Sprintf("AI 模型：%v", status["ai_model"]),
		fmt.Sprintf("运行时长：%v 分钟", status["runtime_minutes"]),
		fmt.Sprintf("扫描周期：%v", status["scan_interval"]),
	}

	if ws, ok := status["user_data_ws"].(string); ok {
		lines = append(lines, fmt.Sprintf("EXCHANGE_WS：%s", strings.ToUpper(ws)))
	}
	
	if ws, ok := status["mark_price_ws"].(string); ok {
		lines = append(lines, fmt.Sprintf("MARK_PRICE_WS：%s", strings.ToUpper(ws)))
	}

	b.reply(chatID, strings.Join(lines, "\n"))
}

func (b *Bot) handleBalance(chatID int64, args []string) {
	at, errMsg := b.resolveTrader(args)
	if errMsg != "" {
		b.reply(chatID, errMsg)
		return
	}

	// Logic similar to previous implementation
	account, err := at.GetAccountInfo()
	if err != nil {
		b.reply(chatID, fmt.Sprintf("❌ 获取账户失败：%v", err))
		return
	}

	lines := []string{
		fmt.Sprintf("💰 交易员：%s", at.GetName()),
		fmt.Sprintf("权益：%s", formatFloat(account["total_equity"], 2)),
		fmt.Sprintf("可用余额：%s", formatFloat(account["available_balance"], 2)),
		fmt.Sprintf("保证金占用：%s (%s%%)", formatFloat(account["margin_used"], 2), formatFloat(account["margin_used_pct"], 2)),
		fmt.Sprintf("未实现盈亏：%s", formatFloat(account["unrealized_profit"], 2)),
		fmt.Sprintf("总盈亏：%s (%s%%)", formatFloat(account["total_pnl"], 2), formatFloat(account["total_pnl_pct"], 2)),
	}

	b.reply(chatID, strings.Join(lines, "\n"))
}

func (b *Bot) handlePositions(chatID int64, args []string) {
	at, errMsg := b.resolveTrader(args)
	if errMsg != "" {
		b.reply(chatID, errMsg)
		return
	}

	positions, err := at.GetPositions()
	if err != nil {
		b.reply(chatID, fmt.Sprintf("❌ 获取持仓失败：%v", err))
		return
	}

	if len(positions) == 0 {
		b.reply(chatID, "📉 当前无持仓。")
		return
	}

	lines := []string{fmt.Sprintf("📊 交易员：%s 当前持仓", at.GetName())}
	for _, pos := range positions {
		symbol := asString(pos["symbol"])
		side := strings.ToUpper(asString(pos["side"]))
		qty := formatFloat(pos["quantity"], 6)
		entry := formatFloat(pos["entry_price"], 4)
		mark := formatFloat(pos["mark_price"], 4)
		pnl := formatFloat(pos["unrealized_pnl"], 2)
		pnlPct := formatFloat(pos["unrealized_pnl_pct"], 2)
		lines = append(lines, fmt.Sprintf("%s %s\nQty: %s | Entry: %s\nMark: %s | PnL: %s (%s%%)", 
			symbol, side, qty, entry, mark, pnl, pnlPct))
	}

	b.reply(chatID, strings.Join(lines, "\n\n"))
}

func (b *Bot) handleOrders(chatID int64, args []string) {
	at, errMsg := b.resolveTrader(args)
	if errMsg != "" {
		b.reply(chatID, errMsg)
		return
	}

	orders, err := at.GetOpenOrders("")
	if err != nil {
		b.reply(chatID, fmt.Sprintf("❌ 获取挂单失败：%v", err))
		return
	}

	if len(orders) == 0 {
		b.reply(chatID, "📝 当前无挂单。")
		return
	}

	lines := []string{fmt.Sprintf("📝 交易员：%s 当前挂单", at.GetName())}
	for _, order := range orders {
		price := order.Price
		if order.StopPrice > 0 {
			price = order.StopPrice
		}
		lines = append(lines, fmt.Sprintf("%s %s %s\nType: %s | Qty: %.6f\nPrice: %.4f | Status: %s",
			order.Symbol, order.Side, order.PositionSide, order.Type, order.Quantity, price, order.Status))
	}
	b.reply(chatID, strings.Join(lines, "\n\n"))
}

func (b *Bot) handleDeferred(chatID int64, args []string) {
	_, errMsg := b.resolveTrader(args)
	if errMsg != "" {
		b.reply(chatID, errMsg)
		return
	}

	b.reply(chatID, "⏳ 暂不支持查看缓存止盈止损。")
	return
}

func (b *Bot) handleDecision(chatID int64, args []string) {
	at, errMsg := b.resolveTrader(args)
	if errMsg != "" {
		b.reply(chatID, errMsg)
		return
	}
	if b.store == nil {
		b.reply(chatID, "❌ 数据库未连接，无法获取决策。")
		return
	}

	records, err := b.store.Decision().GetLatestRecords(at.GetID(), 1)
	if err != nil || len(records) == 0 {
		b.reply(chatID, "📭 当前无决策记录。")
		return
	}

	rec := records[len(records)-1]
	b.reply(chatID, b.formatDecisionNotify(at.GetName(), rec))
}

func (b *Bot) handleAlerts(chatID int64, args []string) {
	at, errMsg := b.resolveTrader(args)
	if errMsg != "" {
		b.reply(chatID, errMsg)
		return
	}
	if b.store == nil {
		b.reply(chatID, "❌ 数据库未连接。")
		return
	}

	records, err := b.store.Decision().GetLatestRecords(at.GetID(), 20)
	if err != nil || len(records) == 0 {
		b.reply(chatID, "📭 当前无告警记录。")
		return
	}

	var alerts []string
	for i := len(records) - 1; i >= 0; i-- {
		rec := records[i]
		if rec.Success && rec.ErrorMessage == "" {
			continue
		}
		msg := fmt.Sprintf("#%d %s", rec.CycleNumber, rec.ErrorMessage)
		if rec.ErrorMessage == "" {
			msg = fmt.Sprintf("#%d 执行失败", rec.CycleNumber)
		}
		alerts = append(alerts, msg)
		if len(alerts) >= 5 {
			break
		}
	}

	if len(alerts) == 0 {
		b.reply(chatID, "✅ 最近无告警/错误。")
		return
	}
	lines := append([]string{fmt.Sprintf("🚨 交易员：%s 最近告警", at.GetName())}, alerts...)
	b.reply(chatID, strings.Join(lines, "\n"))
}

func (b *Bot) handleTraders(chatID int64) {
	traders := b.manager.GetAllTraders()
	if len(traders) == 0 {
		b.reply(chatID, "📭 当前没有交易员。")
		return
	}
	lines := []string{"👥 交易员列表："}
	ids := make([]string, 0, len(traders))
	for id := range traders {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		t := traders[id]
		status := t.GetStatus()
		running := false
		if v, ok := status["is_running"].(bool); ok {
			running = v
		}
		state := "🔴 停止"
		if running {
			state = "🟢 运行中"
		}
		lines = append(lines, fmt.Sprintf("%s | %s | %s", id[:8], t.GetName(), state))
	}
	b.reply(chatID, strings.Join(lines, "\n"))
}

func (b *Bot) resolveTrader(args []string) (*trader.AutoTrader, string) {
	traders := b.manager.GetAllTraders()
	if len(traders) == 0 {
		return nil, "📭 当前没有交易员。"
	}

	var pickedID string
	if len(args) > 0 {
		pickedID = strings.TrimSpace(args[0])
	}

	if pickedID == "" {
		b.stateMu.Lock()
		id := b.currentTraderID
		b.stateMu.Unlock()

		if id != "" {
			if t, ok := traders[id]; ok {
				return t, ""
			}
		}
		// If only one trader, default to it
		if len(traders) == 1 {
			for id, t := range traders {
				b.setCurrentTrader(id)
				return t, ""
			}
		}
		return nil, b.traderPickHint(traders)
	}

	if t, ok := traders[pickedID]; ok {
		b.setCurrentTrader(pickedID)
		return t, ""
	}

	var matchedID string
	for id, t := range traders {
		if strings.HasPrefix(id, pickedID) {
			matchedID = id
			_ = t
			break
		}
	}
	if matchedID != "" {
		b.setCurrentTrader(matchedID)
		return traders[matchedID], ""
	}

	for id, t := range traders {
		if strings.EqualFold(t.GetName(), pickedID) {
			b.setCurrentTrader(id)
			return t, ""
		}
	}

	return nil, b.traderPickHint(traders)
}

func (b *Bot) setCurrentTrader(id string) {
	b.stateMu.Lock()
	b.currentTraderID = id
	b.stateMu.Unlock()
}

func (b *Bot) traderPickHint(traders map[string]*trader.AutoTrader) string {
	lines := []string{"⚠️ 请指定交易员 ID 或名称，例如：/status <trader_id>。"}
	ids := make([]string, 0, len(traders))
	for id := range traders {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		t := traders[id]
		lines = append(lines, fmt.Sprintf("%s | %s", id[:8], t.GetName()))
	}
	return strings.Join(lines, "\n")
}

func (b *Bot) helpText() string {
	return strings.Join([]string{
		"Available Commands:",
		"/start - 启动机器人",
		"/help - 显示帮助",
		"/status [id] - 查看状态",
		"/balance [id] - 查看余额",
		"/positions [id] - 查看持仓",
		"/orders [id] - 查看挂单",
		"/deferred [id] - 查看缓存止盈止损",
		"/decision [id] - 查看最新 AI 决策",
		"/alerts [id] - 查看告警",
		"/traders - 列出所有交易员",
	}, "\n")
}

func (b *Bot) reply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.DisableWebPagePreview = true
	if _, err := b.bot.Send(msg); err != nil {
		logger.Warnf("⚠️ Telegram send failed: %v", err)
	}
}

func (b *Bot) initDecisionCursors() {
	if b.store == nil {
		return
	}
	for _, id := range b.manager.GetTraderIDs() {
		records, err := b.store.Decision().GetLatestRecords(id, 1)
		if err != nil || len(records) == 0 {
			continue
		}
		b.decisionAfter[id] = records[len(records)-1].ID
	}
}

func (b *Bot) watchDecisions() {
	if b.store == nil {
		return
	}
	// Poll every 10 seconds for new AI decisions
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Verify admin ID is set
		if b.adminID == 0 {
			continue
		}

		traders := b.manager.GetAllTraders()
		for id, t := range traders {
			lastID := b.getDecisionCursor(id)
			// Fetch new records since last cursor
			records, err := b.store.Decision().GetRecordsAfterID(id, lastID, 20)
			if err != nil || len(records) == 0 {
				continue
			}
			for _, rec := range records {
				// We don't skip empty decisions if there is an error message
				// or if user wants to see everything. For now, we skip "boring" cycles with no actions and no errors
				if b.shouldSkipDecisionNotify(rec) {
					b.setDecisionCursor(id, rec.ID)
					continue
				}
				
				msg := b.formatDecisionNotify(t.GetName(), rec)
				b.reply(b.adminID, msg)
				
				b.setDecisionCursor(id, rec.ID)
			}
		}
	}
}

func (b *Bot) shouldSkipDecisionNotify(rec *store.DecisionRecord) bool {
	if rec == nil {
		return true
	}
	// If there are actions (trades), satisfy user Request 3 "AI polling info and system results" 
	if len(rec.Decisions) > 0 {
		return false
	}
	// If there is an error, definitely notify
	if rec.ErrorMessage != "" {
		return false
	}
	// If Success is false, notify
	if !rec.Success {
		return false
	}
	
	// Otherwise (no actions, no error, success=true), it's a "silent" cycle (e.g. checking price, doing nothing)
	// To avoid spamming, we skip these unless configured otherwise.
	return true
}

func (b *Bot) formatDecisionNotify(traderName string, rec *store.DecisionRecord) string {
	status := "✅ 成功"
	if rec != nil && !rec.Success {
		status = "❌ 失败"
	}
	
	lines := []string{
		fmt.Sprintf("🧠 AI 决策通知 | %s", traderName),
		fmt.Sprintf("Cycle #%d | %s", rec.CycleNumber, status),
	}
	
	if rec.ErrorMessage != "" {
		lines = append(lines, fmt.Sprintf("⚠️ 错误：%s", rec.ErrorMessage))
	}
	
	if len(rec.Decisions) == 0 {
		// Just in case we are notifying about a failure with no actions
		if rec.ErrorMessage == "" {
			lines = append(lines, "无交易执行")
		}
	} else {
		for _, action := range rec.Decisions {
			flag := "✅"
			if !action.Success {
				flag = "❌"
			}
			// Don't truncate reasoning in Telegram message, show full reasoning
			reason := action.Reasoning
			lines = append(lines, fmt.Sprintf("%s %s %s\nReason: %s", flag, action.Symbol, action.Action, reason))
		}
	}

	// Add Chain of Thought (CoT) if available
	if rec.CoTTrace != "" {
		lines = append(lines, "")
		lines = append(lines, "🧠 AI思维链:")
		// Truncate if too long (Telegram message limit is 4096 chars)
		// We reserve some space for other parts of the message
		const maxCoTLength = 3000
		cot := rec.CoTTrace
		if len(cot) > maxCoTLength {
			cot = cot[:maxCoTLength] + "\n...(truncated)"
		}
		lines = append(lines, cot)
	}
	
	return strings.Join(lines, "\n")
}

func (b *Bot) getDecisionCursor(traderID string) int64 {
	b.decisionMu.Lock()
	defer b.decisionMu.Unlock()
	return b.decisionAfter[traderID]
}

func (b *Bot) setDecisionCursor(traderID string, id int64) {
	b.decisionMu.Lock()
	b.decisionAfter[traderID] = id
	b.decisionMu.Unlock()
}

// Helper functions (duplicated from util to avoid circular imports or just simple local helpers)

func formatFloat(value interface{}, precision int) string {
	switch v := value.(type) {
	case float64:
		return fmt.Sprintf("%.*f", precision, v)
	case float32:
		return fmt.Sprintf("%.*f", precision, v)
	case int:
		return fmt.Sprintf("%.*f", precision, float64(v))
	case int64:
		return fmt.Sprintf("%.*f", precision, float64(v))
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return v
		}
		return fmt.Sprintf("%.*f", precision, parsed)
	default:
		return "0"
	}
}

func asString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case fmt.Stringer:
		return v.String()
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return fmt.Sprintf("%.4f", v)
	default:
		return ""
	}
}
