# API 测试指南

## 环境准备

```bash
# 1. 启动依赖服务
make docker-up

# 2. 执行数据库迁移
make migrate-up

# 3. 启动API服务
make run
```

## 测试已实现的功能

### 1. 健康检查

```bash
curl http://localhost:8080/health
```

预期响应：
```json
{
  "status": "ok",
  "time": 1710508800
}
```

---

### 2. 管理员登录

```bash
curl -X POST http://localhost:8080/api/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

预期响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "admin_id": 1,
    "admin_info": {
      "id": 1,
      "username": "admin",
      "real_name": "超级管理员",
      "role": {
        "id": 1,
        "name": "超级管理员",
        "code": "super_admin"
      }
    }
  }
}
```

**保存返回的 token，后续请求需要使用！**

---

### 3. 会员列表

```bash
TOKEN="你的管理员token"

curl http://localhost:8080/api/admin/members \
  -H "Authorization: Bearer $TOKEN"
```

---

### 4. 创建积分兑换规则

```bash
TOKEN="你的管理员token"

curl -X POST http://localhost:8080/api/admin/points/rules \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "store_id": 1,
    "title": "100积分兑换咖啡",
    "description": "使用100积分兑换一杯美式咖啡",
    "points_cost": 100,
    "product_type": 1,
    "product_name": "美式咖啡",
    "total_stock": 100,
    "user_daily_limit": 1,
    "status": 1,
    "sort_order": 1
  }'
```

---

### 5. 查询兑换规则列表

```bash
curl http://localhost:8080/api/admin/points/rules \
  -H "Authorization: Bearer $TOKEN"
```

---

### 6. 手动调整用户积分

```bash
# 给用户ID=1的用户增加500积分
curl -X POST http://localhost:8080/api/admin/members/1/adjust-points \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "points": 500,
    "remark": "新用户注册奖励"
  }'
```

---

## 小程序端测试（需要真实微信code）

### 1. 微信登录

```bash
curl -X POST http://localhost:8080/api/miniapp/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "code": "WECHAT_CODE",
    "nickname": "测试用户",
    "avatar": "https://example.com/avatar.jpg",
    "gender": 1
  }'
```

**注意：** 需要真实的微信 code，或者修改代码暂时跳过微信验证。

---

### 2. 模拟微信登录（开发测试）

为了方便测试，可以暂时修改代码跳过微信验证：

修改 `internal/application/auth/service.go`:

```go
func (s *Service) WeChatLogin(ctx context.Context, req *WeChatLoginRequest) (*WeChatLoginResponse, error) {
    // 开发环境：模拟 openid
    openID := "test_openid_" + req.Code
    sessionKey := "test_session_key"
    unionID := ""

    // 正式环境取消注释
    // openID, sessionKey, unionID, err := wechat.Code2Session(req.Code)
    // if err != nil {
    //     return nil, errors.ErrWeChatAPIFailed.Wrap(err)
    // }

    // ... 其余代码不变
}
```

测试登录：

```bash
curl -X POST http://localhost:8080/api/miniapp/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "code": "test123",
    "nickname": "测试用户",
    "avatar": "https://example.com/avatar.jpg",
    "gender": 1
  }'
```

保存返回的用户 token。

---

### 3. 获取个人资料

```bash
USER_TOKEN="你的用户token"

curl http://localhost:8080/api/miniapp/profile \
  -H "Authorization: Bearer $USER_TOKEN"
```

---

### 4. 获取资产概览

```bash
curl http://localhost:8080/api/miniapp/asset \
  -H "Authorization: Bearer $USER_TOKEN"
```

---

### 5. 查询积分余额

```bash
curl http://localhost:8080/api/miniapp/points/balance \
  -H "Authorization: Bearer $USER_TOKEN"
```

---

### 6. 查询积分历史

```bash
curl http://localhost:8080/api/miniapp/points/history \
  -H "Authorization: Bearer $USER_TOKEN"
```

---

### 7. 查询兑换规则

```bash
curl http://localhost:8080/api/miniapp/points/exchange/rules \
  -H "Authorization: Bearer $USER_TOKEN"
```

---

### 8. 积分兑换

```bash
curl -X POST http://localhost:8080/api/miniapp/points/exchange \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "rule_id": 1
  }'
```

预期响应：
```json
{
  "code": 0,
  "message": "exchange successful",
  "data": {
    "id": 1,
    "record_no": "EX20260315120000123456",
    "user_id": 1,
    "rule_id": 1,
    "points_cost": 100,
    "product_name": "美式咖啡",
    "status": 1,
    "verify_code": "ABC12345"
  }
}
```

**保存核销码（verify_code）！**

---

### 9. 查询兑换记录

```bash
curl http://localhost:8080/api/miniapp/points/exchange/records \
  -H "Authorization: Bearer $USER_TOKEN"
```

