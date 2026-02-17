# 演练 - 系统修复与增强

## 1. 持仓时间计算修复 (Critical)

### 问题描述
系统重启或订单更新时，仓位的 `EntryTime` 会丢失，导致持仓时间重置为零。
- **根本原因**: 数据库存储的方向字段为大写 `LONG`/`SHORT`，但系统使用小写 `long`/`short` 进行查询，导致 `GetOpenPositionBySymbol` 返回 `record not found`。
- **影响**: 仓位管理逻辑无法追踪持仓时间，可能导致过早离场或盈亏追踪错误。

### 解决方案
- **文件**: `store/position.go`
- **更改**: 在 `GetOpenPositionBySymbol` 中将 `side` 参数统一标准化为大写。
```go
func (s *PositionStore) GetOpenPositionBySymbol(symbol, side string) (*TraderPosition, error) {
    side = strings.ToUpper(side) // Fix: Normalize case
    // ...
}
```

## 2. 15分钟涨跌幅实现 (Feature)

### 问题描述
AI 之前需要从原始 K 线数据计算涨跌幅，容易产生幻觉，且系统仅提供 1h 和 4h 涨跌幅，导致短线趋势判断不准。

### 解决方案
- **市场数据**: 在 `market.Data` 结构体中新增 `PriceChange15m` 字段。
- **计算逻辑**: 在 `market/data.go` 中实现计算逻辑（基于最近 5 根 3 分钟 K 线或直接计算）。
- **Prompt 工程**: 更新 `kernel/formatter.go`，显式向 AI 提供此数据。
  - **修改前**: `Current Price: 100.00`
  - **修改后**: `Current Price: 100.00 | 15m Change: +1.20% | 1h Change: +3.50%`
- **文档**: 更新了 `docs/wiki/INDICATORS.md` 和 `DATA_DICTIONARY.md`。

## 3. 全局市场背景 (Feature)

### 问题描述
`prompt23.yaml` 中的 "Global Command" 逻辑需要 BTC 状态（ADX, 涨跌幅）来判断市场气候。但如果 BTC 不在 `CandidateCoins` 列表（例如仅分析 ETH），AI 将无法看到 BTC 数据，导致无法执行全局指挥逻辑。

### 解决方案
- **引擎**: 修改 `kernel/engine.go`，强制**始终**获取 `BTCUSDT` 市场数据并添加到 `MarketDataMap`，即使它不是候选币种。
- **格式化**: 更新 `kernel/formatter.go` (`formatContextData`)，在系统提示词顶部增加 **"Global Market Context"** 版块。
- **结果**: 即使分析山寨币，AI 也能看到：
  ```
  ## Global Market Context (Reference Only, NOT for Trading)
  **BTCUSDT**:
  - Price: 68000.00
  - Change: 15m -0.50% ...
  - ADX: 45.2
  ```
- **文档**: 更新 `docs/wiki/INDICATORS.md`。

## 4. 文档更新 (Wiki)

### 更新内容
- **`docs/wiki/INDICATORS.md`**: 增加 `Global Market Context` 章节和 `PriceChange15m` 详情。
- **`docs/wiki/DATA_DICTIONARY.md`**: 在市场数据表中增加 `PriceChange15m`。
- **`docs/wiki/MODIFICATIONS.md`**: 记录了 Section 6 的新功能。

## 5. 时区验证 (Audit)

### 审计结果
- **分析**: 确认所有 AI 分析均基于 K 线的 **UTC 时间戳** (`time.UnixMilli(k.Time).UTC()`)。
- **结论**: 用户机器上的本地时间显示 (UTC+8) 不影响内部逻辑或 AI 决策。

## 6. 功能废弃 (Deprecation)

### 行动
- **废弃**: "Global Command" 功能文档 (`task/task_archives/`)。
  - 原因: 该功能曾有文档但未在当前代码库中实现。

## 7. 配置驱动架构 (Refactor)

### 问题描述
之前 EMA20, RSI7 等指标是硬编码的。修改周期（如 EMA60）需要修改代码。
AI 对 "日涨跌幅" 的逻辑判断存在歧义。

### 解决方案
- **数据结构**: `market/types.go` 现支持动态 `PriceChanges`, `EMAs`, `ATRs` 映射。
- **引擎**: `kernel/engine.go` 读取 `strategy.json` 数组 (`selected_timeframes`, `ema_periods` 等) 并驱动计算。
- **格式化**: `kernel/formatter.go` 遍历这些映射生成 Prompt 内容。
- **Schema**: 更新 `kernel/schema.go` 定义 `Change_{tf}`, `EMA{period}`, `ATR{period}`。
- **结果**: 用户只需在配置中添加 `1d`，AI 即可立即获得 `Change_1d` 数据。

## 8. K线显示优化 (Enhancement)

### 问题描述
系统在 Prompt 中将 K 线显示数量硬编码限制为 **30 根**（可能是为了节省 Token）。这导致 AI 只能看到 15m 周期约 7.5 小时的数据，丢失了关键的日内结构（如亚盘高低点或昨日收盘价）。

### 解决方案
- **配置**: 在 `RightSide_Flow_Hunter_V3.0.json` 中添加 `display_count`。
- **逻辑**: 更新 `kernel/formatter.go` 使用此配置值（若未设置则默认为 30）。
- **默认值**: 在策略文件中将默认值设为 **60** (15 小时)，使 AI 能在 15m 图表上看到完整的日内结构。

## 9. 修复全局锁失效问题 (Global Lock Fix)

### 问题描述
用户反馈在策略设定“任意持仓 < 45m 禁开新仓”的情况下，AI 仍然尝试开新仓。

### 根本原因
提供给 AI 的 Prompt 中，持仓信息只包含价格和盈亏，缺失了 **持仓时长 (Hold Duration)**。AI 无法判断持仓是否满足 45分钟的限制。

### 修复方案
- 修改 `kernel/formatter.go`，在 `formatCurrentPositions` 函数中增加时长计算。
- 现在 AI 会看到如下信息：
  ```text
  1. INITUSDT LONG | ... | ⏱️ 持仓时间: 10m (进场: 12:40)
  ```
- AI 现在可以准确执行 `if duration < 45m then STOP` 的逻辑。

### 相关文件
- `kernel/formatter.go` [查看修改](file:///d:/code/nofx-0210/kernel/formatter.go)
