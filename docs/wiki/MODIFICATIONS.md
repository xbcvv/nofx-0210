# 项目修改详细记录 (Wiki)

本文档详细记录了基于 `NoFxAiOS/nofx` main 分支（原始版本）进行的自定义修改。

## 1. 核心交易逻辑修改

### 1.1 支持 Hold 状态下的止盈止损更新
**修改文件**: `trader/auto_trader.go`
**原逻辑**: 当 AI 决策为 `hold` 或 `wait` 时，直接返回，不执行任何操作。
**新逻辑**: 在 `executeHoldWithRecord` 方法中增加了逻辑，如果 AI 决策中包含了 `StopLoss` 或 `TakeProfit`，即使动作为 `hold`，也会尝试更新当前持仓的止盈止损设置。
**代码片段**:
```go
// executeHoldWithRecord executes hold action but checks for updates (SL/TP)
func (at *AutoTrader) executeHoldWithRecord(decision *kernel.Decision, actionRecord *store.DecisionAction) error {
    // If StopLoss or TakeProfit is provided, update them
    if decision.StopLoss > 0 || decision.TakeProfit > 0 {
        // ... (更新止盈止损逻辑)
    }
    return nil
}
```

### 1.2 支持分批止盈 (Partial Take Profit)
**修改文件**: `trader/auto_trader.go`
**原逻辑**: 仅支持完全平仓 (`close_long`, `close_short`)。
**新逻辑**: 
1.  在 `executeDecisionWithRecord` 中增加了 `partial_close` 动作的支持。
2.  新增 `executePartialCloseWithRecord` 方法，根据 `decision.ClosePercentage` 计算平仓数量。
3.  在 `executeCloseLongWithRecord` 和 `executeCloseShortWithRecord` 中增加了处理 `ClosePercentage` 的逻辑，支持按比例平仓。

### 1.3 禁用强制回撤平仓机制
**修改文件**: `trader/auto_trader.go`
**原逻辑**: 在 `checkPositionDrawdown` 方法中，如果持仓盈利超过 5% 且后续回撤超过 40%，会强制触发紧急平仓（"Hidden Stop-Loss"）。
**新逻辑**: 该强制平仓逻辑已被注释掉，仅保留日志记录。系统不再自动因回撤而强制平仓，完全交由 AI 决策或设置的硬性止损来控制。
**代码片段**:
```go
// Check close position condition: profit > 5% and drawdown >= 40%
// [REMOVED BY USER REQUEST] "Hidden Stop-Loss" lock removed
// if currentPnLPct > 5.0 && drawdownPct >= 40.0 {
//     ...
// }
```

## 2. 新增功能模块

### 2.1 决策寄存器 (Registers)
**新增文件**: `kernel/register.go`
**修改文件**: `trader/auto_trader.go`
**功能描述**: 
引入了“寄存器”概念，用于将 AI 的历史决策、执行状态和当时的市场数据保存为 JSON 文件（存储在 `data/decision_history/`）。这允许系统在后续的 Prompt 中引用之前的决策上下文，增强 AI 的连续性思维能力。
**实现细节**:
-   `kernel/register.go`: 定义了 `Register` 结构体，实现了记录的加载、保存和 Prompt 构建。
-   `trader/auto_trader.go`: 在 `runCycle` 结束前，调用寄存器保存当前周期的完整上下文。

### 2.2 Telegram 交互机器人
**新增目录**: `telegram/`
**新增文件**: `telegram/bot.go`
**修改文件**: `main.go`
**功能描述**: 
集成了一个 Telegram Bot，用于实时接收交易通知和进行简单的交互查询。
**主要功能**:
-   **实时通知**: 每当 AI 完成一次决策周期，Bot 会向指定 Chat ID 发送包含决策结果、理由和思维链（CoT）的通知。
-   **命令交互**: 支持 `/status` (查看状态), `/balance` (查看余额), `/positions` (查看持仓) 等命令。
**实现细节**:
-   `telegram/bot.go`: 实现了 Bot 的核心逻辑，包括消息处理和格式化。
-   `main.go`: 在系统启动时初始化并启动 Telegram Bot。

## 3. 其他修改

### 3.1 移除 WebSocket 行情依赖
**修改文件**: `main.go`
**描述**: 注释掉了 `market.NewWSMonitor` 的启动代码，并添加日志 "Using CoinAnk API for all market data (WebSocket cache disabled)"。现在系统完全依赖 CoinAnk API 获取 K 线数据，减少了连接维护开销。

