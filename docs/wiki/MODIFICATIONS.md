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

---
*文档生成时间: 2026-02-12*
