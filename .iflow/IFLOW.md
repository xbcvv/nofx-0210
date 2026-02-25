# NOFX Trading System Project

## 项目概述
NOFX 是一个开源 AI 交易操作系统，AI 驱动金融交易的基础设施层。

## 项目信息
- **GitHub**: https://github.com/xbcvv/nofx-0210
- **本地路径**: /root/nofx-0210
- **Fork来源**: NoFxAiOS/nofx (2026-02-10版本)
- **技术栈**: Go 1.21+, React 18+, TypeScript 5.0+

## 核心功能
- **多 AI 支持**: DeepSeek、通义千问、GPT、Claude、Gemini、Grok、Kimi
- **多交易所**: Binance、Bybit、OKX、Bitget、KuCoin、Gate、Hyperliquid、Aster DEX、Lighter
- **策略工作室**: 可视化策略构建器
- **AI 竞赛模式**: 多个 AI 交易员实时竞争
- **Web 配置**: 通过 Web 界面完成所有配置
- **实时仪表板**: 持仓、盈亏追踪、AI 决策日志

## 项目结构
```
nofx-0210/
├── trader/          # 交易核心模块
├── strategy/        # 策略管理
├── provider/        # 交易所接入
├── backtest/        # 回测系统
├── web/             # Web 前端
├── config/          # 配置文件
├── docker/          # Docker 部署
└── docs/            # 文档
```

## 关键配置
- **Docker Compose**: docker-compose.yml
- **环境变量**: .env.example
- **数据库**: SQLite (trader_orders, trader_positions)
- **安装脚本**: install.sh / install-stable.sh

## OpenClaw 集成
- **NOFX cron 任务**: 每小时检查交易状态
  - 任务ID: f4ba1f4f-b228-4aea-8541-8d34cb00f87a
  - 调度: 0 * * * * (Asia/Shanghai)
- **每日策略优化**: 每天 00:00 执行
  - 任务ID: ff11e012-ad0a-4ef5-afd0-3d61f20e40c6
  - 调度: 0 0 * * * (Asia/Shanghai)

## 交易代理技能
- **技能库**: /root/.openclaw/workspace/agents/nofx_trading_agent/skills/SKILL.md
- **学习来源**: Moltbook社区 (TradingLobster, madeleine-cupcake等)
- **核心技能**: Confluence Scoring, Dual-Source Verification, AI Discipline

## 常用命令
```bash
# 启动 NOFX
docker-compose up -d

# 检查容器
docker ps | grep nofx

# 查看日志
docker logs -f nofx-trading

# 数据库查询
docker exec nofx-trading sqlite3 /app/data/data.db "SELECT * FROM trader_orders ORDER BY created_at DESC LIMIT 10;"
```

## 官方链接
- **官网**: https://nofxai.com
- **数据站点**: https://nofxos.ai/dashboard
- **API 文档**: https://nofxos.ai/api-docs
- **社区**: https://t.me/nofx_dev_community