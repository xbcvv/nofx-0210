package kernel

// ============================================================================
// Trading Data Schema - 交易数据字典 (Optimized for Tokens)
// ============================================================================
// 使用精简英文定义以减少 Token 消耗，保留关键逻辑说明
// ============================================================================

const (
	SchemaVersion = "1.1.0"
)

// Language 语言类型
type Language string

const (
	LangChinese Language = "zh-CN"
	LangEnglish Language = "en-US"
)

// ========== 极简 Prompt 生成函数 ==========

// GetSchemaPrompt 生成Schema说明文本，用于AI Prompt
// 无论 lang 是中文还是英文，统一返回精简的英文 Schema，以节省 Token。
// AI 的理解能力足以处理英文 Schema，回复语言由 System Prompt 中的角色定义控制。
func GetSchemaPrompt(lang Language) string {
	return getSchemaPromptConcise()
}

// getSchemaPromptConcise 生成极简版英文 Prompt
func getSchemaPromptConcise() string {
	return `# Data Schema & Rules (Units: USDT, %)

## Field Definitions
- **Account**: Equity(Net Worth), Balance(Available), Margin(Used Rate), PnL(Total Return)
- **Trade**: Entry(Avg Price), Exit(Avg Price), Profit(Realized PnL), HoldDuration(Time)
- **Position**: 
  - UnrealizedPnL(Floating), PeakPnL(Max Historical Floating, vital for trailing stop)
  - Drawdown(Pullback from Peak), LiqPrice(0=Safe)
  - Action: open_long/short, close_long/short, partial_close(use ClosePercentage 0.0-1.0), hold, wait
  - ClosePercentage: 0.5=50%, 1.0=100% (Only for partial_close)
- **Register**: Cycle(Seq Num), MarketRegime(AI View), ExecutionStatus(Last Result)
- **Market Data**:
- **Market Data**:
  - Change_{tf}: 周期涨跌幅/Price Change % (e.g. Change_1d=日涨跌幅/24h change)
  - EMA{period}: 指数均线/Exp Moving Avg (Value=点位, Slope=斜率, Spread=乖离/Spread)
  - ATR{period}: 平均真实波幅/Average True Range

## Rules
- **OI Logic**: OI↑+Price↑=Bull Inflow; OI↑+Price↓=Bear Inflow; OI↓=Covering/Reversal
`
}
