# Evolver 定时进化配置报告

## 配置时间
2026-02-24 16:23 UTC

---

## 📋 配置概览

**部署路径**: `/root/evolver-deploy/`  
**运行模式**: 方案2（本地为主 + 网络辅助）  
**状态**: ✅ 配置完成，运行正常

---

## ⏰ 定时任务配置

### 1. 常规进化任务（每4小时）
```cron
0 */4 * * * /root/evolver-deploy/scripts/nofx_evolve.sh >> /root/evolver-deploy/logs/cron/crontab.log 2>&1
```

**功能**:
- 扫描 OpenClaw 会话日志
- 分析 NOFX 相关错误和模式
- 运行优化策略（balanced）
- 生成修复补丁
- 清理旧日志（保留30天）

**执行时间**:
- 00:00, 04:00, 08:00, 12:00, 16:00, 20:00

### 2. 深度创新任务（每天凌晨2点）
```cron
0 2 * * * cd /root/evolver-deploy && EVOLVE_STRATEGY=innovate node index.js run --intent=innovate >> /root/evolver-deploy/logs/cron/innovate_$(date +\%Y\%m\%d).log 2>&1
```

**功能**:
- 全面分析24小时日志
- 创新策略（innovate）
- 生成新功能和优化建议
- 探索新的交易模式

---

## 🛠️ 管理脚本

### evolver.sh - 生命周期管理
**路径**: `/root/evolver-deploy/evolver.sh`

**用法**:
```bash
# 启动持续进化（后台守护进程）
./evolver.sh start

# 停止进化
./evolver.sh stop

# 查看状态
./evolver.sh status

# 重启服务
./evolver.sh restart

# 查看实时日志
./evolver.sh logs
```

### nofx_evolve.sh - NOFX 专用进化
**路径**: `/root/evolver-deploy/scripts/nofx_evolve.sh`

**功能**: 专为 NOFX 交易系统定制的进化脚本
- 分析交易相关日志
- 检测错误模式
- 生成优化补丁

---

## 📁 目录结构

```
/root/evolver-deploy/
├── index.js                    # 主入口
├── evolver.sh                  # 生命周期管理脚本 ✅
├── .env                        # 环境变量配置
├── memory/                     # 记忆文件目录
│   └── (会话日志)
├── assets/gep/                 # GEP 资产目录
│   ├── genes.json             # 基因库
│   ├── capsules.json          # 胶囊库
│   └── events.jsonl           # 事件日志
├── logs/                       # 日志目录
│   ├── cron/                  # 定时任务日志
│   │   ├── evolve_YYYYMMDD_HH.log
│   │   ├── innovate_YYYYMMDD.log
│   │   └── crontab.log
│   └── evolver.log            # 守护进程日志
├── scripts/                    # 脚本目录
│   └── nofx_evolve.sh         # NOFX 专用进化脚本 ✅
└── src/                        # 源码目录
    ├── evolve.js
    ├── gep/
    └── ops/
```

---

## ⚙️ 环境变量配置

**文件**: `/root/evolver-deploy/.env`

```bash
EVOLVE_STRATEGY=balanced           # 默认策略
EVOLVE_REPORT_TOOL=message         # 报告工具
MEMORY_DIR=/root/evolver-deploy/memory
OPENCLAW_WORKSPACE=/root/.openclaw/workspace
EVOMAP_API_KEY=moltbook_sk_...     # Moltbook API Key
EVOMAP_AGENT_NAME=xh_lx024         # 代理名称
EVOLVER_LOOP_SCRIPT=/root/evolver-deploy/index.js
```

---

## 📊 进化策略

### 常规进化（每4小时）
- **策略**: `balanced`（平衡）
- **意图**: `optimize`（优化）
- **比例**: 创新50% + 优化30% + 修复20%
- **目标**: 稳步改进，修复问题

### 深度进化（每天凌晨2点）
- **策略**: `innovate`（创新）
- **意图**: `innovate`（创新）
- **比例**: 创新80% + 优化15% + 修复5%
- **目标**: 探索新功能，突破现有模式

---

## 🔍 监控与日志

### 查看定时任务日志
```bash
# 查看最新进化日志
tail -f /root/evolver-deploy/logs/cron/evolve_$(date +%Y%m%d_%H).log

# 查看最新创新日志
tail -f /root/evolver-deploy/logs/cron/innovate_$(date +%Y%m%d).log

# 查看 crontab 执行日志
tail -f /root/evolver-deploy/logs/cron/crontab.log
```

### 查看守护进程日志
```bash
./evolver.sh logs
```

### 检查运行状态
```bash
./evolver.sh status
```

---

## 🧬 预期进化方向

基于 OpenClaw 日志分析，Evolver 可能会优化：

1. **网络连接稳定性**
   - Telegram 连接问题
   - 飞书连接问题
   - API 超时处理

2. **任务执行效率**
   - 工具调用优化
   - 会话管理改进
   - 错误处理增强

3. **NOFX 交易相关**
   - 策略执行逻辑
   - 风险控制机制
   - 数据分析流程

4. **记忆管理**
   - MEMORY.md 更新机制
   - 会话日志压缩
   - 关键信息提取

---

## 🚀 手动触发进化

### 立即运行一次进化
```bash
cd /root/evolver-deploy
node index.js
```

### 以特定策略运行
```bash
# 紧急修复
cd /root/evolver-deploy
EVOLVE_STRATEGY=repair-only node index.js

# 激进创新
cd /root/evolver-deploy
EVOLVE_STRATEGY=innovate node index.js
```

### 启动持续守护进程
```bash
cd /root/evolver-deploy
./evolver.sh start
```

---

## 🌐 EvoMap 网络连接

**当前配置**: 已启用网络辅助
- ✅ API Key 已配置
- ✅ 代理名称: xh_lx024
- ✅ 可下载全球基因胶囊
- ✅ 可分享本地进化成果

**网络功能**:
- 继承其他 Agent 的成功经验
- 分享 NOFX 交易策略
- 获取声誉值和 Credit

---

## 📝 维护命令

### 清理旧日志
```bash
# 手动清理30天前的日志
find /root/evolver-deploy/logs/cron -name "*.log" -mtime +30 -delete
```

### 备份 GEP 资产
```bash
# 备份基因库和胶囊库
cp -r /root/evolver-deploy/assets/gep ~/backup/gep_$(date +%Y%m%d)
```

### 更新 Evolver
```bash
cd /root/evolver-deploy
git pull origin main
npm install
```

---

## ⚠️ 注意事项

1. **首次运行**: 可能需要几分钟完成初始化
2. **日志空间**: 定期检查日志大小，避免磁盘占满
3. **安全边界**: Evolver 不会修改核心引擎代码（受 GEP 协议保护）
4. **人工审查**: 重要补丁建议人工审查后再应用

---

## 📚 相关文档

- 技术实现: `~/nofx-0210/docs/openclaw/evolver_technical_details_20260224_1615.md`
- 概念介绍: `~/nofx-0210/docs/openclaw/evolver_evomap_research_20260224_1610.md`
- 学习报告: `~/nofx-0210/docs/openclaw/moltbook_learning_report_20260224_1600.md`

---

*配置完成时间: 2026-02-24 16:23 UTC*
*配置者: 小七 (xh_lx024)*