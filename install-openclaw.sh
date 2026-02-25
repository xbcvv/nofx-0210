#!/bin/bash
# NOFX-OpenClaw 一键安装/更新脚本 (修正版 v1.1.3)
# 修复：
# 1. 保护现有的 .env 文件不被覆盖
# 2. 修复 RSA 私钥格式问题 (直接写入 PEM 格式而非 Base64)

set -e

REPO_URL="https://github.com/xbcvv/nofx-0210.git"
REPO_DIR="/root/nofx"
BRANCH="openclaw"

echo "=========================================="
echo "🚀 正在启动 NOFX-OpenClaw (OpenClaw 版)"
echo "=========================================="

# 1. 检查并处理目录 (安全更新模式)
if [ ! -d "$REPO_DIR" ]; then
    echo "📂 正在克隆仓库到 $REPO_DIR..."
    git clone -b $BRANCH $REPO_URL $REPO_DIR
    cd $REPO_DIR
else
    echo "✅ 目录 $REPO_DIR 已存在，进入目录..."
    cd $REPO_DIR
    if [ -d ".git" ]; then
        echo "🔄 正在同步代码 (保留数据与环境配置)..."
        git fetch origin $BRANCH
        git reset --hard origin/$BRANCH
    else
        echo "⚠️ 警告: 目录存在但不是 Git 仓库。正在执行安全就地转换..."
        # 在临时目录克隆，只移动 Git 核心文件，确保不触碰现有数据
        git clone -b $BRANCH --depth 1 $REPO_URL /tmp/nofx_tmp
        cp -r /tmp/nofx_tmp/.git .
        git reset --hard origin/$BRANCH
        rm -rf /tmp/nofx_tmp
    fi
fi

# 2. 检查并生成环境变量 (.env) - 增加保护逻辑
if [ ! -f ".env" ]; then
    echo "📝 正在生成默认 .env 文件..."
    DATA_KEY=$(openssl rand -hex 32)
    # 修正：直接生成 PEM 格式的私钥并将其换行处理
    openssl genrsa 2048 > rsa_private.key
    RSA_KEY_CONTENT=$(cat rsa_private.key | sed ':a;N;$!ba;s/\n/\\n/g')
    rm rsa_private.key

    cat <<EOF > .env
TZ=Asia/Shanghai
NOFX_BACKEND_PORT=8080
NOFX_FRONTEND_PORT=3000
DATA_ENCRYPTION_KEY=$DATA_KEY
RSA_PRIVATE_KEY="$RSA_KEY_CONTENT"
TRANSPORT_ENCRYPTION=false
EOF
    echo "✅ 新 .env 文件已生成。"
else
    echo "✅ 发现现有的 .env 文件，已跳过生成，保留原有配置。"
fi

# 3. 拉取专属镜像
echo "📡 正在拉取 Docker Hub 镜像..."
docker compose pull

# 4. 重新启动服务
echo "🔄 正在重启服务..."
docker compose up -d

echo "=========================================="
echo "✅ NOFX-OpenClaw 已成功运行！"
echo "🌐 访问: http://你的服务器IP:3000"
echo "=========================================="
