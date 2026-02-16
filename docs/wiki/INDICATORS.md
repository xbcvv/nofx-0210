# 技术指标指南 (Technical Indicators Guide)

本文档列出了交易系统中所有可用的技术指标，包括其配置参数、默认值以及在 AI Prompt 中的呈现格式。

## 全局市场背景 (Global Market Context)
- **描述**: 在 Prompt 顶部强制展示 BTCUSDT 的关键数据，用于判断整体市场气候（"全局指挥"）。即使 BTC 不在候选交易列表中，该数据也会被包含。
- **包含数据**:
    - **价格**: 最新成交价
    - **涨跌幅 (Change)**: 15m, 1h, 4h
    - **ADX**: 趋势强度
- **Prompt 呈现格式**:
    ```text
    ## Global Market Context (Reference Only, NOT for Trading)
    
    **BTCUSDT**:
    - Price: 68000.00
    - Change: 15m -0.50% | 1h +1.20% | 4h +3.50%
    - ADX: 45.2
    ```

## 趋势指标 (Trend Indicators)

### ADX (平均方向指数)
- **描述**: 衡量趋势的强度，无论方向如何。
- **配置参数**:
    - **开关**: `"enable_adx": true`
# ... (rest of the file) ...

### Quant Data
- **配置参数**:
    - **开关**: `"enable_quant_data": true` (默认关闭)
- **Prompt 呈现格式**:
    ```text
    Price Change: 15m: +0.5% | 1h: +1.2% | 4h: +3.5%
    ```
    *注: 15m 涨跌幅 (`PriceChange15m`) 为系统核心字段，用于判断短期剧烈波动（如"恶性暴跌"）。*
    - **周期**: `"adx_periods": [14]` (默认值)
- **Prompt 呈现格式**:
    ```text
    ADX(14): 25.5
    ```
    *注：Prompt 目前仅提供 ADX 主线值，不提供 DI+/DI-。*
- **用法**:
    - ADX > 25 表示趋势强劲，适合趋势跟踪策略。
    - ADX < 20 表示市场无趋势（震荡），建议观望或使用震荡策略。

### EMA (指数移动平均线)
- **描述**: 对近期价格数据赋予更高权重的加权移动平均线。
- **配置参数**:
    - **开关**: `"enable_ema": true`
    - **周期**: `"ema_periods": [20, 50]` (默认值)
- **Prompt 呈现格式**:
    ```text
    EMA20: [65000.1, 65100.2, ...]
    EMA50: [64000.5, 64200.8, ...]
    ```
- **用法**:
    - **趋势判断**: 价格在 EMA20 上方为多头，下方为空头。
    - **交叉信号**: EMA20 上穿 EMA50 (金叉) 看涨，下穿 (死叉) 看跌。

### MACD (异同移动平均线)
- **描述**: 趋势跟踪动量指标。
- **配置参数**:
    - **开关**: `"enable_macd": true`
    - **参数**: 固定为 (12, 26, 9)
- **Prompt 呈现格式**:
    ```text
    MACD: [120.5, 125.3, ...]
    ```
    *注：系统提供的 MACD 为 DIF 线 (快线 - 慢线)，非 histogram。*
- **用法**:
    - MACD > 0 暗示多头势能。
    - MACD 持续上升表示动能增强。

## 动量 / 波动率指标 (Momentum / Volatility Indicators)

### RSI (相对强弱指数)
- **描述**: 衡量价格变动的速度和变化，用于识别超买超卖。
- **配置参数**:
    - **开关**: `"enable_rsi": true`
    - **周期**: `"rsi_periods": [7, 14]` (默认值)
- **Prompt 呈现格式**:
    ```text
    RSI7: [75.2, 78.5, ...]
    RSI14: [65.1, 68.2, ...]
    ```
- **用法**:
    - **超买**: RSI > 70 (可能回调)。
    - **超卖**: RSI < 30 (可能反弹)。
    - **背离**: 价格创新高(低)但 RSI 未创新高(低)，强烈暗示反转。

### ATR (平均真实波幅)
- **描述**: 衡量市场波动率，不指示方向。
- **配置参数**:
    - **开关**: `"enable_atr": true`
    - **周期**: `"atr_periods": [14]` (默认值)
    - **多周期**: 若启用 `"enable_multi_timeframe": true`，将显示更多周期的 ATR。
- **Prompt 呈现格式**:
    ```text
    ATR14: 120.5000
    ```
- **用法**:
    - **动态止损**: 常用 2倍或 3倍 ATR 作为止损距离 (如 `StopLoss = Price - 2.5 * ATR`)。
    - **仓位计算**: 波动率越高，建议仓位越小。

### Bollinger Bands (布林带)
- **描述**: 结合移动平均线和标准差的波动率通道。
- **配置参数**:
    - **开关**: `"enable_boll": true` (默认关闭)
    - **周期**: `"boll_periods": [20]` (默认值，标准差倍数固定为 2)
- **Prompt 呈现格式**:
    ```text
    BOLL Upper: [66000.1, ...]
    BOLL Middle: [65000.5, ...]
    BOLL Lower: [64000.9, ...]
    ```
- **用法**:
    - **收口 (Squeeze)**: 带宽收窄预示即将变盘。
    - **触轨**: 价格触及上轨可能回调，触及下轨可能反弹 (需结合趋势判定)。

## 市场数据 / 资金流 (Market Data / Fund Flow)

### Volume (成交量)
- **配置参数**:
    - **开关**: `"enable_volume": true`
- **Prompt 呈现格式**:
    ```text
    Volume: [100.5M, 120.1M, ...]
    ```

### Open Interest (持仓量 OI)
- **配置参数**:
    - **开关**: `"enable_oi": true`
- **Prompt 呈现格式**:
    ```text
    Open Interest: Latest: 500.2M Average: 480.0M
    ```
- **用法**: 价格上涨 + OI 上涨 = 真突破；价格上涨 + OI 下降 = 空头回补(弱势)。

### Funding Rate (资金费率)
- **配置参数**:
    - **开关**: `"enable_funding_rate": true` (默认关闭)
- **Prompt 呈现格式**:
    ```text
    Funding Rate: 0.0001
    ```

## 量化增强数据 (NofxOS Quantitative Data)

### Quant Data
- **配置参数**:
    - **开关**: `"enable_quant_data": true` (默认关闭)
- **Prompt 呈现格式**:
    ```text
    Price Change: 15m: +0.5% | 1h: +1.2% | 4h: +3.5%
    ```
    *注: 15m 涨跌幅 (`PriceChange15m`) 为系统核心字段，用于判断短期剧烈波动（如"恶性暴跌"）。*

### Quant Netflow (资金净流向)
- **配置参数**:
    - **开关**: `"enable_quant_netflow": true` (默认关闭)
- **Prompt 呈现格式**:
    ```text
    Fund Flow (Netflow):
      Institutional Futures: 5m: +1.2M | 1h: -5.0M ...
    ```
