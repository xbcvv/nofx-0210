# Linux 系统级 Cron 备用方案

## 概述

由于 OpenClaw Gateway 存在 Event Loop 阻塞问题，导致其内置的 Cron 任务不可靠，本方案配置 Linux 系统级 Cron 作为备用，确保 NOFX 交易监控任务能够稳定执行。

## 配置时间
- 创建时间: 2026-02-24
- 备用方案状态: ✅ 已启用

## 任务配置

### 任务1: 每小时检查
```
执行时间: 每小时整点 (0 * * * *)
脚本: /root/nofx-0210/scripts/cron/nofx_hourly_check.sh
功能:
  - 检查 Docker 容器状态
  - 获取最近订单和持仓
  - 统计盈亏数据
  - 异常时发送通知
日志: /root/nofx-0210/logs/cron/hourly_YYYYMMDD_HH.log
```

### 任务2: 每日分析
```
执行时间: 每天 00:00 (0 0 * * *)
脚本: /root/nofx-0210/scripts/cron/nofx_daily_analysis.sh
功能:
  - 汇总过去24小时交易数据
  - 分析各币种表现
  - 识别问题模式
  - 生成策略建议
日志: /root/nofx-0210/logs/cron/daily_YYYYMMDD.log
```

## 系统配置

### Crontab 配置
```bash
# 查看当前 crontab
crontab -l

# 编辑 crontab
crontab -e
```

当前配置:
```
# NOFX Cron 备用方案 - 每小时检查
0 * * * * /root/nofx-0210/scripts/cron/nofx_hourly_check.sh >> /root/nofx-0210/logs/cron/cron.log 2>&1

# NOFX Cron 备用方案 - 每日分析
0 0 * * * /root/nofx-0210/scripts/cron/nofx_daily_analysis.sh >> /root/nofx-0210/logs/cron/cron.log 2>&1
```

### 脚本位置
```
/root/nofx-0210/scripts/cron/
├── nofx_hourly_check.sh    # 每小时任务脚本
└── nofx_daily_analysis.sh  # 每日任务脚本
```

### 日志位置
```
/root/nofx-0210/logs/cron/
├── hourly_YYYYMMDD_HH.log  # 每小时日志
├── daily_YYYYMMDD.log      # 每日日志
└── cron.log                # Cron 执行日志
```

## 手动测试

### 测试每小时任务
```bash
/root/nofx-0210/scripts/cron/nofx_hourly_check.sh
```

### 测试每日任务
```bash
/root/nofx-0210/scripts/cron/nofx_daily_analysis.sh
```

## 监控和故障排除

### 检查 Cron 服务状态
```bash
# 检查 cron 服务
service cron status
# 或
systemctl status cron
```

### 查看日志
```bash
# 查看最近日志
tail -f /root/nofx-0210/logs/cron/cron.log

# 查看每小时日志
tail -f /root/nofx-0210/logs/cron/hourly_$(date +%Y%m%d_%H).log

# 查看每日日志
tail -f /root/nofx-0210/logs/cron/daily_$(date +%Y%m%d).log
```

### 检查 Cron 执行历史
```bash
# 查看系统日志中的 cron 记录
grep CRON /var/log/syslog | tail -20
```

## 与 OpenClaw Cron 的关系

| 特性 | OpenClaw Cron | Linux Cron (备用) |
|------|--------------|-------------------|
| 可靠性 | ⚠️ 不稳定 (Event Loop 阻塞) | ✅ 稳定 |
| 执行环境 | OpenClaw Gateway | Linux 系统 |
| 日志位置 | OpenClaw 会话 | /root/nofx-0210/logs/cron/ |
| 通知方式 | Telegram | 文件日志 + 可选 Telegram |
| 当前状态 | 作为备用 | 作为主要方案 |

## 双保险策略

1. **Linux Cron 作为主方案**: 确保任务稳定执行
2. **OpenClaw Cron 作为备用**: 在 Heartbeat 中继续检查，如发现问题可触发修复
3. **日志对比**: 定期对比两个系统的日志，确保一致性

## 维护操作

### 修改执行时间
```bash
crontab -e
# 修改时间后保存
```

### 禁用任务
```bash
# 注释掉不需要的任务
crontab -e
# 在行首添加 #
```

### 清理旧日志
脚本会自动清理:
- 每小时日志: 保留7天
- 每日日志: 保留30天

手动清理:
```bash
# 清理7天前的每小时日志
find /root/nofx-0210/logs/cron -name "hourly_*.log" -mtime +7 -delete

# 清理30天前的每日日志
find /root/nofx-0210/logs/cron -name "daily_*.log" -mtime +30 -delete
```

## 安全注意事项

1. **不修改运行中的 Docker**: 脚本只读取数据，不修改容器
2. **权限控制**: 脚本以 root 运行，确保 Docker 命令执行权限
3. **日志保护**: 日志文件包含敏感交易数据，注意权限设置

## Telegram 通知配置 (可选)

如需启用 Telegram 通知，设置环境变量:
```bash
export TELEGRAM_BOT_TOKEN="your_bot_token"
export TELEGRAM_CHAT_ID="your_chat_id"
```

添加到 ~/.bashrc 使其持久化:
```bash
echo 'export TELEGRAM_BOT_TOKEN="your_bot_token"' >> ~/.bashrc
echo 'export TELEGRAM_CHAT_ID="your_chat_id"' >> ~/.bashrc
```

## 参考

- OpenClaw Cron 问题: https://www.moltbook.com/posts/c408d2ee-96f1-4827-a5be-1bafe3b7d2c6
- Linux Cron 文档: man 5 crontab