### 3.2 依赖库更新
**修改文件**: `go.mod`, `go.sum`
**描述**: 引入了 `github.com/go-telegram-bot-api/telegram-bot-api/v5` 等新依赖以支持 Telegram 功能。

### 3.3 数据字典与交易规则更新
**修改文件**: `kernel/schema.go`, `kernel/engine.go`
**描述**: 
1.  **扩展 Action 字段**: 更新了数据字典中的 `Action` 描述，明确支持 `open_long`, `open_short`, `close_long`, `close_short`, `partial_close`, `hold`, `wait` 等动作。
2.  **新增 ClosePercentage 字段**: 在 `PositionMetrics` 下新增了 `ClosePercentage` 字段定义，用于支持分批止盈策略（0.5 表示平仓 50%）。
3.  **新增 RegisterMetrics 数据字典**: 定义了 `Cycle`, `MarketRegime`, `ExecutionStatus`, `Decisions` 等寄存器字段，使 AI 能够理解历史记录（记忆模块）的含义。
4.  **Prompt 同步**: 上述修改会自动反映在 `GetSchemaPrompt` 生成的系统 Prompt 中，确保 AI 能够理解最新的交易指令集和寄存器内容。
5.  **Trader ID 修复**: 修复了 `kernel/engine.go` 和 `trader/auto_trader.go` 中上下文传递 Trader ID 的问题，确保每个交易员读取正确的寄存器历史。

### 3.4 寄存器 Prompt 优化
**修改文件**: `kernel/register.go`
**描述**: 将寄存器历史记录的标题从 `## 历史决策记录` 修改为 **`## 🧠 决策寄存器 (Memory Bank)`**，并配合 System Prompt 中的数据字典，使 AI 能够更精准地识别和利用历史决策数据作为长期记忆。

### 3.5 AI Prompt 性能优化与自适应显示
**修改文件**: `kernel/engine.go`, `store/strategy.go`
**描述**:
1.  **Token 消耗优化**: 重构了市场数据格式化逻辑 (`formatTimeframeSeriesData`)，将 Prompt 中展示的 K 线数量与后端指标计算所需的 K 线数量解耦。
    *   后端计算：始终使用配置的完整数量（如 100 根）以确保 EMA/MACD 指标精准。
    *   前端展示：默认截断为主周期 24 根、辅周期 12 根，大幅降低 Token 消耗（~60-70%）。
2.  **自适应参数 (`DisplayCount`)**:
    *   在 `StrategyConfig` 中新增 `display_count` 字段。
    *   支持用户在策略 JSON 中自定义展示数量（如 50 或 300），系统会自动适配。若未设置，则启用默认的智能截断逻辑。
3.  **去重**: 修复了 System Prompt 与 Custom Prompt 规则重复导致上下文冗余的问题。

### 3.6 Prompt Data Format Optimization
**修改文件**: `kernel/engine.go`
**描述**:
1.  **紧凑 CSV 格式**: 将 K 线数据的展示格式从对齐表格（固定宽度）改为紧凑的 CSV 格式（逗号分隔），大幅减少了由空格填充导致的 Token 消耗。
2.  **智能浮点数格式化**: 实现了 `formatPriceSmart` 和 `formatVolumeSmart` 函数，根据数值大小动态调整小数精度（如小数值保留 8 位，大数值保留 2 位），避免了冗余精度（如 `0.000000` -> `0`）。
3.  **数组压缩**: 优化了指标数组（EMA, MACD, RSI）的展示，移除了方括号和逗号，改用空格分隔，进一步节省 Token。

### 3.7 System Prompt Schema Optimization
**修改文件**: `kernel/schema.go`, `kernel/engine.go`
**描述**:
1.  **英文数据字典**: 将 System Prompt 中的数据字典（Schema）从详细的双语对照改为精简的英文定义（如 `UnrealizedPnL(Floating)`），减少了约 60% 的 Schema Token 消耗。
2.  **强制中文回复**: 在 `engine.go` 的 Prompt 构建逻辑中，增加了针对中文环境的强制指令 `IMPORTANT: Please analyze and respond STRICTLY in Chinese language`，确保即使输入英文定义，AI 仍用中文回复。
3.  **Wiki 文档支持**: 自动生成了 `docs/wiki/DATA_DICTIONARY.md`，提供中英文术语对照表，方便开发者查阅。

