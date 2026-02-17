# Configuration-Driven Architecture Preview (Aligned with RightSide_Flow_Hunter_V1.0)

> **目标**: 基于系统现有的配置文件结构 (`RightSide_Flow_Hunter_V1.0.json`)，展示如何通过修改参数来驱动后端生成数据，实现零代码策略迭代。

## 1. 策略配置文件 (Based on your JSON)

这是您现有的配置文件结构。您只需在对应的数组中**添加数字或字符串**，后端就会自动为您生产数据。

```json
{
  "name": "RightSide_Flow_Hunter_V1.0",
  "config": {
    "indicators": {
      // 1. K线周期配置
      "klines": {
        "primary_timeframe": "15m",
        "primary_count": 100,
        "enable_multi_timeframe": true,
        // [核心修改点]: 在这里添加 "1d" 或 "5m"
        "selected_timeframes": [
          "15m",
          "4h",
          "1d"   // <--- 添加这个，后端就会生成 Change_1d
        ],
        // [新增]: K线显示数量 (0=默认30, 建议60)
        "display_count": 60
      },

      // 2. 技术指标周期配置 (Array 数组)
      // [核心修改点]: 在数组中添加您想要的周期数字

      // 想要 EMA60 和 EMA120？直接加:
      "ema_periods": [
        20,
        50,
        60,    // <--- 新增
        120    // <--- 新增
      ],

      // 想要 ATR7？直接加:
      "atr_periods": [
        7,     // <--- 新增
        14
      ],

      "rsi_periods": [7, 14],
      "adx_periods": [14],

      // 3. 基础开关 (保持开启)
      "enable_ema": true,
      "enable_atr": true,
      "enable_adx": true,
      "enable_volume": true
    }
  }
}
```

---

## 2. 后端自动生成的原语 (Data Primitives)

根据上述配置，后端会自动计算并向 Prompt 注入以下数据结构。
*(AI 会直接看到这些数据)*

```json
{
  "symbol": "BTCUSDT",
  "price": 68000.00,
  
  // 1. 涨跌幅 (对应 selected_timeframes)
  "changes": {
    "15m": -0.50,
    "4h":  3.50,
    "1d":  -2.10   // <--- 因为配置了 "1d"，所以自动生成
  },

  // 2. EMA 数据 (对应 ema_periods)
  // 包含: 数值(Value), 角度(Slope), 乖离率(Spread)
  "ema": {
    "ema20": { "value": 67500.00, "slope": 15.5, "spread": 0.007 },
    "ema50": { "value": 66500.00, "slope": 8.1,  "spread": 0.022 },
    "ema60": { "value": 66000.00, "slope": 5.2,  "spread": 0.029 }, // <--- 新增
    "ema120":{ "value": 64000.00, "slope": 2.1,  "spread": 0.058 }  // <--- 新增
  },

  // 3. ATR 数据 (对应 atr_periods)
  "atr": {
    "atr7":  150.50, // <--- 新增
    "atr14": 180.20
  },

  // 4. 其他基础数据
  "adx": 45.2,  // 对应 adx_periods[0]
  "rsi": {
      "rsi7": 75.5,
      "rsi14": 65.2
  },
  "volume_ratio": 1.2
}
```

---

## 3. 您的 Prompt 调用方式 (User Logic)

数据有了，您在 `custom_prompt` 或 `prompt_sections` 里直接写逻辑即可。

### 场景 A: 橡皮筋策略 (Z-Score)
*使用新增的 EMA60 和 ATR7*

```yaml
B. 豁免:
  - 橡皮筋_V2: 
    # AI 会自动去 ema.ema60 和 atr.atr7 取值进行计算
    公式: (Price - EMA60) / ATR7
    条件: abs(公式) > 2.5
    执行: >>> 开仓("反转", "橡皮筋V2: 偏离EMA60过大")
```

### 场景 B: 每日复盘 (Change_1d)
*使用新增的 1d 涨跌幅*

```yaml
气候判定:
  - [日线暴跌]: 
    # AI 自动读取 changes.1d
    条件: Change_1d < -5.0%
    策略: 【防守】
```

### 场景 C: 趋势强度 (Slope)
*使用 EMA120 的斜率判断长期趋势*

```yaml
战略 (1D):
  # AI 自动读取 ema.ema120.slope
  - 牛市: EMA120.Slope > 5
  - 熊市: EMA120.Slope < -5
```

---

## 4. 总结

您现有的 `RightSide_Flow_Hunter_V1.0.json` 已经具备了完美的数组结构。我们不需要改变它的格式，只需要：

1.  **在后端 (Go)** 实现一个通用逻辑：遍历这些数组 (`ema_periods`, `selected_timeframes`)，为每个元素计算数据。
2.  **在前端 (JSON)** 往数组里填您想用的数字。

**结果**: 您填什么，系统就算什么。您想用 EMA999？填进去，系统就算 EMA999 给您用。完全不需要改代码。
