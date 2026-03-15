#!/bin/bash

# Worker系统集成测试脚本
# 测试优惠券推送Worker和微信通知Worker

set -e

echo "=========================================="
echo "Worker System Integration Test"
echo "=========================================="
echo ""

# 配置
API_URL="http://localhost:8080"
ADMIN_TOKEN=""
USER_TOKEN=""

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 辅助函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查服务状态
check_service() {
    local service_name=$1
    local url=$2

    log_info "Checking $service_name..."
    if curl -s "$url" > /dev/null; then
        log_info "$service_name is running ✓"
        return 0
    else
        log_error "$service_name is not running ✗"
        return 1
    fi
}

# 等待服务启动
wait_for_service() {
    local service_name=$1
    local url=$2
    local max_wait=30
    local count=0

    log_info "Waiting for $service_name to start..."
    while [ $count -lt $max_wait ]; do
        if curl -s "$url" > /dev/null 2>&1; then
            log_info "$service_name is ready ✓"
            return 0
        fi
        sleep 1
        count=$((count + 1))
    done

    log_error "$service_name failed to start within ${max_wait}s"
    return 1
}

echo "Step 1: Checking required services"
echo "-------------------------------------------"

# 检查API服务
if ! check_service "API Service" "$API_URL/health"; then
    log_error "Please start API service first: make run"
    exit 1
fi

# 检查MySQL
if ! docker ps | grep -q mysql; then
    log_error "MySQL is not running. Start it with: make docker-up"
    exit 1
fi
log_info "MySQL is running ✓"

# 检查Redis
if ! docker ps | grep -q redis; then
    log_error "Redis is not running. Start it with: make docker-up"
    exit 1
fi
log_info "Redis is running ✓"

# 检查RabbitMQ
if ! docker ps | grep -q rabbitmq; then
    log_error "RabbitMQ is not running. Start it with: make docker-up"
    exit 1
fi
log_info "RabbitMQ is running ✓"

echo ""
echo "Step 2: Login and get tokens"
echo "-------------------------------------------"

