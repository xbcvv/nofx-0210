#!/bin/bash
# NOFX-OpenClaw 一键安装/更新脚本
# 版本: 1.0.0
# 镜像: qdw1010/nofx-*:openclaw

set -e

echo "=========================================="
echo "🚀 正在启动 NOFX-OpenClaw (OpenClaw 版)"
echo "=========================================="

# 1. 强制切换到 openclaw 分支
git checkout openclaw

# 2. 拉取最新代码
git pull origin openclaw || echo "⚠️ 远程分支可能尚未推送，跳过 pull"

# 3. 拉取专属镜像
echo "📡 正在拉取 qdw1010/nofx-*:openclaw 镜像..."
docker-compose pull

# 4. 重新启动容器
echo "🔄 正在重启服务..."
docker-compose up -d

echo "=========================================="
echo "✅ NOFX-OpenClaw 已成功运行！"
echo "📊 请通过浏览器访问前端面板。"
echo "=========================================="
