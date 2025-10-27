#!/bin/bash

# TinyGo 启动脚本
# 此脚本帮助用户安全地启动 TinyGo 短链接服务

set -e

echo "🔗 TinyGo 短链接服务启动脚本"
echo "================================"

# 检查是否设置了认证环境变量
if [ -z "$TINYGO_AUTH_USERNAME" ] || [ -z "$TINYGO_AUTH_PASSWORD" ]; then
    echo "❌ 错误：未设置认证凭据！"
    echo ""
    echo "请设置以下环境变量："
    echo "  export TINYGO_AUTH_USERNAME=\"your_username\""
    echo "  export TINYGO_AUTH_PASSWORD=\"your_password\""
    echo ""
    echo "或者复制并编辑环境变量示例文件："
    echo "  cp env.example .env"
    echo "  # 编辑 .env 文件设置你的凭据"
    echo "  source .env"
    echo ""
    exit 1
fi

echo "✅ 认证凭据已设置"
echo "📦 编译项目..."

# 编译项目
go build -o bin/urlshort ./cmd/server

if [ $? -ne 0 ]; then
    echo "❌ 编译失败！"
    exit 1
fi

echo "✅ 编译成功"
echo "🚀 启动服务..."

# 启动服务
./bin/urlshort
