# 交易所特性差异 (Exchange Features Matrix)

在 NOFX 多交易所集成中，不同交易所的底层 API 机制存在客观差异，这导致了部分高级交易功能在某些交易所原生支持，而在另一些交易所受到限制。

## 功能支持矩阵

| 功能项 | Binance (币安) | Bitget | 备注说明 |
| :--- | :--- | :--- | :--- |
| **基础开平仓** | ✅ 支持 | ✅ 支持 | 市价单全量开平仓均完美支持 |
| **全量止损止盈 (Full TP/SL)** | ✅ 支持 | ✅ 支持 | 开仓时或持仓中附带 TP/SL 价格 |
| **部分平仓 (Partial Close)** | ✅ 支持 | ✅ 支持 | 均可通过给 `ClosePosition` 传递具体 `quantity` 数量并附加 `reduceOnly` (只减仓) 实现局部减仓 |
| **部分止损止盈 (Partial TP/SL)** | ❌ **不受支持** | ✅ **支持** | **区别标记**：Binance 默认的止盈止损逻辑通常绑定整个仓位；而 **Bitget** 计划委托原生支持 `size` 参数，允许单独为仓位的一部分设定退出计划。 |

## Bitget 特有参数说明 [Bitget-Specific]

针对 Bitget 的 TP/SL 设置，底层系统使用了以下特有参数来实现和 Binance 的行为隔离及高级功能支持：

1. **`planType` (计划委托类型与智能降级)**
   - Bitget V2 API 要求，对**已有持仓**设置触发式平仓单时，需尽可能匹配持仓使用局部平盘 (`pos_profit` 和 `pos_loss`)。
   - **智能降级 (Fallback)**：由于开仓瞬间存在微小的异步延迟，底层系统（`trader.go`）会在首选 `pos_loss` 遭遇 `400172 Illegal type` 拒绝时，毫秒级自动**无缝降级**为原生的无持仓约束条件单 (`loss_plan` / `profit_plan`)。从而保证在任何震荡市和物理延迟下，止盈止损指令 **100% 挂载成功**。
   - 注意：V2 API 已废除使用 `normal_plan` 等闲杂类型挂载止盈止损单，违规会触发 `43011 delegateType error`。
2. **`size` (局部止盈止损数量)**
   - 区别于 Binance 的全仓绑定，Bitget `pos_profit` 和 `pos_loss` API **原生支持接收 `size` 参数**。
   - 这意味着，当你调用系统或 AI 下达带有明确 `quantity` 的止盈止损指令时，系统能够调用底层 `SetTakeProfit` 和 `SetStopLoss` 把 `quantity` 发送给 Bitget。
   - **效果**：实现阶梯止盈、分批止损等高级交易策略。币安在统一账户架构下目前较难做到一次性提交原生的分批止盈止损订单。