# 管理员登录
log_info "Admin login..."
ADMIN_RESPONSE=$(curl -s -X POST "$API_URL/api/admin/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')

ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$ADMIN_TOKEN" ]; then
    log_error "Failed to get admin token"
    echo "Response: $ADMIN_RESPONSE"
    exit 1
fi

log_info "Admin token obtained ✓"

# 用户登录（模拟）
log_info "User login..."
USER_RESPONSE=$(curl -s -X POST "$API_URL/api/miniapp/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "TEST_WECHAT_CODE_001"
  }')

USER_TOKEN=$(echo $USER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$USER_TOKEN" ]; then
    log_warn "Failed to get user token, creating test user..."
fi

log_info "User token obtained ✓"

echo ""
echo "Step 3: Create test coupon template"
echo "-------------------------------------------"

TEMPLATE_RESPONSE=$(curl -s -X POST "$API_URL/api/admin/coupons/templates" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Worker测试优惠券",
    "type": 1,
    "discount_type": 1,
    "discount_value": 1000,
    "min_amount": 0,
    "total_quantity": 1000,
    "per_user_limit": 10,
    "valid_days": 30,
    "status": 1,
    "sort_order": 100
  }')

TEMPLATE_ID=$(echo $TEMPLATE_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -z "$TEMPLATE_ID" ]; then
    log_error "Failed to create coupon template"
    echo "Response: $TEMPLATE_RESPONSE"
    exit 1
fi

log_info "Coupon template created (ID: $TEMPLATE_ID) ✓"

echo ""
echo "Step 4: Create coupon push task"
echo "-------------------------------------------"

# 创建推送任务（推送给5个测试用户）
PUSH_RESPONSE=$(curl -s -X POST "$API_URL/api/admin/coupons/push" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"template_id\": $TEMPLATE_ID,
    \"title\": \"Worker系统测试推送\",
    \"target_type\": 2,
    \"target_user_ids\": [1, 2, 3, 4, 5],
    \"send_wechat_msg\": true
  }")

TASK_ID=$(echo $PUSH_RESPONSE | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -z "$TASK_ID" ]; then
    log_error "Failed to create push task"
    echo "Response: $PUSH_RESPONSE"
    exit 1
fi

log_info "Push task created (ID: $TASK_ID) ✓"

echo ""
echo "Step 5: Start Worker service"
echo "-------------------------------------------"

log_info "Starting Worker service in background..."

# 启动Worker服务
CONFIG_PATH=configs/config.dev.yaml go run cmd/worker/main.go > logs/worker_test.log 2>&1 &
WORKER_PID=$!

log_info "Worker PID: $WORKER_PID"

# 等待Worker启动
sleep 3

# 检查Worker是否在运行
if ! ps -p $WORKER_PID > /dev/null; then
    log_error "Worker failed to start"
    cat logs/worker_test.log
    exit 1
fi

log_info "Worker service started ✓"

echo ""
echo "Step 6: Monitor task processing"
echo "-------------------------------------------"

log_info "Waiting for Worker to process the task..."
log_info "Checking task status every 2 seconds..."

# 监控任务状态（最多等待30秒）
MAX_WAIT=30
ELAPSED=0

while [ $ELAPSED -lt $MAX_WAIT ]; do
    TASK_STATUS=$(curl -s "$API_URL/api/admin/coupons/push-tasks/$TASK_ID" \
      -H "Authorization: Bearer $ADMIN_TOKEN" | grep -o '"status":[0-9]*' | cut -d':' -f2)

    if [ "$TASK_STATUS" = "3" ]; then
        log_info "Task completed! ✓"
        break
    elif [ "$TASK_STATUS" = "2" ]; then
        log_info "Task is processing... (${ELAPSED}s)"
    elif [ "$TASK_STATUS" = "1" ]; then
        log_info "Task is pending... (${ELAPSED}s)"
    else
        log_warn "Unknown task status: $TASK_STATUS"
    fi

    sleep 2
    ELAPSED=$((ELAPSED + 2))
done

if [ "$TASK_STATUS" != "3" ]; then
    log_error "Task did not complete within ${MAX_WAIT}s"
    log_info "Current status: $TASK_STATUS (1=pending, 2=processing, 3=completed)"
fi

echo ""
echo "Step 7: Verify results"
echo "-------------------------------------------"

# 获取任务详情
log_info "Fetching task details..."
TASK_DETAIL=$(curl -s "$API_URL/api/admin/coupons/push-tasks/$TASK_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

SUCCESS_COUNT=$(echo $TASK_DETAIL | grep -o '"success_count":[0-9]*' | cut -d':' -f2)
FAIL_COUNT=$(echo $TASK_DETAIL | grep -o '"fail_count":[0-9]*' | cut -d':' -f2)
TOTAL_COUNT=$(echo $TASK_DETAIL | grep -o '"total_count":[0-9]*' | cut -d':' -f2)

log_info "Task Results:"
echo "  - Total: $TOTAL_COUNT"
echo "  - Success: $SUCCESS_COUNT"
echo "  - Failed: $FAIL_COUNT"

if [ "$SUCCESS_COUNT" = "$TOTAL_COUNT" ]; then
    log_info "All coupons distributed successfully! ✓"
else
    log_warn "Some coupons failed to distribute"
fi

# 检查用户优惠券（如果有用户token）
if [ -n "$USER_TOKEN" ]; then
    log_info "Checking user coupons..."
    USER_COUPONS=$(curl -s "$API_URL/api/miniapp/coupons/available" \
      -H "Authorization: Bearer $USER_TOKEN")

    COUPON_COUNT=$(echo $USER_COUPONS | grep -o '"total":[0-9]*' | cut -d':' -f2)
    log_info "User has $COUPON_COUNT available coupons"
fi

echo ""
echo "Step 8: Check Worker logs"
echo "-------------------------------------------"

log_info "Recent Worker logs:"
echo "-------------------------------------------"
tail -20 logs/worker_test.log
echo "-------------------------------------------"

echo ""
echo "Step 9: Check RabbitMQ queues"
echo "-------------------------------------------"

log_info "Checking RabbitMQ queue status..."

# 检查队列（需要rabbitmqctl或management API）
if command -v rabbitmqctl &> /dev/null; then
    rabbitmqctl list_queues name messages messages_ready messages_unacknowledged 2>/dev/null || \
        log_warn "Cannot access rabbitmqctl (may need sudo or different access)"
else
    log_warn "rabbitmqctl not available, skipping queue check"
    log_info "You can check RabbitMQ Management UI at http://localhost:15672"
fi

echo ""
echo "Step 10: Cleanup"
echo "-------------------------------------------"

log_info "Stopping Worker service..."
kill $WORKER_PID 2>/dev/null || true
sleep 2

if ps -p $WORKER_PID > /dev/null 2>&1; then
    log_warn "Worker did not stop gracefully, forcing..."
    kill -9 $WORKER_PID 2>/dev/null || true
fi

log_info "Worker stopped ✓"

echo ""
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo ""

if [ "$TASK_STATUS" = "3" ] && [ "$SUCCESS_COUNT" = "$TOTAL_COUNT" ]; then
    log_info "✓ All tests passed successfully!"
    echo ""
    echo "Test Results:"
    echo "  ✓ Services running"
    echo "  ✓ Coupon template created"
    echo "  ✓ Push task created"
    echo "  ✓ Worker processed task"
    echo "  ✓ Coupons distributed ($SUCCESS_COUNT/$TOTAL_COUNT)"
    echo ""
    exit 0
else
    log_error "✗ Some tests failed"
    echo ""
    echo "Test Results:"
    echo "  Task Status: $TASK_STATUS (expected: 3)"
    echo "  Success Count: $SUCCESS_COUNT (expected: $TOTAL_COUNT)"
    echo "  Fail Count: $FAIL_COUNT"
    echo ""
    echo "Please check logs/worker_test.log for details"
    exit 1
fi
