#!/bin/bash
# NOFX-OpenClaw 极简安装/更新脚本 (v1.1.5)
# 仅负责代码同步与镜像启动，不触碰任何用户数据

set -e

REPO_DIR="/root/nofx"
BRANCH="openclaw"

echo "=========================================="
echo "🚀 正在更新/启动 NOFX-OpenClaw"
echo "=========================================="

# 1. 确保在正确的目录下
if [ ! -d "$REPO_DIR" ]; then
    echo "📂 正在执行首次安装..."
    git clone -b $BRANCH https://github.com/xbcvv/nofx-0210.git $REPO_DIR
    cd $REPO_DIR
else
    cd $REPO_DIR
    echo "🔄 正在同步代码..."
    git fetch origin $BRANCH
    git reset --hard origin/$BRANCH
fi

# 2. 仅在缺失时提醒，但不强制生成（由应用内核处理）
if [ ! -f ".env" ]; then
    echo "⚠️ 警告: 未发现 .env 文件。如果这是首次安装，请稍后手动配置或查看说明。"
fi

# 3. 拉取镜像并启动 (利用 Docker Volume 保证持久化)
echo "📡 正在同步镜像并重启容器..."
docker compose pull
docker compose up -d

echo "=========================================="
echo "✅ NOFX-OpenClaw 已就绪"
echo "📊 数据持久化于: $REPO_DIR/data"
echo "=========================================="
