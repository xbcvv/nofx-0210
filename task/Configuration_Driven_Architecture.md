# 配置驱动架构 (Configuration-Driven Architecture) 修改说明

> **版本**: V1.1.0
> **日期**: 2026-02-17
> **目标**: 实现零代码策略迭代。通过 JSON 配置文件驱动后端生成数据，前端 Prompt 直接调用。

## 1. 修改摘要 (Modification Summary)

### 1.1 核心数据结构升级 (Data Structure)
- **文件**: `market/types.go`
- **新增字段**:
  - `PriceChanges`: 动态存储任意时间周期 (1d, 15m, 5m...) 的涨跌幅。
  - `DynamicEMAs`: 动态存储任意周期 (20, 60, 120...) 的 EMA 数据 (含斜率、乖离率)。
  - `DynamicATRs`: 动态存储任意周期 (7, 14...) 的 ATR 数据。

### 1.2 计算引擎泛化 (Calculation Engine)
- **文件**: `market/data.go`
- **功能增强**:
  - 重构 `GetWithTimeframes` 函数，使其不再硬编码周期，而是接受 `emaPeriods` 和 `atrPeriods` 数组。
  - 实现了 `calculateEMASlope` 算法，为每条均线计算归一化斜率。
  - 实现了动态 `Change_{tf}` 计算，支持 `selected_timeframes` 中的所有周期。

### 1.3 数据注入与Schema (Injection & Schema)
- **文件**: `kernel/formatter.go`, `kernel/engine.go`, `kernel/schema.go`
- **功能增强**:
  - `Schema` 明确定义了双语对照 (Bilingual definitions)。
    - `Change_{tf}` 对应 "周期涨跌幅" (e.g. "日涨跌幅").
    - `EMA{period}` 对应 "指数均线" (Value=点位, Slope=斜率, Spread=乖离).
    - `ATR{period}` 对应 "平均真实波幅".
    - `ATR{period}` 对应 "平均真实波幅".
    - 确保 AI 能完美理解中文 Prompt (如 "EMA60斜率 > 10")。

### 1.4 K线显示优化 (K-Line Display Optimization)
- **文件**: `store/strategy.go`, `kernel/formatter.go`
- **功能增强**:
  - `klines` 配置中新增 `display_count` 参数 (默认 30)。
  - 允许将显示给 AI 的 K 线数量提升至 60-96 根，以覆盖更完整的日内市场结构 (Day Structure)，解决 AI 视野过短的问题。
  - 移除了代码中硬编码的 30 根限制。

---

## 2. 使用方法 (Usage Guide)

### 2.1 如何配置 `strategy.json`

您只需要修改 `RightSide_Flow_Hunter_V1.0.json` (或您的策略配置文件) 中的数组即可。

```json
{
  "config": {
    "indicators": {
      // [1] 定义 K 线与涨跌幅周期
      "klines": {
        "selected_timeframes": [
          "15m", 
          "4h", 
          "1d"   // <--- 添加 "1d"，后端自动生成 Change_1d
        ]
      },

      // [2] 定义技术指标周期
      "ema_periods": [
        20,
        50,
        60,    // <--- 添加 60，自动生成 EMA60 (含斜率/乖离)
        120
      ],

      "atr_periods": [
        7,     // <--- 添加 7，自动生成 ATR7
        14
      ]
    },
    // [3] K线显示数量配置
    // 默认为 30，建议设为 60+ 以覆盖日内结构
    "display_count": 60
  }
}
```

### 2.2 如何在 Prompt 中调用

由于 `kernel/schema.go` 已经定义了通用术语，AI 现在能理解 `Change_{周期}` 和 `EMA{周期}`。

#### 场景 1: 每日复盘 (使用 1d 涨跌幅)
```yaml
气候判定:
  - [日线暴跌]: 
    # 只要您在 json 里配了 "1d"，这里就能用 Change_1d
    条件: Change_1d < -5.0%
    策略: 【防守】
```

#### 场景 2: 橡皮筋反转 (使用 EMA60 + ATR7)
```yaml
特快通道:
  - 橡皮筋_V2:
    # 只要您在 json 里配了 60 和 7，这里就能用
    公式: (Price - EMA60) / ATR7
    条件: abs(公式) > 2.5
```

#### 场景 3: 趋势斜率 (使用 EMA120 斜率)
```yaml
战略 (1D):
  # 系统会自动输出 EMA120 (Slope: x.x)
  - 强趋势: EMA120.Slope > 10
```

---

## 3. 注意事项 (Notes)

1.  **重启生效**: 修改 `json` 配置文件后，必须重启程序才能生效。
2.  **Schema 定义**: 系统 Prompt 头部已包含如下定义，无需您手动解释：
    > Change_{tf}: Price change percent over timeframe
    > EMA{period}: Exponential Moving Average (Value, Slope, Spread)
3.  **性能损耗**: 理论上添加过多的周期会增加计算量，但对于几十个周期以内，Golang 的计算速度是毫秒级的，可忽略不计。
