#!/bin/bash
# NOFX-OpenClaw 一键安装/更新脚本 (增强版 - 支持全新安装)
# 版本: 1.1.1
# 镜像: qdw1010/nofx-backend:openclaw, qdw1010/nofx-frontend:openclaw

set -e

REPO_URL="https://github.com/xbcvv/nofx-0210.git"
REPO_DIR="nofx-0210"
BRANCH="openclaw"

echo "=========================================="
echo "🚀 正在启动 NOFX-OpenClaw (OpenClaw 版)"
echo "=========================================="

# 1. 检查并处理目录
if [ ! -d "$REPO_DIR" ]; then
    echo "📂 正在克隆仓库..."
    git clone -b $BRANCH $REPO_URL $REPO_DIR
    cd $REPO_DIR
else
    echo "✅ 目录已存在，进入目录..."
    cd $REPO_DIR
    if [ ! -d ".git" ]; then
        echo "⚠️ 警告: 目录存在但不是 Git 仓库。正在清理并重新克隆..."
        cd ..
        rm -rf $REPO_DIR
        git clone -b $BRANCH $REPO_URL $REPO_DIR
        cd $REPO_DIR
    else
        echo "🔄 正在同步分支..."
        git checkout $BRANCH
        git pull origin $BRANCH || echo "⚠️ 远程分支可能尚未推送，跳过 pull"
    fi
fi

# 2. 拉取专属镜像
echo "📡 正在拉取 qdw1010/nofx-*:openclaw 镜像..."
docker compose pull

# 3. 重新启动服务
echo "🔄 正在重启服务..."
docker compose up -d

echo "=========================================="
echo "✅ NOFX-OpenClaw 已成功运行！"
echo "📊 请通过浏览器访问前端面板。"
echo "=========================================="
