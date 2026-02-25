#!/bin/bash
# NOFX-OpenClaw 一键安装/更新脚本 (增强版 v1.1.2)
# 修复路径问题：统一使用 /root/nofx
# 修复变量问题：自动生成必要的加密 KEY

set -e

REPO_URL="https://github.com/xbcvv/nofx-0210.git"
REPO_DIR="/root/nofx"
BRANCH="openclaw"

echo "=========================================="
echo "🚀 正在启动 NOFX-OpenClaw (OpenClaw 版)"
echo "=========================================="

# 1. 检查并处理目录 (强制使用 /root/nofx)
if [ ! -d "$REPO_DIR" ]; then
    echo "📂 正在克隆仓库到 $REPO_DIR..."
    git clone -b $BRANCH $REPO_URL $REPO_DIR
    cd $REPO_DIR
else
    echo "✅ 目录 $REPO_DIR 已存在，进入目录..."
    cd $REPO_DIR
    if [ ! -d ".git" ]; then
        echo "⚠️ 警告: 目录存在但不是 Git 仓库。正在备份并重新克隆..."
        mv $REPO_DIR ${REPO_DIR}_backup_$(date +%s)
        git clone -b $BRANCH $REPO_URL $REPO_DIR
        cd $REPO_DIR
    else
        echo "🔄 正在同步代码..."
        git fetch origin $BRANCH
        git reset --hard origin/$BRANCH
    fi
fi

# 2. 检查并生成环境变量 (.env)
if [ ! -f ".env" ]; then
    echo "📝 正在生成默认 .env 文件..."
    # 生成随机的加密 KEY
    DATA_KEY=$(openssl rand -hex 32)
    # 生成 RSA 私钥 (如果系统需要)
    openssl genrsa -out rsa_private.key 2048
    RSA_KEY=$(cat rsa_private.key | base64 -w 0)
    rm rsa_private.key

    cat <<EOF > .env
TZ=Asia/Shanghai
NOFX_BACKEND_PORT=8080
NOFX_FRONTEND_PORT=3000
DATA_ENCRYPTION_KEY=$DATA_KEY
RSA_PRIVATE_KEY=$RSA_KEY
TRANSPORT_ENCRYPTION=false
EOF
    echo "✅ .env 文件已生成，包含自动生成的安全密钥。"
fi

# 3. 拉取专属镜像
echo "📡 正在拉取 Docker Hub 镜像..."
docker compose pull

# 4. 重新启动服务
echo "🔄 正在重启服务..."
docker compose up -d

echo "=========================================="
echo "✅ NOFX-OpenClaw 已成功运行！"
echo "📊 路径: $REPO_DIR"
echo "🌐 访问: http://你的服务器IP:3000"
echo "=========================================="
