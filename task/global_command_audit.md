# 全局指挥 (相对强弱气候) 策略审计报告 (Global Command Strategy Audit)

## 1. 现状评估 (Current Status Assessment)

### 数据可用性 (Data Availability)
-   **缺失 (Missing)**: 当前系统 (`kernel/engine.go`) 仅获取“持仓币种”和“候选币种”的行情数据。如果不交易 BTC，且 BTC 不在推荐列表中，**Context 中将完全没有 BTC 的 K 线数据**。
-   **影响 (Impact)**: 无法计算 BTC 的涨跌幅 (`btc_pnl_pct`)，也就无法计算“相对强弱” (`relative_strength`)，策略无法落地。

### 逻辑可行性 (Logic Feasibility)
-   **可行 (Feasible)**: 用户提出的 7 种气候状态 (从[狂暴牛市]到[良性洗盘]) 逻辑清晰，完全可以通过 Go 代码实现。
-   **优越性 (Advantage)**: 相比于让 AI 去看图，**代码级预计算**具有压倒性优势：
    -   **精准**: 严格执行 `0.3%` 阈值，无幻觉。
    -   **省 Token**: 直接告诉 AI "当前是[吸血独涨]气候"，AI 无需消耗 Token 去推理。
    -   **全局一致**: 所有币种共享同一个“天气预报”。

## 2. 实施方案 (Implementation Plan)

### 第一步：数据源增强 (Global Market Monitor)
修改 `kernel/engine.go` 中的 `fetchMarketDataWithStrategy`：
-   **强制获取**: 无论当前交易什么币，**强制获取 BTCUSDT 的 15m K线数据**。
-   **结构体扩展**: 在 `Context` 中新增 `GlobalMarket` 字段。

```go
type GlobalMarket struct {
    BTCPriceChange   float64 // BTC 15m 涨跌幅
    MarketRegime     string  // 气候状态 (e.g., "Raging Bull")
    RegimeReason     string  // 状态判定原因
    BTCVolumeSpike   bool    // 是否放量 (用于判定恶性暴跌)
}
```

### 第二步：气候判定引擎 (Climate Engine)
在 `kernel/strategy_climate.go` (新建) 中实现 7 态判定逻辑：

1.  **计算基准**:
    -   `btc_change`: BTC 15m 涨跌幅
    -   `symbol_change`: 当前币种 15m 涨跌幅
2.  **状态机**:
    -   `if Abs(btc_change) < 0.003`: **[静区]** -> 判定 [资金溢出]
    -   `if btc_change >= 0.003`: **[异动涨]** -> 对比 `symbol_change` 判定 [狂暴] 或 [吸血]
    -   `if btc_change <= -0.003`: **[异动跌]** -> 对比 `symbol_change` 判定 [阴跌]
    -   `if btc_change <= -0.02`: **[暴跌]** -> 判定 [恶性] 或 [良性]

### 第三步：Prompt 注入 (Context Injection)
在 `kernel/engine.go` 构建 Prompt 时，将“天气预报”置顶：

```text
## 🌍 Global Command (Market Regime)
Current Climate: [Vampire Market] (BTC +0.8%, Symbol -0.2%)
Strategy: DEFENSIVE (Capital is flowing into BTC, Alts are draining)
```

### 第四步：策略文件适配 (Strategy Adaptation)
修改 `fsdownload/RightSide_Flow_Hunter_V1.0.json` 的 Prompt：
-   **移除** 原有的 "Step 1 全局指挥" 中关于暴跌/震荡的手动计算逻辑（交给 Go 代码算）。
-   **新增** 对 `Global Market Regime` 的响应规则。
    -   Example: `If MarketRegime == "Vampire Market" -> Action: Reduce Risk / Wait`

## 3. 结论 (Conclusion)
该策略**高度可行且必要**。它填补了系统“只见树木不见森林”的缺陷。
建议立即实施开发，预计涉及 3 个文件的修改。
