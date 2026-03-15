#!/bin/bash

# 会员管理系统前端快速启动脚本

set -e

echo "================================"
echo "会员管理系统 - 管理后台前端"
echo "================================"
echo ""

# 检查 Node.js
if ! command -v node &> /dev/null; then
    echo "❌ 错误: 未检测到 Node.js，请先安装 Node.js (推荐 v18+)"
    echo "下载地址: https://nodejs.org/"
    exit 1
fi

NODE_VERSION=$(node -v)
echo "✅ Node.js 版本: $NODE_VERSION"

# 检查 npm
if ! command -v npm &> /dev/null; then
    echo "❌ 错误: 未检测到 npm"
    exit 1
fi

NPM_VERSION=$(npm -v)
echo "✅ npm 版本: $NPM_VERSION"
echo ""

# 进入前端目录
cd "$(dirname "$0")"

# 检查是否已安装依赖
if [ ! -d "node_modules" ]; then
    echo "📦 正在安装依赖..."
    npm install
    echo ""
else
    echo "✅ 依赖已安装"
    echo ""
fi

# 显示启动信息
echo "================================"
echo "🚀 启动开发服务器"
echo "================================"
echo ""
echo "访问地址: http://localhost:3000"
echo "默认账号: admin"
echo "默认密码: admin123"
echo ""
echo "提示: 请确保后端 API 服务已在 http://localhost:8080 启动"
echo ""
echo "按 Ctrl+C 停止服务器"
echo ""

# 启动开发服务器
npm run dev
