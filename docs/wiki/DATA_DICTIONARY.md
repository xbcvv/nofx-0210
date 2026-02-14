# 📖 交易数据字典与规则 (Data Dictionary & Rules)

本文档是系统 System Prompt 中数据字典的中英文对照参考，旨在帮助用户理解精简后的英文定义。

## 1. 字段定义 (Field Definitions)

### 账户指标 (Account)
| 英文 Key | 中文含义 | 详细说明 |
| :--- | :--- | :--- |
| **Equity** | 总权益 | 账户净值，包含未实现盈亏 (Available + Unrealized PnL) |
| **Balance** | 可用余额 | 可用于开仓的资金 (不含已用保证金) |
| **Margin** | 保证金使用率 | 已用保证金占总权益的比例 (Used Margin / Equity) |
| **PnL** | 总收益率 | 账户自启动以来的总投资回报率 |

### 交易指标 (Trade)
| 英文 Key | 中文含义 | 详细说明 |
| :--- | :--- | :--- |
| **Entry** | 进场均价 | 开仓的平均价格 |
| **Exit** | 出场均价 | 平仓的平均价格 |
| **Profit** | 已实现盈亏 | 该笔交易的实际盈利金额 (USDT) |
| **HoldDuration** | 持仓时长 | 从开仓到平仓经历的时间 |

### 持仓指标 (Position)
| 英文 Key | 中文含义 | 详细说明 |
| :--- | :--- | :--- |
| **UnrealizedPnL** | 未实现盈亏 | 当前持仓的浮动盈亏 |
| **PeakPnL** | 历史峰值盈亏 | 该持仓曾达到的最高浮盈% (用于移动止盈判断) |
| **Drawdown** | 利润回撤 | 当前盈亏距离峰值盈亏的回撤幅度 (Current - Peak) |
| **LiqPrice** | 强平价格 | 预估爆仓价格 (0 表示无风险) |
| **Action** | 交易动作 | `open_long/short` (开多/空), `close_long/short` (平多/空), `partial_close` (分批平), `hold` (持有) |
| **ClosePercentage** | 平仓比例 | 仅配合 `partial_close` 使用，0.5 代表平掉 50% 仓位 |

### 寄存器/记忆 (Register)
| 英文 Key | 中文含义 | 详细说明 |
| :--- | :--- | :--- |
| **Cycle** | 决策周期 | 决策序列号，用于追踪连续性 |
| **MarketRegime** | 市场状态 | AI 对当前市场环境的定性 (如：牛市、震荡、熊市) |
| **ExecutionStatus**| 执行状态 | 上一周期的决策是否执行成功 |

## 2. 核心规则 (Core Rules)

### OI (持仓量) 逻辑
| 现象 | 含义 | 市场解读 |
| :--- | :--- | :--- |
| **OI↑ + Price↑** | Bull Inflow | **多头资金流入**：看涨，趋势增强 |
| **OI↑ + Price↓** | Bear Inflow | **空头资金流入**：看跌，抛压增强 |
| **OI↓ + Price↑** | Shorts Covering | **空头回补**：空头止损离场，可能反转或反弹 |
| **OI↓ + Price↓** | Longs Liquidation | **多头平仓**：多头止损离场，可能反转或回调 |

### 风险控制 (Risk Management)
*   **Max Margin (最大仓位)**: 30% (保留 70% 资金应对风险)
*   **Stop Loss (硬止损)**: -5% (单笔亏损不超过 5%)
*   **Daily Loss (日亏损限额)**: -10% (触及则停止当日交易)

### 策略原则 (Strategy)
*   **Scale In**: 只在盈利时加仓 (Only add to winners)，绝不摊平亏损。
*   **Trailing Stop (移动止盈)**: 当浮动盈亏从 **PeakPnL** 回撤超过 **30%** 时，触发止盈平仓，锁定利润。
