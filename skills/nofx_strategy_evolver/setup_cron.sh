#!/bin/bash
# NOFX 策略进化定时任务配置

# 每小时分析 - 05分执行
echo "配置每小时分析任务..."
(crontab -l 2>/dev/null; echo "5 * * * * cd /root/nofx-0210 && python3 skills/nofx_strategy_evolver/scripts/hourly_analyze.py >> logs/evolver_hourly.log 2>&1") | crontab -

# 每日进化 - 00:10执行
echo "配置每日进化任务..."
(crontab -l 2>/dev/null; echo "10 0 * * * cd /root/nofx-0210 && python3 skills/nofx_strategy_evolver/scripts/daily_evolve.py >> logs/evolver_daily.log 2>&1") | crontab -

echo "✅ 定时任务配置完成！"
echo ""
echo "查看当前crontab:"
crontab -l | grep evolver