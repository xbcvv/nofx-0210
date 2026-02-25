# NOFX Strategy Evolver Skill

## 概述

NOFX 交易策略智能进化与 EvoMap 全球基因融合系统。

通过每小时数据分析和每日深度优化，结合 EvoMap 全球策略基因库，实现交易策略的持续自我进化。

## 功能

- **每小时分析**: 交易数据收集、决策质量评估、异常检测
- **每日优化**: 策略表现回顾、全球基因检索、策略融合进化
- **EvoMap 集成**: 基因胶囊发布、全球策略共享、声誉经济参与
- **自动报告**: 生成 Markdown 分析报告，推送到 docs/openclaw/

## 触发方式

### 1. 每小时分析（自动）
```bash
# Crontab: 05 * * * *
cd ~/nofx-0210 && python3 skills/nofx_strategy_evolver/scripts/hourly_analyze.py
```

### 2. 每日进化（自动）
```bash
# Crontab: 10 0 * * *
cd ~/nofx-0210 && python3 skills/nofx_strategy_evolver/scripts/daily_evolve.py
```

### 3. 手动触发
```bash
# 每小时分析
python3 skills/nofx_strategy_evolver/scripts/hourly_analyze.py

# 每日进化
python3 skills/nofx_strategy_evolver/scripts/daily_evolve.py

# 同步 EvoMap
python3 skills/nofx_strategy_evolver/scripts/evomap_sync.py --publish --search
```

## 配置

配置文件: `skills/nofx_strategy_evolver/config/evolver.yaml`

```yaml
evolution:
  # 本地配置
  data_dir: "data"
  output_dir: "docs/openclaw"
  strategy_dir: "docs/openclaw/strategies"
  
  # EvoMap 配置
  evomap:
    enabled: true
    api_base: "https://evomap.ai/api/v1"
    agent_name: "xh_lx024"
    auto_publish: true
    min_gene_rating: 4.0
    
  # 分析参数
  analysis:
    hourly_interval: 5  # 每小时第5分钟执行
    daily_time: "00:10" # 每天00:10执行
    min_confidence_threshold: 76
    max_drawdown_alert: 10
    consecutive_losses_alert: 3
    
  # 基因融合参数
  fusion:
    auto_apply_low_risk: true
    manual_review_high_risk: true
    ab_test_threshold: 0.7
    max_daily_fusions: 3
```

## 依赖

- Python 3.9+
- sqlite3
- requests
- pandas (可选，用于高级分析)

```bash
pip install -r skills/nofx_strategy_evolver/requirements.txt
```

## 数据结构

### 输入
- `data/raw/decisions_YYYYMMDD_HH.json` - AI原始决策
- `data/data.db` - SQLite交易数据库
- `data/nofx_*.log` - 系统日志

### 输出
- `docs/openclaw/hourly_analysis_YYYYMMDD_HH.md` - 每小时报告
- `docs/openclaw/daily_analysis_YYYYMMDD.md` - 每日深度分析
- `docs/openclaw/strategy_optimization_YYYYMMDD.md` - 优化建议
- `docs/openclaw/strategies/optimized/*.json` - 优化后策略

## 安全

- 所有策略修改需人工确认后才应用
- 高风险变更自动标记为待审核
- 完整备份机制，支持一键回滚
- API密钥存储在环境变量，不提交到仓库

## 作者

小七 (xh_lx024)

## 版本

v1.0.0 - 2026-02-25

## 许可证

MIT