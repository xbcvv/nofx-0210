package kernel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// RegisterConfig 寄存器配置
type RegisterConfig struct {
	Enabled       bool `json:"enabled"`        // 是否启用寄存器
	MaxRecords    int  `json:"max_records"`    // 最大记录数
	IncludeDecisions bool `json:"include_decisions"` // 是否包含完整决策
	IncludeMarketData bool `json:"include_market_data"` // 是否包含市场数据
}

// RegisterRecord 寄存器记录
type RegisterRecord struct {
	Timestamp       string     `json:"timestamp"`        // 时间戳
	Cycle           int        `json:"cycle"`            // 周期数
	Decisions       []Decision `json:"decisions,omitempty"` // 决策记录
	MarketData      map[string]interface{} `json:"market_data,omitempty"` // 市场数据
	ExecutionStatus string     `json:"execution_status"` // 执行状态
	MarketRegime    string     `json:"market_regime"`    // 市场状态
}

// Register 寄存器管理
type Register struct {
	traderID string
	config   RegisterConfig
}

// NewRegister 创建新的寄存器实例
func NewRegister(traderID string, config RegisterConfig) *Register {
	return &Register{
		traderID: traderID,
		config:   config,
	}
}

// GetRegisterPath 获取寄存器文件路径
func (r *Register) GetRegisterPath() string {
	dataDir := "data/decision_history"
	// 确保目录存在
	os.MkdirAll(dataDir, 0755)
	return filepath.Join(dataDir, fmt.Sprintf("%s.json", r.traderID))
}

// LoadRecords 加载寄存器记录
func (r *Register) LoadRecords() ([]RegisterRecord, error) {
	path := r.GetRegisterPath()
	
	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []RegisterRecord{}, nil
	}
	
	// 读取文件
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read register file: %w", err)
	}
	
	// 解析JSON
	var records []RegisterRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("failed to parse register file: %w", err)
	}
	
	return records, nil
}

// SaveRecord 保存寄存器记录
func (r *Register) SaveRecord(record RegisterRecord) error {
	// 加载现有记录
	records, err := r.LoadRecords()
	if err != nil {
		return err
	}
	
	// 添加新记录到开头
	records = append([]RegisterRecord{record}, records...)
	
	// 限制记录数量
	if len(records) > r.config.MaxRecords {
		records = records[:r.config.MaxRecords]
	}
	
	// 按时间戳排序（最新的在前）
	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp > records[j].Timestamp
	})
	
	// 保存到文件
	path := r.GetRegisterPath()
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal register records: %w", err)
	}
	
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write register file: %w", err)
	}
	
	return nil
}

// GetRecentRecords 获取最近的记录
func (r *Register) GetRecentRecords(limit int) ([]RegisterRecord, error) {
	records, err := r.LoadRecords()
	if err != nil {
		return nil, err
	}
	
	if limit <= 0 || limit > len(records) {
		limit = len(records)
	}
	
	return records[:limit], nil
}

// ClearRecords 清空寄存器记录
func (r *Register) ClearRecords() error {
	path := r.GetRegisterPath()
	
	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	
	// 删除文件
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to clear register records: %w", err)
	}
	
	return nil
}

// BuildRegisterPrompt 构建寄存器提示词
func (r *Register) BuildRegisterPrompt() (string, error) {
	if !r.config.Enabled {
		return "", nil
	}
	
	records, err := r.GetRecentRecords(5) // 默认获取最近5条
	if err != nil {
		return "", err
	}
	
	if len(records) == 0 {
		return "", nil
	}
	
	var prompt string
	prompt += "## 历史决策记录\n"
	
	for i, record := range records {
		prompt += fmt.Sprintf("%d. 时间: %s | 状态: %s | 市场: %s\n", 
			i+1, record.Timestamp, record.ExecutionStatus, record.MarketRegime)
		
		if r.config.IncludeDecisions {
			for _, decision := range record.Decisions {
				prompt += fmt.Sprintf("   - %s %s", decision.Symbol, decision.Action)
				if decision.StopLoss > 0 {
					prompt += fmt.Sprintf(" SL:%.2f", decision.StopLoss)
				}
				if decision.TakeProfit > 0 {
					prompt += fmt.Sprintf(" TP:%.2f", decision.TakeProfit)
				}
				prompt += "\n"
			}
		}

		if r.config.IncludeMarketData && len(record.MarketData) > 0 {
			prompt += "   - Market Context:\n"
			var keys []string
			for k := range record.MarketData {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			
			for _, k := range keys {
				if dataMap, ok := record.MarketData[k].(map[string]interface{}); ok {
					price, _ := dataMap["price"].(float64)
					rsi, _ := dataMap["rsi"].(float64)
					change4h, _ := dataMap["change_4h"].(float64)
					prompt += fmt.Sprintf("     * %s: Price=%.2f, RSI=%.2f, 4h=%.2f%%\n", 
						k, price, rsi, change4h)
				}
			}
		}
	}
	
	prompt += "\n"
	return prompt, nil
}

// CreateRecordFromContext 从上下文创建寄存器记录
func CreateRecordFromContext(ctx *Context, decisions []Decision, executionStatus string) RegisterRecord {
	record := RegisterRecord{
		Timestamp:       time.Now().Format(time.RFC3339),
		Cycle:           ctx.CallCount,
		Decisions:       decisions,
		ExecutionStatus: executionStatus,
		MarketRegime:    "normal", // 默认市场状态
		MarketData:      make(map[string]interface{}),
	}
	
	// 简单判断市场状态
	if btcData, hasBTC := ctx.MarketDataMap["BTCUSDT"]; hasBTC {
		if btcData.PriceChange4h > 5 {
			record.MarketRegime = "strong_uptrend"
		} else if btcData.PriceChange4h < -5 {
			record.MarketRegime = "strong_downtrend"
		} else if btcData.CurrentRSI7 > 70 {
			record.MarketRegime = "overbought"
		} else if btcData.CurrentRSI7 < 30 {
			record.MarketRegime = "oversold"
		}
		
		// 记录BTC市场数据作为基准
		record.MarketData["BTCUSDT"] = map[string]interface{}{
			"price":     btcData.CurrentPrice,
			"rsi":       btcData.CurrentRSI7,
			"change_4h": btcData.PriceChange4h,
		}
	}
	
	// 记录决策涉及币种的市场数据
	for _, decision := range decisions {
		if decision.Symbol == "BTCUSDT" {
			continue // 已经记录过
		}
		if data, ok := ctx.MarketDataMap[decision.Symbol]; ok {
			record.MarketData[decision.Symbol] = map[string]interface{}{
				"price":     data.CurrentPrice,
				"rsi":       data.CurrentRSI7,
				"change_4h": data.PriceChange4h,
			}
		}
	}
	
	return record
}
