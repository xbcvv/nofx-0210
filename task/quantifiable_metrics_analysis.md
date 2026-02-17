# Optimized Quantifiable Analysis (Decoupled Architecture)

> **核心思想**: **数据原语 (Data Primitives) + 动态逻辑 (Dynamic Logic)**
> **目标**: 后端 Go 代码只负责提供“乐高积木”（标准数学属性），前端 Prompt/Config 负责定义“拼搭规则”（业务逻辑）。
> **优势**: 调整策略参数（如均线周期、判定阈值）无需修改代码。

## 1. 趋势与形态 (Trend & Pattern)

### 1.1 均线属性 (EMA Attributes)
- **解耦方案**:
  - **Go 端**: 读取 `strategy.json` 中的 `ema_periods` (e.g. `[20, 60, 120]`)，为每条均线计算并输出对象。
  - **原语 (Primitives)**:
    - `Value`: 均线数值
    - `Angle`: 归一化斜率 (Normalized Slope)
    - `Spread`: 与价格的乖离率 `(Price - EMA) / Price`
- **Prompt 逻辑 (示例)**:
  - *"Check if EMA20 Angle > 15 (Steep Slope)"*
  - *"Check if EMA60 Angle > 5 (Long Term Bull)"*
  - *"Check if EMA_Spread < 0.1% (Entangled)"*

### 1.2 动能属性 (MACD Attributes)
- **解耦方案**:
  - **Go 端**: 标准输出 `DIF`, `DEA`, `Hist`。
  - **原语 (Primitives)**:
    - `Hist_Value`: 柱状图数值
    - `Hist_Slope`: 柱状图变化率 `(Hist_t - Hist_{t-1})`
    - `Cross_Distance`: 距离上次金叉/死叉的 K 线数
- **Prompt 逻辑 (示例)**:
  - *"Check if Hist_Slope > 0 AND Hist_Value > 0 (Momentum Expansion)"*
  - *"Check if Hist_Value < 0 AND Hist_Slope > 0 (Rebound within Trend)"*

## 2. 回归与反转 (Mean Reversion)

### 2.1 橡皮筋指数 (Rubber Band Index)
- **解耦方案**:
  - **Go 端**: 不判断“是否触发”，只提供偏离度数据。
  - **原语 (Primitives)**:
    - `Z_Score`: `(Price - EMA20) / ATR(14)`
    - `RSI_Value`: 标准 RSI(14) 数值
- **Prompt 逻辑 (示例)**:
  - *"If Z_Score > 2.5 AND RSI > 80 (Overbought Reversion)"*
  - *"If Z_Score < -3.0 (Extreme Panic)"*

## 3. 市场状态 (Market Regime)

### 3.1 波动率状态 (Volatility State)
- **解耦方案**:
  - **Go 端**: **移除** `IsStagnant` / `IsDead` 等布尔判断。只提供基础波动数据。
  - **原语 (Primitives)**:
    - `ADX_Value`: 趋势强度
    - `Change_Abs_1h`: 1小时绝对涨跌幅
    - `Volume_Ratio`: 量比
- **Prompt 逻辑 (示例)**:
  - **死寂 (Stagnancy)**: *"If ADX < 20 AND Change_Abs_1h < 0.3%"*
  - **吸血 (Siphon)**: *"If BTC_ADX > 30 AND Symbol_Change < 0%"*

### 3.2 极端行情 (Extreme Events)
- **解耦方案**:
  - **Go 端**: 提供多周期涨跌幅。
  - **原语 (Primitives)**:
    - `Change_15m`: 15分钟涨跌幅
    - `Change_5m`: (可选) 5分钟涨跌幅
- **Prompt 逻辑 (示例)**:
  - **暴跌 (Crash)**: *"If Change_15m < -2.0%"*
  - **急拉 (Pump)**: *"If Change_5m > 3.0%"*

### 3.3 趋势延续 (Trend Continuation)
- **解耦方案**:
  - **Go 端**: 提供足够长度的历史数据 (DisplayCount)。
  - **原语 (Primitives)**:
    - `Display_Count`: K线显示数量 (默认60，约15小时)
    - `Structure_High/Low`: 日内结构高低点
- **Prompt 逻辑 (示例)**:
  - **结构破坏**: *"Check if Price breaks the highest point of the last 60 bars"*

## 4. 量能分析 (Volume Analysis)

### 4.1 资金流向 (Flow Primitives)
- **解耦方案**:
  - **Go 端**: 提供 OI 和价格变化的原语。
  - **原语 (Primitives)**:
    - `OI_Change_1h`: 1小时 OI 变化量
    - `Price_Change_1h`: 1小时价格变化量
- **Prompt 逻辑 (示例)**:
  - **真流入 (Inflow)**: *"If OI_Change > 0 AND Price_Change > 0"*
  - **离场 (Outflow)**: *"If OI_Change < 0"*

### 4.2 量价背离 (Volume Divergence)
- **解耦方案**:
  - **Go 端**: 提供量比。
  - **原语 (Primitives)**:
    - `Vol_MA_Ratio`: `Volume / MA(Vol, 20)`
- **Prompt 逻辑 (示例)**:
  - **缩量突破 (Fake Breakout)**: *"If Price breaks High AND Vol_MA_Ratio < 0.8"*

## 5. 架构总结 (Architectural Summary)

| 层级 (Layer) | 职责 (Responsibility) | 示例 (Examples) | 修改频率 (Frequency) |
| :--- | :--- | :--- | :--- |
| **Go Backend** | **计算器 (Calculator)** | `CalculateEMA`, `CalculateADX`, `GetPriceChange` | 极低 (Stable) |
| **Strategy Config** | **参数配置 (Parameters)** | `ema_periods: [20, 60]`, `rsi_period: 14` | 中 (Configurable) |
| **System Prompt** | **逻辑判定 (Judge)** | *"If ADX < 20 then Dead"*, *"If Slope > 10 then Bull"* | 高 (Flexible) |

---
**下一步建议**:
1.  **Refactor**: 将后端 `Indicator` 模块重构为通用计算服务。
2.  **Expose**: 修改 `Formatter`，将计算出的 `Primitives` 结构化注入 Prompt。
3.  **Config**: 在 `strategy.json` 中开放更多参数配置数组。
