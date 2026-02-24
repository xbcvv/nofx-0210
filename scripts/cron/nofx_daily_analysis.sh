#!/bin/bash
# NOFX 每日策略优化脚本 - Linux 系统级 Cron 备用方案
# 执行时间: 每天 00:00
# 创建时间: 2026-02-24

# 配置
LOG_DIR="/root/nofx-0210/logs/cron"
DATE=$(date +%Y%m%d)
LOG_FILE="${LOG_DIR}/daily_${DATE}.log"
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

# 记录开始
echo "[${TIMESTAMP}] ====== NOFX 每日策略优化分析开始 ======" >> "${LOG_FILE}"

# 1. 汇总过去24小时的交易
echo "[${TIMESTAMP}] 1. 汇总过去24小时交易数据..." >> "${LOG_FILE}"

DAILY_STATS=$(docker exec nofx-trading sqlite3 /app/data/data.db "
SELECT 
  COUNT(*) as total_trades,
  SUM(CASE WHEN realized_pnl > 0 THEN 1 ELSE 0 END) as win_count,
  SUM(CASE WHEN realized_pnl < 0 THEN 1 ELSE 0 END) as loss_count,
  ROUND(SUM(realized_pnl), 4) as total_pnl,
  ROUND(AVG(realized_pnl), 4) as avg_pnl,
  ROUND(MAX(realized_pnl), 4) as max_win,
  ROUND(MIN(realized_pnl), 4) as max_loss,
  ROUND(SUM(fee), 4) as total_fee
FROM trader_positions 
WHERE status = 'CLOSED' 
  AND created_at >= strftime('%s', 'now', '-24 hours') * 1000;" 2>/dev/null)

echo "[${TIMESTAMP}] 24小时统计:" >> "${LOG_FILE}"
echo "${DAILY_STATS}" | while IFS='|' read -r total wins losses pnl avg max_win max_loss fee; do
    echo "  总交易: ${total}" >> "${LOG_FILE}"
    echo "  盈利: ${wins}" >> "${LOG_FILE}"
    echo "  亏损: ${losses}" >> "${LOG_FILE}"
    if [ "$total" -gt 0 ] 2>/dev/null; then
        WIN_RATE=$(echo "scale=2; $wins * 100 / $total" | bc 2>/dev/null || echo "N/A")
        echo "  胜率: ${WIN_RATE}%" >> "${LOG_FILE}"
    fi
    echo "  总盈亏: ${pnl} USDT" >> "${LOG_FILE}"
    echo "  平均盈亏: ${avg} USDT" >> "${LOG_FILE}"
    echo "  最大盈利: ${max_win} USDT" >> "${LOG_FILE}"
    echo "  最大亏损: ${max_loss} USDT" >> "${LOG_FILE}"
    echo "  手续费: ${fee} USDT" >> "${LOG_FILE}"
done

# 2. 分析各币种表现
echo "[${TIMESTAMP}] 2. 各币种表现分析..." >> "${LOG_FILE}"
docker exec nofx-trading sqlite3 /app/data/data.db "
SELECT symbol, side,
  COUNT(*) as trades,
  SUM(CASE WHEN realized_pnl > 0 THEN 1 ELSE 0 END) as wins,
  ROUND(SUM(realized_pnl), 4) as total_pnl,
  ROUND(AVG(realized_pnl), 4) as avg_pnl
FROM trader_positions
WHERE status = 'CLOSED' 
  AND created_at >= strftime('%s', 'now', '-24 hours') * 1000
GROUP BY symbol, side
ORDER BY total_pnl DESC;" 2>/dev/null >> "${LOG_FILE}"

# 3. 识别问题模式
echo "[${TIMESTAMP}] 3. 问题模式识别..." >> "${LOG_FILE}"

# 检查连续亏损币种
echo "[${TIMESTAMP}] 连续亏损币种 (2笔以上):" >> "${LOG_FILE}"
docker exec nofx-trading sqlite3 /app/data/data.db "
SELECT symbol, side,
  COUNT(*) as trades,
  SUM(CASE WHEN realized_pnl > 0 THEN 1 ELSE 0 END) as wins,
  SUM(CASE WHEN realized_pnl < 0 THEN 1 ELSE 0 END) as losses
FROM trader_positions
WHERE status = 'CLOSED' 
  AND created_at >= strftime('%s', 'now', '-24 hours') * 1000
GROUP BY symbol, side
HAVING losses >= 2
ORDER BY losses DESC;" 2>/dev/null >> "${LOG_FILE}"

# 4. 生成策略建议
echo "[${TIMESTAMP}] 4. 策略优化建议..." >> "${LOG_FILE}"

# 根据数据分析生成建议
WIN_COUNT=$(echo "$DAILY_STATS" | cut -d'|' -f2)
LOSS_COUNT=$(echo "$DAILY_STATS" | cut -d'|' -f3)
TOTAL_PNL=$(echo "$DAILY_STATS" | cut -d'|' -f4)

if [ -n "$WIN_COUNT" ] && [ -n "$LOSS_COUNT" ]; then
    if [ "$LOSS_COUNT" -gt "$WIN_COUNT" ]; then
        echo "[${TIMESTAMP}] ⚠️ 建议: 胜率偏低 (${WIN_COUNT}/${LOSS_COUNT})，建议:" >> "${LOG_FILE}"
        echo "  - 收紧入场条件，提高 Confluence 分数要求" >> "${LOG_FILE}"
        echo "  - 检查趋势判断逻辑" >> "${LOG_FILE}"
        echo "  - 考虑暂停表现差的币种" >> "${LOG_FILE}"
    fi
    
    if (( $(echo "$TOTAL_PNL < 0" | bc -l 2>/dev/null || echo 0) )); then
        echo "[${TIMESTAMP}] 🚨 警告: 今日亏损，建议回顾策略!" >> "${LOG_FILE}"
    else
        echo "[${TIMESTAMP}] ✅ 今日盈利，策略表现良好" >> "${LOG_FILE}"
    fi
fi

# 5. 记录到 OpenClaw 工作区
OPENCLAW_LOG="/root/.openclaw/workspace/memory/nofx_daily_summary_${DATE}.md"
echo "# NOFX 每日交易汇总 - ${TIMESTAMP}" > "${OPENCLAW_LOG}"
echo "" >> "${OPENCLAW_LOG}"
echo "## 24小时统计" >> "${OPENCLAW_LOG}"
echo "${DAILY_STATS}" >> "${OPENCLAW_LOG}"
echo "" >> "${OPENCLAW_LOG}"
echo "## 分析" >> "${OPENCLAW_LOG}"
echo "- 生成时间: ${TIMESTAMP}" >> "${OPENCLAW_LOG}"
echo "- 来源: Linux Cron 备用方案" >> "${OPENCLAW_LOG}"

# 记录结束
echo "[${TIMESTAMP}] ====== NOFX 每日策略优化分析完成 ======" >> "${LOG_FILE}"
echo "" >> "${LOG_FILE}"

# 发送 Telegram 通知 (如配置了)
if [ -n "$TELEGRAM_BOT_TOKEN" ] && [ -n "$TELEGRAM_CHAT_ID" ]; then
    SUMMARY="📊 NOFX 每日报告 (${DATE}):\n总盈亏: ${TOTAL_PNL} USDT\n胜: ${WIN_COUNT} | 负: ${LOSS_COUNT}"
    curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
        -d "chat_id=${TELEGRAM_CHAT_ID}" \
        -d "text=${SUMMARY}" 2>/dev/null
fi

# 保留最近30天的日志
find "${LOG_DIR}" -name "daily_*.log" -mtime +30 -delete 2>/dev/null

exit 0