#!/bin/bash

# 前端诊断和修复脚本

PROJECT_ROOT="/Users/niumc/Downloads/code/gostudy/memberhub/web/admin"
cd "$PROJECT_ROOT"

echo "========================================"
echo "  前端诊断和修复"
echo "========================================"
echo ""

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 1. 检查 node_modules
echo "1️⃣  检查依赖安装..."
if [ ! -d "node_modules" ]; then
    echo -e "${RED}❌ node_modules 目录不存在${NC}"
    echo -e "${YELLOW}正在安装依赖...${NC}"
    npm install
else
    echo -e "${GREEN}✅ 依赖已安装${NC}"
fi
echo ""

# 2. 检查关键文件
echo "2️⃣  检查关键文件..."
FILES=(
    "src/main.js"
    "src/App.vue"
    "src/router/index.js"
    "src/layouts/MainLayout.vue"
    "src/views/dashboard/Index.vue"
)

for file in "${FILES[@]}"; do
    if [ -f "$file" ]; then
        echo -e "${GREEN}✅ $file${NC}"
    else
        echo -e "${RED}❌ $file 不存在${NC}"
    fi
done
echo ""

# 3. 检查 vite.config.js
echo "3️⃣  检查 Vite 配置..."
if [ -f "vite.config.js" ]; then
    echo -e "${GREEN}✅ vite.config.js 存在${NC}"
else
    echo -e "${RED}❌ vite.config.js 不存在${NC}"
fi
echo ""

# 4. 清理缓存
echo "4️⃣  清理缓存..."
rm -rf node_modules/.vite
rm -rf dist
echo -e "${GREEN}✅ 缓存已清理${NC}"
echo ""

# 5. 检查端口占用
echo "5️⃣  检查端口 3000..."
PORT_IN_USE=$(lsof -ti:3000 || echo "")
if [ ! -z "$PORT_IN_USE" ]; then
    echo -e "${YELLOW}⚠️  端口 3000 已被占用 (PID: $PORT_IN_USE)${NC}"
    echo "是否要终止该进程？(y/n)"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        kill -9 $PORT_IN_USE
        echo -e "${GREEN}✅ 已终止进程${NC}"
    fi
else
    echo -e "${GREEN}✅ 端口 3000 可用${NC}"
fi
echo ""

# 6. 验证 package.json
echo "6️⃣  验证 package.json..."
if [ -f "package.json" ]; then
    # 检查关键依赖
    DEPS=("vue" "element-plus" "pinia" "vue-router" "axios" "vite")
    for dep in "${DEPS[@]}"; do
        if grep -q "\"$dep\"" package.json; then
            echo -e "${GREEN}  ✅ $dep${NC}"
        else
            echo -e "${RED}  ❌ $dep 未在 package.json 中${NC}"
        fi
    done
else
    echo -e "${RED}❌ package.json 不存在${NC}"
fi
echo ""

echo "========================================"
echo "  诊断完成"
echo "========================================"
echo ""
echo "📝 建议操作："
echo "1. 运行: npm install"
echo "2. 运行: npm run dev"
echo "3. 访问: http://localhost:3000"
echo ""
echo "如果问题仍然存在，请检查浏览器控制台的错误信息"
echo ""
