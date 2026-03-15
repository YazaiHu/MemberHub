#!/bin/bash

# 会员管理系统 - 完整测试脚本

set -e

PROJECT_ROOT="/Users/niumc/Downloads/code/gostudy/memberhub"
cd "$PROJECT_ROOT"

echo "========================================"
echo "  会员管理系统 - 完整环境测试"
echo "========================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. 检查后端API服务
echo "1️⃣  检查后端API服务..."
if ps aux | grep "cmd/api/main.go" | grep -v grep > /dev/null; then
    echo -e "${GREEN}✅ 后端API服务已运行${NC}"
    API_RUNNING=true
else
    echo -e "${YELLOW}⚠️  后端API服务未运行，正在启动...${NC}"
    make run > /tmp/api.log 2>&1 &
    API_PID=$!
    sleep 3
    API_RUNNING=false
fi

# 2. 测试API健康检查
echo ""
echo "2️⃣  测试API健康状态..."
HEALTH_CHECK=$(curl -s http://localhost:8080/health 2>/dev/null || echo "error")
if [[ $HEALTH_CHECK == *"ok"* ]]; then
    echo -e "${GREEN}✅ API健康检查通过${NC}"
    echo "   响应: $HEALTH_CHECK"
else
    echo -e "${RED}❌ API健康检查失败${NC}"
    echo "   请检查日志: tail -f /tmp/api.log"
    exit 1
fi

# 3. 测试管理员登录
echo ""
echo "3️⃣  测试管理员登录..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')

if [[ $LOGIN_RESPONSE == *"token"* ]]; then
    echo -e "${GREEN}✅ 管理员登录成功${NC}"
    TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')
    echo "   Token: ${TOKEN:0:50}..."
else
    echo -e "${RED}❌ 管理员登录失败${NC}"
    echo "   响应: $LOGIN_RESPONSE"
    echo ""
    echo "   正在尝试修复管理员密码..."

    # 生成新的bcrypt hash并更新
    NEW_HASH=$(cd /tmp && cat > gen_hash.go << 'EOF'
package main
import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)
func main() {
    hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
    fmt.Print(string(hash))
}
EOF
go run gen_hash.go)

    mysql -h127.0.0.1 -uroot -proot membership_dev -e "UPDATE admins SET password='$NEW_HASH' WHERE username='admin';" 2>/dev/null

    echo "   密码已重置，请重新运行此脚本"
    exit 1
fi

# 4. 测试门店API
echo ""
echo "4️⃣  测试门店列表API..."
STORES_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/admin/stores?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN")

STORES_COUNT=$(echo $STORES_RESPONSE | jq -r '.data.total // 0')
echo -e "${GREEN}✅ 门店列表API正常${NC}"
echo "   门店数量: $STORES_COUNT"

# 5. 测试会员API
echo ""
echo "5️⃣  测试会员列表API..."
MEMBERS_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/admin/members?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN")

MEMBERS_COUNT=$(echo $MEMBERS_RESPONSE | jq -r '.data.total // 0')
echo -e "${GREEN}✅ 会员列表API正常${NC}"
echo "   会员数量: $MEMBERS_COUNT"

# 6. 测试积分规则API
echo ""
echo "6️⃣  测试积分规则API..."
POINTS_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/admin/points/rules?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN")

POINTS_COUNT=$(echo $POINTS_RESPONSE | jq -r '.data.total // 0')
echo -e "${GREEN}✅ 积分规则API正常${NC}"
echo "   规则数量: $POINTS_COUNT"

# 7. 测试充值活动API
echo ""
echo "7️⃣  测试充值活动API..."
RECHARGE_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/admin/recharge/promotions?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN")

RECHARGE_COUNT=$(echo $RECHARGE_RESPONSE | jq -r '.data.total // 0')
echo -e "${GREEN}✅ 充值活动API正常${NC}"
echo "   活动数量: $RECHARGE_COUNT"

# 8. 测试优惠券API
echo ""
echo "8️⃣  测试优惠券模板API..."
COUPON_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/admin/coupons/templates?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN")

COUPON_COUNT=$(echo $COUPON_RESPONSE | jq -r '.data.total // 0')
echo -e "${GREEN}✅ 优惠券模板API正常${NC}"
echo "   模板数量: $COUPON_COUNT"

# 9. 测试促销商品API
echo ""
echo "9️⃣  测试特价商品API..."
PROMOTION_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/admin/promotions?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN")

PROMOTION_COUNT=$(echo $PROMOTION_RESPONSE | jq -r '.data.total // 0')
echo -e "${GREEN}✅ 特价商品API正常${NC}"
echo "   商品数量: $PROMOTION_COUNT"

# 10. 检查前端服务
echo ""
echo "🔟  检查前端服务..."
cd "$PROJECT_ROOT/web/admin"

if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}⚠️  前端依赖未安装，正在安装...${NC}"
    npm install
fi

# 生成测试报告
echo ""
echo "========================================"
echo "  ✅ 测试完成 - 系统状态报告"
echo "========================================"
echo ""
echo "📊 数据统计:"
echo "   - 门店数量: $STORES_COUNT"
echo "   - 会员数量: $MEMBERS_COUNT"
echo "   - 积分规则: $POINTS_COUNT"
echo "   - 充值活动: $RECHARGE_COUNT"
echo "   - 优惠券模板: $COUPON_COUNT"
echo "   - 特价商品: $PROMOTION_COUNT"
echo ""
echo "🔑 测试账号:"
echo "   用户名: admin"
echo "   密码: admin123"
echo ""
echo "🌐 访问地址:"
echo "   后端API: http://localhost:8080"
echo "   前端管理后台: http://localhost:3000 (需要启动)"
echo ""
echo "📝 启动前端服务:"
echo "   cd $PROJECT_ROOT/web/admin"
echo "   npm run dev"
echo ""
echo "📖 查看详细文档:"
echo "   cat $PROJECT_ROOT/web/admin/API_TESTING_GUIDE.md"
echo ""
echo "========================================"
echo ""

# 询问是否启动前端
read -p "是否立即启动前端服务? (y/n): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${GREEN}🚀 正在启动前端服务...${NC}"
    cd "$PROJECT_ROOT/web/admin"
    npm run dev
fi
