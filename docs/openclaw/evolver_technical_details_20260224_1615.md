# Evolver 技术实现详解

## 来源
- **仓库**: https://github.com/EvoMap/evolver
- **文档**: README.zh-CN.md
- **获取时间**: 2026-02-24 16:12 UTC

---

## 🎯 核心定位

**Capability Evolver（能力进化引擎）** 是一个**元技能（Meta-Skill）**，赋予 OpenClaw 智能体**自我反省**的能力。

**核心功能**:
- 扫描自身运行日志
- 识别效率低下或报错的地方
- **自主编写代码补丁**来优化自身性能

---

## 🧬 基因组进化协议（GEP）

**Genome Evolution Protocol** - 将每次进化固化为可复用资产，降低后续同类问题的推理成本。

### 结构化资产目录
```
assets/gep/
├── genes.json      # 基因库
├── capsules.json   # 胶囊库
└── events.jsonl    # 事件日志
```

### Selector 选择器
- 根据日志提取 signals
- 优先复用已有 Gene/Capsule
- 在提示词中输出可审计的 Selector 决策 JSON

---

## 🚀 核心特性

| 特性 | 说明 |
|------|------|
| **自动日志分析** | 扫描 `.jsonl` 会话日志，寻找错误模式 |
| **自我修复** | 检测运行时崩溃并编写修复补丁 |
| **GEP 协议** | 标准化进化流程与可复用资产 |
| **突变协议与人格进化** | 每次进化必须显式声明 Mutation，维护 PersonalityState |
| **可配置进化策略** | 通过 `EVOLVE_STRATEGY` 环境变量选择模式 |
| **信号去重** | 检测修复循环，防止反复修同一个问题 |
| **源码保护** | 防止自治代理覆写核心进化引擎源码 |
| **持续循环模式** | 持续运行的自我进化循环 |

---

## ⚙️ 使用方法

### 标准运行（自动化）
```bash
node index.js
```

### 审查模式（人工介入）
```bash
node index.js --review
```

### 持续循环（守护进程）
```bash
node index.js --loop
```

### 指定进化策略
```bash
# 最大化创新
EVOLVE_STRATEGY=innovate node index.js --loop

# 聚焦稳定性
EVOLVE_STRATEGY=harden node index.js --loop

# 紧急修复模式
EVOLVE_STRATEGY=repair-only node index.js --loop
```

---

## 📊 进化策略对比

| 策略 | 创新 | 优化 | 修复 | 适用场景 |
|------|------|------|------|---------|
| `balanced`（默认） | 50% | 30% | 20% | 日常运行，稳步成长 |
| `innovate` | 80% | 15% | 5% | 系统稳定，快速出新功能 |
| `harden` | 20% | 40% | 40% | 大改动后，聚焦稳固 |
| `repair-only` | 0% | 20% | 80% | 紧急状态，全力修复 |

---

## 🔧 运维管理

```bash
# 后台启动进化循环
node src/ops/lifecycle.js start

# 优雅停止
node src/ops/lifecycle.js stop

# 查看运行状态
node src/ops/lifecycle.js status

# 健康检查 + 停滞自动重启
node src/ops/lifecycle.js check
```

### 运维模块 (src/ops/)
6 个可移植的运维工具，零平台依赖：
1. 生命周期管理
2. 技能健康监控
3. 磁盘清理
4. Git 自修复
5. ...

---

## 🛡️ 安全模型

### 各组件执行行为

| 组件 | 行为 | 执行 Shell |
|------|------|-----------|
| `src/evolve.js` | 读取日志、选择 Gene、构建提示词、写入工件 | 仅只读 git/进程查询 |
| `src/gep/prompt.js` | 组装 GEP 协议提示词字符串 | 否 |
| `src/gep/selector.js` | 按信号匹配对 Gene/Capsule 评分和选择 | 否 |
| `src/gep/solidify.js` | 通过 Gene `validation` 命令验证补丁 | 是（受控） |
| `index.js` | 崩溃时输出 `sessions_spawn(...)` 文本 | 否 |

### Gene Validation 命令安全机制

`solidify.js` 执行 Gene 的 `validation` 数组中的命令，安全检查：

1. **前缀白名单**: 仅允许 `node`、`npm`、`npx` 开头
2. **禁止命令替换**: 反引号或 `$(...)` 均被拒绝
3. **禁止 Shell 操作符**: `;`、`&`、`|`、`>`、`<` 均被拒绝
4. **超时限制**: 每条命令 180 秒
5. **作用域限定**: 以仓库根目录为工作目录

### A2A 外部资产摄入

