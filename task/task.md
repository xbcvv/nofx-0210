# 当前任务追踪 (Current Task Tracking)

## 任务归档与整理 [Archived]
- [x] 创建 `task/` 和 `task/task_archives/` 目录
- [x] 归档旧的实施计划和演练文档
- [x] 创建自动化归档脚本 `scripts/archive_task.ps1`
- [x] 重置 `task.md` 为当前任务状态

## 修复持仓时间归零问题 [Fix]
- [x] **问题分析**:
    - [x] 确认 `kernel/register.go` 和 `trader/binance/futures.go` 的逻辑缺陷
    - [x] 确认 `Runtime` 重置导致的时间归零
    - [x] 确认数据库存储为大写 `LONG/SHORT` 而系统使用小写 `long/short` 导致的读取失败
- [x] **方案制定**:
    - [x] 选择系统级修复方案 (修复数据库读取逻辑)
- [x] **代码实施**:
    - [x] 修改 `store/position.go` 的 `GetOpenPositionBySymbol` 方法
    - [x] 添加 `side = strings.ToUpper(side)` 以兼容大小写
- [x] **验证与审计**:
    - [x] 反向审计修复逻辑 (确保 DB 写入和读取的一致性)
    - [x] 确认代码已正确应用

## 系统增强 [Enhancement]
- [x] **K线数据增强**:
    - [x] 在 `market/types.go` 和 `market/data.go` 中增加 `PriceChange15m` 字段
    - [x] 更新 `kernel/formatter.go` 提示词，明确展示 15m 涨跌幅，防止 AI 幻觉
- [x] **浮盈加仓支持 (Pyramiding)**:
    - [x] 修改 `trader/auto_trader.go`: 允许同方向持仓加仓 (Add-on)
    - [x] 升级风控: `enforcePositionValueRatio` 支持计算总仓位价值 (Total Exposure)
    - [x] 升级止损: `SetStopLoss` 自动应用到总仓位 (Existing + New)

## Prompt 验证 [Verification]
- [x] **prompt23.yaml 分析**:
    - [x] 确认 "全局指挥" 逻辑依赖 BTC 数据
    - [x] 修复 `kernel/engine.go` 强制获取 BTC 数据用于上下文
    - [x] 更新 `kernel/formatter.go` 增加 "Global Market Context" 版块
    - [x] 确认 BTC 在 AI500 列表时可正常交易 (非硬编码排除)

## 功能废弃 [Deprecation]
- [x] **Global Command**:
    - [x] 确认代码库中不存在 Global Command 相关实现
    - [x] 在 `task/task_archives/` 相关文档中添加 `DEPRECATED` 警告


## 量化指标分析 [Analysis]
- [x] **Prompt23.yaml 量化拆解**:
    - [x] 提取所有可量化的逻辑 (如 EMA 斜率, 橡皮筋, 动能背离)
    - [x] 映射为标准数学公式 (Standard Formulas)
    - [x] 输出思维导图文档 `quantifiable_metrics_analysis.md`

## 配置驱动架构实施 [Implementation]
- [x] **数据结构升级**:
    - [x] `market/types.go`: 增加 `EMAs`, `ATRs`, `PriceChanges` 等动态 Map
    - [x] `market/data.go`: 更新 `GetWithTimeframes` 支持动态计算
- [x] **引擎与格式化**:
    - [x] `kernel/engine.go`: 传递 `strategy.json`配置参数
    - [x] `kernel/formatter.go`: 遍历 Map 注入 System Prompt
    - [x] `kernel/schema.go`: 更新 Schema 包含新字段定义
- [ ] **文档**:
    - [x] 创建 `task/Configuration_Driven_Architecture.md` 交付文档

## 动态 Prompt 整合 [Integration]
- [x] **Prompt 逻辑更新**:
    - [x] 修改 `prompt23.yaml` 使用 `PriceChange15m`, `Change_{tf}`, `EMA{period}` 等动态字段
- [ ] **验证**:
    - [ ] 验证 Prompt 逻辑与后端数据的对齐

## K线显示数量优化 [Optimization]
- [x] **配置增强**:
    - [x] `store/strategy.go`: 在 `KlineConfig` 中增加 `DisplayCount` (默认60)
    - [x] `kernel/engine.go`: 将 `StrategyConfig` 注入 `Context`
    - [x] `kernel/formatter.go`: 使用配置的 `DisplayCount` 替代硬编码的 30

## 修复全局锁失效 [Fix]
- [x] **Prompt 缺少持仓时长**:
    - [x] `kernel/formatter.go`: 在 `formatCurrentPositions` 中增加 `Hold Duration` 计算与显示
    - [x] 验证: AI 可见 `⏱️ 持仓时间: 15m`，从而遵守 `< 45m 禁开新仓` 规则