---

### 10. 核销兑换（管理员）

```bash
curl -X POST http://localhost:8080/api/admin/points/exchange/verify \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "verify_code": "ABC12345"
  }'
```

---

## 完整测试流程示例

```bash
#!/bin/bash

# 1. 管理员登录
echo "=== 管理员登录 ==="
ADMIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}')

ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | jq -r '.data.token')
echo "管理员 Token: $ADMIN_TOKEN"
echo ""

# 2. 创建兑换规则
echo "=== 创建兑换规则 ==="
curl -s -X POST http://localhost:8080/api/admin/points/rules \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "store_id": 1,
    "title": "100积分兑换咖啡",
    "points_cost": 100,
    "product_type": 1,
    "product_name": "美式咖啡",
    "total_stock": 100,
    "user_daily_limit": 1,
    "status": 1,
    "sort_order": 1
  }' | jq '.'
echo ""

# 3. 用户登录（模拟）
echo "=== 用户登录 ==="
USER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/miniapp/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "code": "test123",
    "nickname": "测试用户",
    "gender": 1
  }')

USER_TOKEN=$(echo $USER_RESPONSE | jq -r '.data.token')
USER_ID=$(echo $USER_RESPONSE | jq -r '.data.user_id')
echo "用户 Token: $USER_TOKEN"
echo "用户 ID: $USER_ID"
echo ""

# 4. 管理员给用户加积分
echo "=== 给用户增加积分 ==="
curl -s -X POST http://localhost:8080/api/admin/members/$USER_ID/adjust-points \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "points": 500,
    "remark": "新用户奖励"
  }' | jq '.'
echo ""

# 5. 查询用户积分
echo "=== 查询用户积分 ==="
curl -s http://localhost:8080/api/miniapp/points/balance \
  -H "Authorization: Bearer $USER_TOKEN" | jq '.'
echo ""

# 6. 用户兑换
echo "=== 用户兑换 ==="
EXCHANGE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/miniapp/points/exchange \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"rule_id": 1}')

echo $EXCHANGE_RESPONSE | jq '.'
VERIFY_CODE=$(echo $EXCHANGE_RESPONSE | jq -r '.data.verify_code')
echo "核销码: $VERIFY_CODE"
echo ""

# 7. 管理员核销
echo "=== 管理员核销 ==="
curl -s -X POST http://localhost:8080/api/admin/points/exchange/verify \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"verify_code\": \"$VERIFY_CODE\"}" | jq '.'
echo ""

# 8. 查询兑换记录
echo "=== 查询兑换记录 ==="
curl -s http://localhost:8080/api/miniapp/points/exchange/records \
  -H "Authorization: Bearer $USER_TOKEN" | jq '.'
```

保存为 `test.sh`，执行：

```bash
chmod +x test.sh
./test.sh
```

---

## 注意事项

1. **数据库连接**：确保 `configs/config.dev.yaml` 中的数据库密码正确
2. **微信登录**：开发环境建议修改代码跳过微信验证
3. **Token过期**：Token 默认7天有效，过期需重新登录
4. **并发测试**：可以使用 `wrk` 或 `vegeta` 进行压力测试

---

## 性能测试

### 测试积分兑换并发性能

```bash
# 安装 wrk
brew install wrk  # macOS
# 或
apt-get install wrk  # Linux

# 准备测试脚本 exchange.lua
cat > exchange.lua << 'EOF'
wrk.method = "POST"
wrk.body   = '{"rule_id": 1}'
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Bearer YOUR_USER_TOKEN"
EOF

# 执行测试：100并发，持续30秒
wrk -t10 -c100 -d30s -s exchange.lua \
  http://localhost:8080/api/miniapp/points/exchange

# 查看结果验证无超卖
mysql -u root -p membership_dev -e "SELECT remaining_stock FROM points_exchange_rules WHERE id = 1;"
```

---

## 故障排查

### 1. 连接数据库失败

```bash
# 检查MySQL是否启动
docker ps | grep mysql

# 查看数据库日志
docker logs membership-mysql

# 测试连接
mysql -h 127.0.0.1 -P 3306 -u root -p
```

### 2. 连接Redis失败

```bash
# 检查Redis是否启动
docker ps | grep redis

# 测试连接
redis-cli ping
```

### 3. Token无效

- 检查token是否过期
- 检查 JWT secret 配置是否一致
- 重新登录获取新token

### 4. 权限不足

- 检查用户角色
- 超级管理员拥有所有权限
- 店铺管理员只能访问本店数据

---

## 下一步

当前系统已支持：
- ✅ 管理员和用户认证
- ✅ 会员管理
- ✅ 积分管理（含高并发兑换）

待实现：
- ⏳ 充值系统
- ⏳ 优惠券系统
- ⏳ 门店管理
- ⏳ Worker异步任务