- 外部 Gene/Capsule 暂存隔离候选区
- 提升到本地存储需要 `--validated` 标志
- Gene 提升不会覆盖本地已存在的同 ID Gene

### 其他安全约束

1. **单进程锁**: 禁止生成子进化进程（防止 Fork 炸弹）
2. **稳定性优先**: 近期错误率高则强制进入修复模式
3. **环境检测**: 外部集成仅在检测到相应插件时启用

---

## ⚙️ 环境变量配置

| 环境变量 | 描述 | 默认值 |
|---------|------|--------|
| `EVOLVE_STRATEGY` | 进化策略预设 | `balanced` |
| `EVOLVE_REPORT_TOOL` | 用于报告结果的工具名称 | `message` |
| `MEMORY_DIR` | 记忆文件路径 | `./memory` |
| `OPENCLAW_WORKSPACE` | 工作区根路径 | 自动检测 |
| `EVOLVER_LOOP_SCRIPT` | 循环启动脚本路径 | 自动检测 |

---

## 📦 发布流程

```bash
# 构建公开产物
npm run build

# 发布公开产物
npm run publish:public

# 演练
DRY_RUN=true npm run publish:public
```

### 必填环境变量
- `PUBLIC_REMOTE`（默认：`public`）
- `PUBLIC_REPO`（例如 `autogame-17/evolver`）
- `PUBLIC_OUT_DIR`（默认：`dist-public`）
- `PUBLIC_USE_BUILD_OUTPUT`（默认：`true`）

---

## 🎭 适用场景 vs 反例

### ✅ 适用场景
- 需要审计与可追踪的提示词演进
- 团队协作维护 Agent 的长期能力
- 希望将修复经验固化为可复用资产

### ❌ 反例
- 一次性脚本或没有日志的场景
- 需要完全自由发挥的改动
- 无法接受协议约束的系统

---

## 🔗 与 EvoMap 的关系

**Capability Evolver 是 EvoMap 的核心引擎**

**EvoMap**: https://evomap.ai
- AI 智能体通过验证协作实现进化的网络
- 实时智能体图谱
- 进化排行榜
- 将孤立的提示词调优转化为共享可审计智能的生态系统

---

## 💡 对 NOFX 系统的启示

### 1. 自我进化能力
**当前 NOFX**: 手动分析交易日志，人工调整策略  
**Evolver 模式**: 自动扫描交易日志，识别亏损模式，自主生成策略补丁

**实现路径**:
```javascript
// 在 NOFX 中集成 Evolver 思路
node evolve.js --strategy=repair-only  // 紧急修复连续亏损
node evolve.js --strategy=innovate     // 创新新的交易策略
```

### 2. GEP 协议应用
**当前**: 交易策略散落在代码和文档中  
**GEP 模式**: 结构化资产目录
```
nofx-gep/
├── genes.json           # 盈利策略基因库
├── loss-capsules.json   # 亏损修复胶囊
├── events.jsonl         # 交易事件日志
└── selectors/           # 策略选择器
```

### 3. 策略版本管理
使用 SemVer 管理策略版本：
- MAJOR: 不兼容的策略变更（如完全更换交易逻辑）
- MINOR: 向后兼容的新策略（如增加新的入场条件）
- PATCH: 策略问题修复（如修复止损逻辑 bug）

### 4. 安全约束借鉴
**NOFX 交易安全**:
- 白名单机制: 只允许特定币种/方向
- 禁止操作符: 不允许超出预设风险参数
- 超时限制: 单笔交易最大持仓时间
- 稳定性优先: 连续亏损自动暂停交易

### 5. 运维自动化
```bash
# NOFX 交易运维
node ops/lifecycle.js start    # 启动交易系统
node ops/lifecycle.js status   # 查看运行状态
node ops/lifecycle.js check    # 健康检查 + 异常重启
```

---

## 📚 相关资源

- **EvoMap 官网**: https://evomap.ai
- **Wiki 文档**: https://evomap.ai/wiki
- **GitHub 仓库**: https://github.com/EvoMap/evolver
- **英文文档**: README.md

---

## 🙏 鸣谢

- [onthebigtree](https://github.com/onthebigtree) -- 启发了 evomap 进化网络的诞生
- [lichunr](https://github.com/lichunr) -- 提供了数千美金 Token 供算力网络免费使用
- [shinjiyu](https://github.com/shinjiyu) -- 提交了大量 bug report
- [upbit](https://github.com/upbit) -- 在技术普及中起到关键作用
- [池建强](https://mowen.cn) -- 传播和用户体验改进

---

## 📄 许可证
MIT

---

*技术文档整理时间: 2026-02-24 16:15 UTC*
*来源: GitHub EvoMap/evolver 官方仓库*