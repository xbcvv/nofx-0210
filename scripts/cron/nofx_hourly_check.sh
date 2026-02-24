#!/bin/bash
# NOFX 每小时检查脚本 - Linux 系统级 Cron 备用方案
# 执行时间: 每小时整点
# 创建时间: 2026-02-24

# 配置
LOG_DIR="/root/nofx-0210/logs/cron"
DATE_HOUR=$(date +%Y%m%d_%H)
LOG_FILE="${LOG_DIR}/hourly_${DATE_HOUR}.log"
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

# 记录开始
echo "[${TIMESTAMP}] ====== NOFX 每小时检查开始 ======" >> "${LOG_FILE}"

# 1. 检查 Docker 容器状态
echo "[${TIMESTAMP}] 1. 检查 Docker 容器状态..." >> "${LOG_FILE}"
CONTAINER_STATUS=$(docker ps --filter "name=nofx" --format "{{.Status}}" 2>/dev/null)
if [ -z "$CONTAINER_STATUS" ]; then
    echo "[${TIMESTAMP}] ⚠️ 警告: nofx 容器未运行!" >> "${LOG_FILE}"
    echo "[${TIMESTAMP}] 尝试启动容器..." >> "${LOG_FILE}"
    cd /root/nofx-0210 && docker-compose up -d 2>&1 >> "${LOG_FILE}"
else
    echo "[${TIMESTAMP}] ✅ 容器状态: ${CONTAINER_STATUS}" >> "${LOG_FILE}"
fi

# 2. 获取最近订单
echo "[${TIMESTAMP}] 2. 获取最近订单..." >> "${LOG_FILE}"
docker exec nofx-trading sqlite3 /app/data/data.db "
SELECT id, symbol, side, type, quantity, status, avg_fill_price, 
       datetime(created_at/1000, 'unixepoch') as created_time 
FROM trader_orders 
WHERE created_at >= strftime('%s', 'now', '-1 hour') * 1000 
ORDER BY created_at DESC LIMIT 10;" 2>/dev/null >> "${LOG_FILE}"

# 3. 获取最近平仓持仓
echo "[${TIMESTAMP}] 3. 获取最近平仓持仓..." >> "${LOG_FILE}"
docker exec nofx-trading sqlite3 /app/data/data.db "
SELECT id, symbol, entry_price, quantity, realized_pnl, status, 
       datetime(created_at/1000, 'unixepoch') as created_time
FROM trader_positions 
WHERE created_at >= strftime('%s', 'now', '-1 hour') * 1000 
  AND status = 'CLOSED' 
ORDER BY created_at DESC LIMIT 10;" 2>/dev/null >> "${LOG_FILE}"

# 4. 统计盈亏
echo "[${TIMESTAMP}] 4. 盈亏统计..." >> "${LOG_FILE}"
PROFIT_STATS=$(docker exec nofx-trading sqlite3 /app/data/data.db "
SELECT 
  COUNT(*) as total,
  SUM(CASE WHEN realized_pnl > 0 THEN 1 ELSE 0 END) as wins,
  SUM(CASE WHEN realized_pnl < 0 THEN 1 ELSE 0 END) as losses,
  ROUND(SUM(realized_pnl), 4) as total_pnl
FROM trader_positions 
WHERE status = 'CLOSED' 
  AND created_at >= strftime('%s', 'now', '-1 hour') * 1000;" 2>/dev/null)

if [ -n "$PROFIT_STATS" ]; then
    echo "[${TIMESTAMP}] 本小时统计: ${PROFIT_STATS}" >> "${LOG_FILE}"
    
    # 检查连续亏损
    LOSSES=$(echo "$PROFIT_STATS" | cut -d'|' -f3)
    if [ "$LOSSES" -ge 3 ] 2>/dev/null; then
        echo "[${TIMESTAMP}] 🚨 警告: 连续亏损 ${LOSSES} 笔!" >> "${LOG_FILE}"
        # 发送 Telegram 通知 (如配置了)
        if [ -n "$TELEGRAM_BOT_TOKEN" ] && [ -n "$TELEGRAM_CHAT_ID" ]; then
            curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
                -d "chat_id=${TELEGRAM_CHAT_ID}" \
                -d "text=⚠️ NOFX 系统警告: 过去1小时亏损 ${LOSSES} 笔交易，请检查策略!" 2>/dev/null
        fi
    fi
    
    # 检查大幅亏损
    TOTAL_PNL=$(echo "$PROFIT_STATS" | cut -d'|' -f4)
    if (( $(echo "$TOTAL_PNL < -100" | bc -l 2>/dev/null || echo 0) )); then
        echo "[${TIMESTAMP}] 🚨 警告: 本小时大幅亏损 ${TOTAL_PNL}!" >> "${LOG_FILE}"
        if [ -n "$TELEGRAM_BOT_TOKEN" ] && [ -n "$TELEGRAM_CHAT_ID" ]; then
            curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
                -d "chat_id=${TELEGRAM_CHAT_ID}" \
                -d "text=🚨 NOFX 系统警告: 过去1小时大幅亏损 ${TOTAL_PNL} USDT!" 2>/dev/null
        fi
    fi
fi

# 记录结束
echo "[${TIMESTAMP}] ====== NOFX 每小时检查完成 ======" >> "${LOG_FILE}"
echo "" >> "${LOG_FILE}"

# 保留最近7天的日志
find "${LOG_DIR}" -name "hourly_*.log" -mtime +7 -delete 2>/dev/null

exit 0