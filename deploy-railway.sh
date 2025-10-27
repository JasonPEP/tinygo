#!/bin/bash

# Railway 部署脚本
# 此脚本帮助用户快速部署到 Railway

set -e

echo "🚀 TinyGo Railway 部署脚本"
echo "=========================="

# 检查 Railway CLI 是否安装
if ! command -v railway &> /dev/null; then
    echo "❌ Railway CLI 未安装！"
    echo ""
    echo "请先安装 Railway CLI："
    echo "  npm install -g @railway/cli"
    echo "  或者"
    echo "  brew install railway"
    echo ""
    exit 1
fi

echo "✅ Railway CLI 已安装"

# 检查是否已登录
if ! railway whoami &> /dev/null; then
    echo "🔐 请先登录 Railway："
    railway login
fi

echo "✅ 已登录 Railway"

# 检查是否已初始化项目
if [ ! -f ".railway/project.json" ]; then
    echo "📦 初始化 Railway 项目..."
    railway init
fi

echo "✅ Railway 项目已初始化"

# 检查环境变量
echo ""
echo "🔧 环境变量检查："
echo "请确保在 Railway Dashboard 中设置了以下环境变量："
echo ""
echo "必需的环境变量："
echo "  TINYGO_AUTH_USERNAME=your_username"
echo "  TINYGO_AUTH_PASSWORD=your_password"
echo "  TINYGO_BASE_URL=https://your-app-name.railway.app"
echo ""
echo "可选的环境变量："
echo "  TINYGO_ADDR=:8080"
echo "  TINYGO_DATABASE_DRIVER=sqlite"
echo "  TINYGO_DATABASE_DSN=data/tinygo.db"
echo "  TINYGO_LOG_LEVEL=info"
echo "  TINYGO_LOG_FORMAT=json"
echo ""

read -p "是否已设置环境变量？(y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "请先在 Railway Dashboard 中设置环境变量，然后重新运行此脚本。"
    exit 1
fi

echo "🚀 开始部署..."

# 部署到 Railway
railway up

echo ""
echo "✅ 部署完成！"
echo ""
echo "📊 查看日志："
echo "  railway logs"
echo ""
echo "🌐 查看应用："
echo "  railway open"
echo ""
echo "📈 监控应用："
echo "  访问 Railway Dashboard 查看应用状态和指标"
