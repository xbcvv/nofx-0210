# Technical Indicators Guide

This document lists all technical indicators available in the trading system and their usage.

## Trend Indicators

### ADX (Average Directional Index)
- **Description**: Measures the strength of a trend, regardless of direction.
- **Components**:
    - **ADX**: The main line. Values > 25 indicate a strong trend.
    - **DI+ / DI-**: Directional indicators. DI+ > DI- suggests bullish trend, DI- > DI+ suggests bearish trend.
- **Usage**: Used to filter out ranging markets. The system uses ADX > 25 as a confirmation for trend-following strategies.

### EMA (Exponential Moving Average)
- **Description**: Weighted moving average that gives more importance to recent price data.
- **Usage**:
    - **EMA20**: Short-term trend baseline.
    - **EMA50**: Medium-term trend baseline.
    - **Crossover**: EMA20 crossing above EMA50 is a bullish signal (Golden Cross).

### MACD (Moving Average Convergence Divergence)
- **Description**: Trend-following momentum indicator.
- **Usage**:
    - **Zero Line**: MACD > 0 suggests bullish momentum.
    - **Signal Line**: MACD line crossing signal line indicates momentum shifts.

## Momentum / Volatility Indicators

### RSI (Relative Strength Index)
- **Description**: Measures the speed and change of price movements.
- **Usage**:
    - **Overbought**: RSI > 70.
    - **Oversold**: RSI < 30.
    - **Divergence**: Price making new high but RSI failing to do so suggests reversal.

### ATR (Average True Range)
- **Description**: Measures market volatility.
- **Usage**:
    - **Stop Loss**: Used to calculate dynamic stop loss distances (e.g., 2 * ATR).
    - **Position Sizing**: Higher volatility (ATR) leads to smaller position sizes.

### Bollinger Bands (BOLL)
- **Description**: Volatility bands placed above and below a moving average.
- **Usage**:
    - **Squeeze**: Bands contracting indicates low volatility and potential breakout.
    - **Mean Reversion**: Price touching outer bands often reverts to the mean (middle band).

## Market Data / Volume

### Volume
- **Description**: Trading volume.
- **Usage**: Confirms trend strength. Rising volume on uptrend validates the move.

### Open Interest (OI)
- **Description**: Total number of outstanding derivative contracts.
- **Usage**: Rising OI during a trend indicates new money entering the market (strong trend). Declining OI suggests liquidation or exiting (trend weakening).

### Funding Rate
- **Description**: Periodic payments to traders that are long or short based on the difference between perpetual and spot prices.
- **Usage**: High positive funding rate suggests over-leveraged longs (potential contrarian short signal).

## Quantitative Data (NofxOS)

### Quant Data
- **Price Change**: Performance across multiple timeframes (5m to 24h).
- **Netflow**: Visualization of institutional vs. retail fund flow.
- **OI Ranking**: Ranking of coins by Open Interest change.
