# 全局指挥配置化 (Global Command Configurability)

## 目标
允许用户通过 Web 前端修改“全局指挥”策略的判定阈值，并支持添加新的市场气候类型。

## 核心设计
采用**基于规则 (Rule-based)** 的系统替代硬编码的逻辑。每个气候类型被定义为一个规则，包含条件 (Condition) 和结果 (Result)。规则按顺序评估，第一个匹配的规则决定当前的市场气候。

## 提议的变更

### 后端

#### [MODIFY] [store/strategy.go](file:///d:/code/nofx-0210/store/strategy.go)
- 添加 `GlobalCommandConfig` 结构体：
  ```go
  type GlobalCommandConfig struct {
      Enabled bool `json:"enabled"`
      Rules   []ClimateRule `json:"rules"`
      DefaultClimate ClimateResult `json:"default_climate"`
  }
  
  type ClimateRule struct {
      Name string `json:"name"`
      Conditions []Condition `json:"conditions"`
      Result ClimateResult `json:"result"`
  }
  
  type Condition struct {
      Metric string `json:"metric"` // "btc_change", "alt_change"
      Operator string `json:"operator"` // ">", "<", ">=", "<=", "between"
      Value float64 `json:"value"`
      Value2 float64 `json:"value2,omitempty"` // for "between"
  }
  
  type ClimateResult struct {
      Climate string `json:"climate"`
      Reason string `json:"reason"`
      Instruction string `json:"instruction,omitempty"` // Custom instruction template
  }
  ```
- 更新 `StrategyConfig` 以包含 `GlobalCommand *GlobalCommandConfig` (替代之前的 simple map check, or keep map for enable/disable and use this for detailed settings).
  - *Decision*: Keep `QuantStrategies map[string]bool` for quick toggles, adds `GlobalCommandConfig` for settings.

#### [MODIFY] [kernel/quant/global_command.go](file:///d:/code/nofx-0210/kernel/quant/global_command.go)
- 重构 `AnalyzeGlobalMarket` 接受 `*store.GlobalCommandConfig`。
- 实现规则评估逻辑 `evaluateRules(rules, btcChange, altChange)`.
- 实现 `GetDefaultGlobalCommandConfig()` 返回默认的硬编码逻辑对应的规则集。

#### [MODIFY] [kernel/engine.go](file:///d:/code/nofx-0210/kernel/engine.go)
- 在调用 `quant.AnalyzeGlobalMarket` 时传递配置。

### 前端

#### [MODIFY] [web/src/types.ts](file:///d:/code/nofx-0210/web/src/types.ts)
- 更新 TS 定义以匹配新的后端结构。

#### [NEW] [web/src/components/GlobalCommandEditor.tsx](file:///d:/code/nofx-0210/web/src/components/GlobalCommandEditor.tsx)
- 创建规则编辑器组件。
- 支持：
  - 规则列表展示（支持拖拽排序或上移/下移）。
  - 编辑规则条件（BTC涨跌幅、山寨币涨跌幅）。
  - 编辑规则结果（气候名称、AI指令）。
  - 添加/删除规则。

#### [MODIFY] [web/src/pages/QuantStrategiesPage.tsx](file:///d:/code/nofx-0210/web/src/pages/QuantStrategiesPage.tsx)
- 在 Global Command 卡片上添加 "Configure" (配置) 按钮。
- 点击按钮打开模态框或抽屉，显示 `GlobalCommandEditor`。

## 验证计划

1.  **单元测试**: 针对 `evaluateRules` 编写测试，确保其逻辑与旧版硬编码逻辑一致。
2.  **前端交互**:
    - 验证可以修改阈值（例如将牛市阈值从 1% 改为 2%）。
    - 验证可以添加新规则（例如 "Super Crash"）。
    - 验证保存后，配置正确持久化到后端。
3.  **集成测试**:
    - 修改配置后，触发一次 Run Test (在 API page 或 logs)，确认 AI 接收到的 Prompt 包含新的气候名称/指令。