### 3.8 Configurable Risk & Strategy Rules
**修改文件**: `kernel/schema.go`, `docs/wiki/DATA_DICTIONARY.md`
**描述**:
1.  **移除硬编码规则**: 从 System Prompt (`kernel/schema.go`) 中彻底移除了硬编码的风险控制 (Risk Management) 和策略原则 (Strategy) 规则（如 "Max Margin 30%", "Only add to winners"）。
2.  **配置驱动**: 将风控和策略逻辑的控制权完全交还给用户。现在，AI 的行为完全取决于策略配置文件 (`strategy_*.json`) 中的 `entry_standards` 和 `decision_process` 字段，实现了真正的策略动态配置。
3.  **文档更新**: 在 Wiki 文档中标记了这些规则已从底层移除，并提供了配置建议。

## 4. 系统稳定性与功能增强 [2026-02-14]

### 4.1 Telegram Bot 重构与配置集成
**修改文件**: `telegram/bot.go`, `config/config.go`, `main.go`, `.env`
**描述**: 
1.  **代码重构**: 优化了 Bot 的初始化逻辑，使其更符合项目整体架构。
2.  **配置集成**: 将 Bot Token 和 User ID 从代码中剥离，统一通过 `config` 包读取环境变量 (`TELEGRAM_BOT_TOKEN`, `TELEGRAM_USER_ID`)。
3.  **启动集成**: 在 `main.go` 中正式集成 Bot 启动流程，确保随主程序一同运行。

### 4.2 修复 Hold 状态下的订单更新逻辑
**修改文件**: `trader/auto_trader.go`
**描述**: 
修复了在 `Hold` 状态下尝试更新止损/止盈时，因未先撤销旧订单而导致 Binance 报错 `-4130` (Existing Order) 的问题。
**新逻辑**: 引入“先撤后设”机制，确保新订单能正确覆盖旧订单。

### 4.3 修复 Partial Close 订单更新与安全隐患
**修改文件**: `trader/auto_trader.go`, `trader/binance/futures.go`
**描述**: 
1.  **逻辑补全**: 完善了分批平仓 (`Partial Close`) 逻辑。现在执行分批平仓后，会自动更新剩余仓位的止损/止盈设置（如有）。
2.  **重大 Bug 修复**: 修正了 Binance 接口驱动中 `CloseLong/CloseShort` 无条件撤销所有挂单的缺陷。
    - **修复前**: 无论平多少仓位，都会撤销该币种所有挂单，导致剩余仓位“裸奔”。
    - **修复后**: 仅在全仓平仓 (`Full Close`) 时撤销挂单；分批平仓时保留现有止损/止盈单，确保安全。

### 4.4 构建系统优化
**修改文件**: `.github/workflows/docker-build.yml`, `Dockerfile`
**描述**: 
1.  **构建修复**: 解决了 `go mod tidy` 在 Docker 构建中的兼容性问题。
2.  **资源优化**: 在 GitHub Actions 中启用了 `max-parallel: 1` 策略，避免了因并行构建导致的资源耗尽和超时失败。
3.  **环境对齐**: 同步了 `main` 与 `test` 分支的构建环境，确保生产环境稳定性。

### 4.5 记忆解析升级 (Memory Audit Upgrade)
**修改文件**: `kernel/engine.go`, `trader/auto_trader.go`, `kernel/register.go`
**描述**: 
1.  **数据层**: 在 `Decision` 结构体中新增 `EntryPrice` 字段，用于记录决策当时的持仓均价。
2.  **执行层**: 在 `auto_trader` 生成决策记录时，若当前持有相应仓位，则自动抓取 `Position.EntryPrice` 并填充到决策记录中。
3.  **认知层**: 优化了 Prompt 中“决策寄存器”的展示格式，在 Hold/Open 记录中显式展示 `Entry: xxx`，使 AI 能够进行准确的盈亏审计和反思。

---

## 5. 新增功能模块 [2026-02-15]

### 5.1 ADX (平均方向指数) 指标集成
**修改文件**: `market/adx.go`, `market/data.go`, `market/types.go`, `kernel/engine.go`, `store/strategy.go`, `web/src/components/strategy/IndicatorEditor.tsx`
**描述**: 
1.  **核心计算**: 实现了 Wilder's Smoothing 算法计算 ADX, DI+, DI- (`market/adx.go`)。
2.  **数据流**: 将 ADX 集成到 Market Data 流程中，支持从 1m 到 1d 所有周期的计算。
3.  **AI 感知**: 更新了 `kernel/engine.go`，当策略开启 `EnableADX` 时，自动将 ADX 值注入 System Prompt。
4.  **前端支持**: 在策略编辑器中增加了 ADX 开关和周期配置（默认 14）。
5.  **配置支持**: 策略配置文件 (`store/strategy.go`) 新增 `EnableADX` 和 `ADXPeriods` 字段。