## 文档中文化与索引修复 [Doc]
- [x] **文档中文化**:
    - [x] 移除所有多语言切换导航 (保留纯中文)
    - [x] 将 `docs/i18n/zh-CN/*.md` 提升为主文件，覆盖英文版
    - [x] 修复 `docs/market-regime-classification.md` 等独立文档
- [x] **索引与链接修复**:
    - [x] 更新 `CONTRIBUTING.md` 增加文档维护规则 (GitHub Flow)
    - [x] 更新 `README.md` 索引链接
- [x] **文档完善与修正**:
    - [x] 恢复 `INDICATORS.md` 遗漏细节 (Global Context, ADX, Funding Rate)
    - [x] 修正系统风控规则 (移除硬编码的 Global Lock)
    - [x] 修复 `getting-started/README.md` 缺失的交易所指南链接
    - [x] 修复 `guides/README.md` 缺失的故障排除链接
    - [x] 修复 `wiki/README.md` 缺失的配置驱动文档链接
    - [x] 补充开发者资源 (MCP, Web, Hook) 到主页索引
    - [x] 修复 AI 客户端超时问题 (Context Deadline Exceeded)

## 修复 AI Timeout [Fix]
- [x] **增加默认超时时间**:
    - [x] `mcp/client.go`: 将 `DefaultTimeout` 从 120s 增加到 360s (已回退到 120s，因无效果)

## 修复 BTC ADX 缺失问题 [Fix]
- [x] 分析 `kernel/engine.go` 中关于 BTC 全局环境数据的提取逻辑
- [x] 修复 `BuildUserPrompt` 方法中对 `BTCUSDT` 数据格式化的硬编码问题，将 `btcData.CurrentADX` 补充进发给 AI 的提示词内
- [x] 提交到 GitHub

## 修复 Bitget 止盈止损及订单同步问题 [Fix]
- [x] **止盈止损修复**:
    - [x] 修改 `trader/bitget/trader.go`，将 `SetStopLoss` 的 `planType` 替换为 `pos_loss`
    - [x] 修改 `trader/bitget/trader.go`，将 `SetTakeProfit` 的 `planType` 替换为 `pos_profit`
    - [x] 修改 `CancelStopLossOrders` 和 `CancelTakeProfitOrders`，使用正确的 `pos_loss` 和 `pos_profit`
    - [x] 清理残留的旧常量
- [x] **局部平仓支持**:
    - [x] 代码审查确认：`ClosePosition` 已传递 `quantity` 且使用了 `reduceOnly`，原生支持部分平仓机制
- [x] **局部止盈止损验证**:
    - [x] 代码审查确认：`SetTakeProfit/SetStopLoss` 原生支持传输 `size`，随 `pos_profit` 等计划类型生效
- [x] **订单同步 JSON 解析修复**:
    - [x] 修改 `trader/bitget/order_sync.go` 针对 `{"fillList":null}` 的处理逻辑，防止强制向 `[]BitgetFill` 转换导致的崩溃

## Wiki 与文档更新 [Doc]
- [x] 创建 `docs/wiki/EXCHANGE_FEATURES.md` 记录交易所特性差异矩阵
- [x] 明确标记 Bitget 支持局部止盈止损 (`pos_profit`/`pos_loss` + `size`)，而 Binance 默认不支持局部退出
- [x] 更新 `docs/wiki/README.md` 索引链接

## 交易策略 Prompt 优化 [Enhancement]
- [x] **Claude 模型提效优化**:
    - [x] 在 `custom_prompt` 顶部增加 `<critical_instruction>`。
    - [x] 首次尝试过度限制了 AI 的推演权（被用户拒收）。
    - [x] **二次修订**：赋予了 Claude “资深交易员”的推演权限，但要求风格向 Gemini 看齐，必须严格按照 `# 元规则` 中的精简结构化模板输出 `1. 气候分析 / 2. 持仓分析 / 3. 机会扫描 / 4. 审计反思`，保证有骨架有逻辑但不啰嗦。
- [x] **Bitget 阶梯止盈优化**:
    - [x] 利用 Bitget 原生支持部分减仓特性，修改 `decision_process` 中关于利润管理的断言。
    - [x] 将全平规则改为：浮盈 > 1.5R 时强制使用 `partial_close 0.2~0.3` 进行阶梯落袋，并推保护损，保留20%底仓让利润奔跑。
- [x] **隐私保护**:
    - [x] `fsdownload` 已通过 `.gitignore` 排除，仅在本地更新策略文件，不上传 GitHub。
