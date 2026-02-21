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

## 11. Frontend Mobile Adaptation
**修改文件**: `web/src/pages/StrategyStudioPage.tsx`, `web/src/pages/TraderDashboardPage.tsx`, `web/src/App.tsx`
**描述**:
1. **交易面板**: 优化了移动端排版，使卡片在小屏幕上自动变为单列布局。
2. **策略工作室**: 让左侧导航栏在移动端变为顶部横向滑动。
3. **全局防溢出**: 修复了 App 级别 X 轴滚动条导致页面晃动的问题。


## 12. Context Prompt Optimization & GitHub Action Fix [2026-02-22]

### 12.1 GitHub Actions Go Version Alignment
**修改文件**: `.github/workflows/test.yml`, `.github/workflows/pr-checks.yml`, `.github/workflows/pr-checks-run.yml`
**描述**: 
1. **统一环境**: 修复了由于 `go.mod` 要求 Go 1.25，而 GitHub Actions 仍停留在老旧版本（1.21/1.23）导致的 `go mod download` 和 `go build` 完全罢工甚至无法识别语法环境的核心问题。
2. **清除无用导入**: 修正了部分源码（如 `filter/coin_filter_manager.go`）中遗留的未使用的导入包（如 `sort`），解决了 Go 编译器严格报错的问题。

### 12.2 Inject 1H Price Change
**修改文件**: `market/data.go`
**描述**:
1. **零性能损耗注入**: 停止通过发送全量 60 根 1H K 线（高达 2000+ Tokens）引发上下文冗余的方式，而是直接在底层利用内存里的分时 K 曲线精确倒推前 60 分钟的价格差。
2. **格式化输出**: 修改 `Format()` 方法。现在无论是针对大盘（BTC/ETH）还是具体山寨币（如 PEPE, WIF），在输送给大模型的单一资产档案最头部，都会强制附带硬编码字符：`(1h Change: +X.XX%)`。
3. **目的**: 完美缝合了 \prompt.yaml\ 气候控制中关于 \[资金溢出]\ 和 \1h涨>0.5%\ 的识别盲点，完全杜绝了 AI 对 1H 涨跌幅的心算幻觉。