### 5.2 Wiki 文档系统升级
**新增文件**: `docs/wiki/INDICATORS.md`, `docs/wiki/README.md`
**修改文件**: `web/src/components/landing/FooterSection.tsx`
**描述**: 
1.  **指标库文档**: 创建了 `INDICATORS.md`，详细记录了所有支持的指标（ADX, EMA, MACD, RSI, ATR, BOLL 等）的配置参数、默认值及 Prompt 呈现格式（中文版）。
2.  **Wiki 索引**: 创建了首页索引 `README.md`。
3.  **入口集成**: 在网站页脚 (Footer) 新增了 "Wiki / 指标说明" 链接，方便用户查阅。

## 6. 新增功能模块 [2026-02-16]

### 6.1 Global Market Context (全局市场背景)
**修改文件**: `kernel/engine.go`, `kernel/formatter.go`
**描述**: 
1.  **强制 BTC 数据获取**: 修改 `kernel/engine.go`，在获取市场数据时，无论当前策略如何配置，**始终强制获取 BTCUSDT 的 K 线数据**并加入 `MarketDataMap`。
2.  **Prompt 头部展示**: 修改 `kernel/formatter.go`，在 System Prompt 的最顶部（Header 之后）插入 "Global Market Context" 版块，展示 BTCUSDT 的价格、15m/1h/4h 涨跌幅和 ADX。
3.  **目的**: 确保 AI 在分析任何币种（如 ETH, SOL）时，都能获得 BTC 的实时数据作为“全局指挥”判断的依据，彻底解决了 `prompt23.yaml` 中全局指挥逻辑的数据依赖问题。

### 6.2 15m Price Change (15分钟涨跌幅)
**修改文件**: `market/types.go`, `market/data.go`, `kernel/formatter.go`
**描述**: 
1.  **字段扩展**: 在 `market.Data` 结构体中新增 `PriceChange15m` 字段。
2.  **计算逻辑**: 在 `market/data.go` 中实现了基于最近 5 根 3m K 线的涨跌幅计算逻辑。
3.  **Prompt 呈现**: 在 Prompt 中明确展示 `15m Change`，为 AI 判断“恶性暴跌”提供精确数值，消除了 AI 需要从 K 线列表自行计算而产生的幻觉风险。

## 7. Configuration-Driven Architecture (V1.1.0)
- **Goal**: Enable zero-code strategy iteration by driving data generation via `strategy.json`.
- **Files**: `market/types.go`, `market/data.go`, `kernel/engine.go`, `kernel/formatter.go`, `kernel/schema.go`
- **Details**:
  1. **Dynamic Data**: `market.Data` now supports `PriceChanges` (Map), `DynamicEMAs` (Map), `DynamicATRs` (Map).
  2. **Config Driven**: `GetWithTimeframes` accepts `emaPeriods` and `selected_timeframes` from config.
  3. **Formatter**: Automatically iterates maps to generate Prompt (e.g. `EMA60: ...`).
  4. **Schema**: Updated System Prompt to define `Change_{tf}`, `EMA{period}`, `ATR{period}` with Chinese support.

## 8. K-Line Display Optimization (Enhancement)
**修改文件**: `fsdownload/RightSide_Flow_Hunter_V3.0.json`, `kernel/formatter.go`, `store/strategy.go`, `kernel/engine.go`
**描述**:
1.  **配置化**: 在策略配置文件的 `klines` 节点下新增 `display_count` 字段（默认 60）。
2.  **逻辑优化**: `kernel/formatter.go` 不再硬编码显示 30 根 K 线，而是读取配置值。
3.  **目的**: 让 AI 能看到更长周期的 K 线结构（如 60根 15m K线 = 15小时），从而更好地判断日内趋势和支撑阻力位，避免因数据截断导致的误判。

## 9. 修复全局锁失效问题 (Global Lock Fix)
**修改文件**: `kernel/formatter.go`
**描述**:
1.  **问题**: AI 无法遵守 "持仓 < 45m 禁开新仓" 的规则，因为 Prompt 中并未提供持仓时长信息。
2.  **修复**: 修改 `formatCurrentPositions`，在持仓信息中显式增加 `⏱️ 持仓时间 (Hold Duration)` 和 `进场时间 (Entry Time)`。
3.  **效果**: AI 现在能明确看到 `持仓时间: 15m`，从而正确触发规则拦截。

## 10. 全局指令配置化 (Global Command Config) - [Planned]
**修改文件**: `store/strategy.go` (预埋)
**描述**:
预埋了 `GlobalCommandConfig` 结构，为后续将“全局指挥”逻辑从代码硬编码迁移到 JSON 配置做准备。



