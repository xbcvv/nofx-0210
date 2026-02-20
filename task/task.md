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
    - [x] 二次修订赋予了资深交易员推演权，但输出了大量罗嗦步骤和Markdown列表。
    - [x] **三次修订（终极版）**：彻底封杀 Markdown 列表和标题排版，将 `<reasoning>` 强制压缩为仅含四行纯文本的“系统日志”结构（气候、持仓、扫描、反思），截断 Claude 的思维链。
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
    - [x] 二次修订赋予了资深交易员推演权，但输出了大量罗嗦步骤和Markdown列表。
    - [x] **三次修订（终极版）**：彻底封杀 Markdown 列表和标题排版，将 `<reasoning>` 强制压缩为仅含四行纯文本的“系统日志”结构（气候、持仓、扫描、反思），截断 Claude 的思维链。
- [x] **Bitget 阶梯止盈优化**:
    - [x] 利用 Bitget 原生支持部分减仓特性，修改 `decision_process` 中关于利润管理的断言。
    - [x] 将全平规则改为：浮盈 > 1.5R 时强制使用 `partial_close 0.2~0.3` 进行阶梯落袋，并推保护损，保留20%底仓让利润奔跑。
- [x] **隐私保护**:
    - [x] `fsdownload` 已通过 `.gitignore` 排除，仅在本地更新策略文件，不上传 GitHub。
- [x] **Bitget开仓挂单时序冲突修复**:
    - [x] 重构底层 API 方法。调用止盈止损时，系统会动态侦测是否具有匹配的现有底仓。匹配则采用 `pos_loss`/`pos_profit` (局部止损/止盈)，否则降级尝试新挂条件单。如被拦截则自动双向进行 fallback 降级补偿。
    - [x] **终极修复**：追踪错误日志确认 Bitget V2 的非持仓条件单依然需要严格的 `loss_plan`/`profit_plan` 枚举，之前被 `pos_loss` 拦截误导而尝试的 `normal_plan` 导致了 `delegateType error`。目前已将底层 fallback 的标准参数替换为官方合法且不受持仓状态约束的 `loss_plan` 和 `profit_plan`，彻底杜绝 `Illegal planType` 错误。
## 修复 AutoTrader 开仓保证金不足错误 (40762) [Fix]
- [x] **问题分析**:
    - [x] 在 open_long 最新日志中出现 40762 The order amount exceeds the balance 错误。
    - [x] 排查认定：引擎在根据 Available Balance 计算最大可购买的头寸价值时，原设定只预留了极小的 2% 价格缓冲下限。
    - [x] 遭遇极端情况组合：1. 实时行情取自第三方 (Binance Fallback) 导致取值有微小偏差，在除法计算时导致对 Bitget 生成了偏大的 quantity。2. Bitget 针对市价单开仓会默认冻结比限价单多达数个百分点的冗余保证金防滑点补偿。二者结合彻底击穿了 2% 的安全边界，导致资金不足而拒单。
- [x] **方案实施**:
    - [x] 修改 	rader/auto_trader.go 中的 maxAffordablePositionSize 缓冲调节机制。
    - [x] 将原本过窄的  .98 (2% 缓冲) 直接拓宽至更为宽容稳健的  .90 (10% 缓冲)，这能彻底吸纳无论多大的跨所差价水分以及市价单巨额锁仓。编译并同步至 GitHub，此后再无越权爆表隐患。
- [x] **修复 AI 引擎无法感知手动修改止盈止损问题 (Bitget Sync)**:
    - [x] **问题根源**：系统的 AI Prompt 依赖 GetPositions 拉取持仓信息，但交易所的仓位接口原生不返回独立的条件委托单 (StopLoss / TakeProfit)。这导致 AI 一直在盲人摸象，它仅凭自己的内部“记忆寄存器” (Memory Bank) 记得上一次设定的 SL/TP，每次 Hold 时都盲目覆盖重发，导致用户在 App 端的修改被冲刷。
    - [x] **修复方案**：重构底层 	rader/auto_trader.go。在给 AI 装配持仓列表之前，强行插入一步 GetOpenOrders 二次查询。抓取当前所有活跃的 STOP_MARKET 和 TAKE_PROFIT_MARKET 条件单，并与各个仓位一一对应。
    - [x] **AI 赋能**：现在，在 AI 的系统中，Current Positions 字符串将带有完整的 | SL 0.35 | TP 0.40 实时尾巴，它终于能看到您手动修改的止损点了！并基于最新的真实风险进行推演。
- [x] **修复 Bitget 最新 V2 环境下止盈止损一直挂单失败 (400172 Illegal Type) 问题**:
    - [x] **问题根源**：之前的代码在给当前开仓部位附加止盈止损时，使用了普通的计划交易终端 /api/v2/mix/order/place-plan-order，但强行塞入了 pos_loss (仓位止损) 这个特权 planType。而自 Bitget V2 以来，普通计划终端强制仅接受 
ormal_plan，直接抛出了 400172 参数结构非法的验证错误。
    - [x] **修复方案**：将 SetStopLoss 和 SetTakeProfit 发送请求的目标 URL 彻底更正为专门用于持仓管理条件的特权端点：/api/v2/mix/order/place-pos-tpsl。
    - [x] **参数映射**：废弃了笼统的 	riggerPrice 字段，严格对标 V2 接口协议，将止损字段映射为 stopLossTriggerPrice，将止盈字段映射为 stopSurplusTriggerPrice。同时显式将执行价格设定为   (市价强平防滑保护)。现在系统在开仓后附加的防护罩将100%立刻挂载。
