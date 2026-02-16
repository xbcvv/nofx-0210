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
