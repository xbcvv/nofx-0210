> [!CAUTION]
> **DEPRECATED / 废弃说明**
> 此文档描述的“全局指挥 (Global Command)”功能已被废弃，且在当前代码库中**不存在**。
> This feature is deprecated and does NOT exist in the current codebase.
> 请勿根据此文档进行开发或调试。
> Do NOT use this document for development or debugging.

# 演练 - 量化策略与全局指挥

## 1. 全局指挥 (Global Command) 重构

### 目标
将“全局指挥”策略从核心引擎中解耦，防止循环依赖，并便于后续扩展量化策略。

### 变更
- **逻辑迁移**: 将 `AnalyzeGlobalMarket` 从 `kernel/strategy_climate.go` 迁移到新的 `kernel/quant` 包中。
- **解耦**: 新的 `quant.AnalyzeGlobalMarket` 函数现在接受市场数据映射表 (Map)，而不是沉重的 `kernel.Context` 对象。
- **集成**: 更新 `kernel/engine.go` 以使用新的 `quant` 包，并根据配置有条件地注入全局指挥提示词。

## 2. 量化策略页面

### 目标
提供一个用户界面来管理位于个体交易代理之上的高级量化策略（如全局指挥）。

### 特性
- **新页面**: `QuantStrategiesPage` (通过顶部的 "Quant/量化" 导航进入)。
- **开关控制**: 用户可以启用/禁用“全局指挥”集成。
- **视觉效果**: 使用 Tailwind CSS 适配应用的暗色主题，带有状态指示器和提示词预览。
- **持久化**: 连接到 `api.updateStrategy` 以保存配置。
- **国际化**: 全面支持中/英双语界面。

### 截图
*(请参考任务 Artifact 中的截图)*

## 3. 关键修复 (Backend Fix)

### 策略配置热重载
- **问题**: `AutoTrader` 在启动时加载配置，导致 API 更新的策略开关无法即时生效。
- **修复**: 修改 `trader/auto_trader.go`，在每个 `runCycle` 开始时从数据库重新读取最新的活跃策略配置。
- **效果**: 用户在前端页面切换 "Global Command" 开关后，正在运行的交易机器人会在下一个周期立即应用新设置，无需重启。

## 4. 文档 (Documentation)

### 新增/更新文档
- **[NEW] 全局指挥策略手册**: `docs/strategies/GLOBAL_COMMAND.zh-CN.md`
  - 详细解释了市场气候判定逻辑、阈值和 AI 指令生成规则。
- **[UPDATE] 策略模块架构**: `docs/architecture/STRATEGY_MODULE.zh-CN.md`
  - 增加了 2.6 节 "全局指挥 (Global Command)"，说明其在系统架构中的位置。

## 5. 验证

### 后端
- 运行 `go build ./kernel/...` 和 `go build ./trader/...` 确保代码编译通过。**结果：通过。**
- 逻辑代码审计确认了热重载机制的正确性。

### 前端
- 使用 `useSWR` 进行状态管理和 `lucide-react` 图标实现了 `QuantStrategiesPage.tsx`。
- 验证了导航、路由和国际化切换功能。

## 下一步
- 用户现在可以通过 UI 实时配置全局指挥策略。
- 所有的架构变更已文档化，便于后续开发。

## 6. 全局指挥配置化 (Global Command Configurability)

### 功能说明
现在，用户可以直接在前端界面配置“全局指挥”策略的判定规则，而无需修改代码。

### 使用指南
1.  **访问页面**: 导航至“量化策略 (Quant Strategies)”页面。
2.  **打开配置**: 点击“全局指挥 (Global Command)”卡片上的“配置 (Configure)”按钮。
3.  **编辑器功能**:
    *   **启用/禁用**: 顶部开关可一键启停该策略。
    *   **默认气候 (Default Climate)**: 当所有规则都不匹配时，系统将使用的默认判断（兜底）。
    *   **规则管理 (Rules Management)**:
        *   **添加规则**: 点击右上角“Add Rule”。
        *   **编辑规则**: 点击规则卡片展开详细编辑器。可以修改规则名称、添加/删除条件、设置结果气候和AI指令。
        *   **调整优先级**: 使用上/下箭头调整规则评估顺序（系统从上到下匹配，取第一个满足的规则）。
        *   **删除规则**: 点击垃圾桶图标删除规则。
4.  **保存生效**: 点击底部“Save Configuration”保存。**更改即时生效**，无需重启后端服务。

### 技术实现
*   **后端**: `GlobalCommandConfig` 结构存储规则，`kernel/quant` 模块实现了基于规则的判定引擎。
*   **前端**: 新增 `GlobalCommandEditor` 组件，提供可视化规则编辑。
*   **热重载**: 自动交易机器人 (AutoTrader) 会在每个交易周期开始时重新加载策略配置，确保规则修改立即应用。

## 7. 前端汉化 (Frontend Localization)

### 目标
为新增的量化策略页面和全局指挥编辑器提供完整的中文支持。

### 变更
- **翻译文件**: 更新 `web/src/i18n/translations.ts`，新增了 `quantStrategies`、`configureGlobalCommand`、气候判定规则相关的中英双语键值。
- **页面适配**: 
  - 重构 `QuantStrategiesPage.tsx`，将硬编码的英文替换为 `t()` 函数调用。
  - 重构 `GlobalCommandEditor.tsx`，实现编辑器内所有标签、按钮、下拉选项的国际化。
- **一致性**: 解决了 `active` 等通用键值的命名冲突，统一了术语（如 "Global Command" -> "全局指挥"）。

### 效果
- 用户在切换语言为中文时，策略配置界面将完全显示为中文，包括复杂的规则条件和操作提示。

## 8. Docker 构建流程调整

### 目标
管理多架构 Docker 镜像的构建策略，根据需要启用或禁用 ARM64 构建。

### 变更记录
1.  **临时禁用 ARM64 (已回滚)**: 
    - 曾临时注释掉 `linux/arm64` 平台配置以加速构建。
2.  **恢复多架构构建 (当前状态)**:
    - 恢复了 `.github/workflows/docker-build.yml` 中的 `linux/arm64` 平台配置。
    - 重新启用了 `create-manifest` 任务。

### 影响
- CI 流程将构建并推送 `linux/amd64` 和 `linux/arm64` 两种架构的镜像。
- 最终生成的 Docker Manifest list 将包含两种架构，确保用户在 x86 服务器和 ARM 设备（如 Raspberry Pi, Apple Silicon Mac）上均可无缝拉取运行。
