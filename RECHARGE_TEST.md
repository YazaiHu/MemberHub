# 充值系统测试指南

## 功能概述

充值系统实现了：
- ✅ 充值活动管理（充值送活动）
- ✅ 充值订单创建
- ✅ 微信支付回调处理（带事务保证）
- ✅ 余额管理
- ✅ 余额流水查询
- ✅ 手动调整余额

## 核心特性

### 1. 事务原子性保证

充值回调处理保证以下操作的原子性：
```
订单状态更新 + 余额增加 + 积分赠送
```

关键代码逻辑：
```go
func ProcessPaymentCallback(orderNo, transactionID string, paidAmount int) error {
    // 1. 分布式锁防重
    lock := redis.AcquireLock("lock:recharge:callback:" + orderNo)
    defer lock.Release()

    // 2. 数据库事务
    tx.Begin()
    {
        // 2.1 幂等性检查（流水表）
        if existsTransaction(orderNo) {
            return nil // 已处理
        }

        // 2.2 更新订单状态
        order.Status = OrderStatusPaid
        order.TransactionID = transactionID
        order.PaidAt = now
        UpdateOrder(order)

        // 2.3 增加余额（乐观锁）
        balance.Balance += order.TotalAmount
        balance.TotalRecharge += order.RechargeAmount
        UpdateBalance(balance)

        // 2.4 记录余额流水
        CreateBalanceTransaction(...)

        // 2.5 赠送积分（如果有）
        if order.BonusPoints > 0 {
            AddPoints(user_id, bonus_points, ...)
        }
    }
    tx.Commit()
}
```

### 2. 幂等性保证

- 分布式锁防止并发处理
- 流水表唯一约束防止重复记录
- 订单状态检查

### 3. 并发安全

- 余额更新使用乐观锁（version字段）
- 分布式锁防止重复回调

---

## API 测试

### 1. 查询充值活动列表（小程序）

```bash
USER_TOKEN="你的用户token"

curl http://localhost:8080/api/miniapp/recharge/promotions \
  -H "Authorization: Bearer $USER_TOKEN"
```

预期响应：
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "title": "充100送10",
      "description": "充值100元送10元余额",
      "recharge_amount": 100,
      "bonus_amount": 10,
      "bonus_points": 0,
      "status": 1
    },
    {
      "id": 2,
      "title": "充200送50",
      "description": "充值200元送50元余额+100积分",
      "recharge_amount": 200,
      "bonus_amount": 50,
      "bonus_points": 100,
      "status": 1
    }
  ]
}
```

---

### 2. 创建充值订单（参加活动）

```bash
# 参加活动ID=1的充值活动
curl -X POST http://localhost:8080/api/miniapp/recharge/create-order \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "promotion_id": 1
  }'
```

预期响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order": {
      "id": 1,
      "order_no": "RC20260315120000123456",
      "user_id": 1,
      "promotion_id": 1,
      "recharge_amount": 10000,
      "bonus_amount": 1000,
      "bonus_points": 0,
      "total_amount": 11000,
      "pay_amount": 10000,
      "status": 1,
      "expired_at": "2026-03-15T12:30:00Z"
    },
    "pay_params": {
      "order_no": "RC20260315120000123456",
      "pay_amount": 10000
    }
  }
}
```

**保存 order_no，用于模拟支付回调！**

---

### 3. 创建充值订单（不参加活动）

```bash
# 直接充值50元
curl -X POST http://localhost:8080/api/miniapp/recharge/create-order \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 50
  }'
```

---

### 4. 模拟微信支付回调

```bash
# 模拟支付成功回调
ORDER_NO="RC20260315120000123456"

curl -X POST http://localhost:8080/api/callback/wechat-pay \
  -H "Content-Type: application/xml" \
  -d "<xml>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <appid><![CDATA[wx1234567890]]></appid>
  <mch_id><![CDATA[1234567890]]></mch_id>
  <openid><![CDATA[test_openid]]></openid>
  <total_fee>10000</total_fee>
  <transaction_id><![CDATA[WX20260315120000]]></transaction_id>
  <out_trade_no><![CDATA[$ORDER_NO]]></out_trade_no>
  <time_end><![CDATA[20260315120000]]></time_end>
</xml>"
```

预期响应：
```xml
<xml>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <return_msg><![CDATA[OK]]></return_msg>
</xml>
```

---

### 5. 查询用户余额

```bash
curl http://localhost:8080/api/miniapp/recharge/balance \
  -H "Authorization: Bearer $USER_TOKEN"
```

预期响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "user_id": 1,
    "balance": 11000,
    "frozen_balance": 0,
    "total_recharge": 10000,
    "total_consume": 0,
    "version": 1
  }
}
```

---

### 6. 查询余额历史

```bash
curl http://localhost:8080/api/miniapp/recharge/balance/history \
  -H "Authorization: Bearer $USER_TOKEN"
```

预期响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "transaction_no": "RC20260315120000123456",
        "user_id": 1,
        "amount": 11000,
        "type": 1,
        "source": "recharge",
        "remark": "充值：¥100.00",
        "created_at": "2026-03-15T12:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 7. 查询充值记录

```bash
curl http://localhost:8080/api/miniapp/recharge/records \
  -H "Authorization: Bearer $USER_TOKEN"
```

---

## 管理端测试

### 1. 创建充值活动

```bash
ADMIN_TOKEN="你的管理员token"

curl -X POST http://localhost:8080/api/admin/recharge/promotions \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "充500送150",
    "description": "充值500元送150元余额+300积分",
    "recharge_amount": 500,
    "bonus_amount": 150,
    "bonus_points": 300,
    "user_limit": 1,
    "status": 1,
    "sort_order": 1
  }'
```

---

### 2. 查询充值活动列表

```bash
curl http://localhost:8080/api/admin/recharge/promotions \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

---

### 3. 更新充值活动

```bash
curl -X PUT http://localhost:8080/api/admin/recharge/promotions/1 \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "充100送20（升级）",
    "description": "活动升级",
    "recharge_amount": 100,
    "bonus_amount": 20,
    "bonus_points": 50,
    "status": 1,
    "sort_order": 1
  }'
```

---

### 4. 手动调整用户余额

```bash
# 给用户ID=1增加5000分（50元）
curl -X POST http://localhost:8080/api/admin/members/1/adjust-balance \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 5000,
    "remark": "补偿余额"
  }'
```

---

### 5. 查询充值订单列表

```bash
# 查询所有已支付的订单
curl "http://localhost:8080/api/admin/recharge/orders?status=2" \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# 查询指定用户的订单
curl "http://localhost:8080/api/admin/recharge/orders?user_id=1" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

---

## 完整测试流程

```bash
#!/bin/bash

# 1. 管理员登录
ADMIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}')
ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | jq -r '.data.token')
echo "管理员 Token: $ADMIN_TOKEN"

# 2. 用户登录
USER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/miniapp/auth/login \
  -H "Content-Type: application/json" \
  -d '{"code": "test123", "nickname": "测试用户"}')
USER_TOKEN=$(echo $USER_RESPONSE | jq -r '.data.token')
USER_ID=$(echo $USER_RESPONSE | jq -r '.data.user_id')
echo "用户 Token: $USER_TOKEN"
echo "用户 ID: $USER_ID"

# 3. 查询充值活动
echo "=== 查询充值活动 ==="
curl -s http://localhost:8080/api/miniapp/recharge/promotions \
  -H "Authorization: Bearer $USER_TOKEN" | jq '.'

# 4. 创建充值订单（参加活动ID=1）
echo "=== 创建充值订单 ==="
ORDER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/miniapp/recharge/create-order \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"promotion_id": 1}')
echo $ORDER_RESPONSE | jq '.'
ORDER_NO=$(echo $ORDER_RESPONSE | jq -r '.data.order.order_no')
PAY_AMOUNT=$(echo $ORDER_RESPONSE | jq -r '.data.order.pay_amount')
echo "订单号: $ORDER_NO"
echo "支付金额: $PAY_AMOUNT"

# 5. 模拟微信支付回调
echo "=== 模拟支付回调 ==="
CALLBACK_XML="<xml>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <total_fee>$PAY_AMOUNT</total_fee>
  <transaction_id><![CDATA[WX$(date +%Y%m%d%H%M%S)]]></transaction_id>
  <out_trade_no><![CDATA[$ORDER_NO]]></out_trade_no>
</xml>"

curl -s -X POST http://localhost:8080/api/callback/wechat-pay \
  -H "Content-Type: application/xml" \
  -d "$CALLBACK_XML"
echo ""

# 6. 查询用户余额
echo "=== 查询用户余额 ==="
curl -s http://localhost:8080/api/miniapp/recharge/balance \
  -H "Authorization: Bearer $USER_TOKEN" | jq '.'

# 7. 查询余额历史
echo "=== 查询余额历史 ==="
curl -s http://localhost:8080/api/miniapp/recharge/balance/history \
  -H "Authorization: Bearer $USER_TOKEN" | jq '.'

# 8. 查询充值记录
echo "=== 查询充值记录 ==="
curl -s http://localhost:8080/api/miniapp/recharge/records \
  -H "Authorization: Bearer $USER_TOKEN" | jq '.'
```

---

## 数据验证

### 1. 验证余额和流水一致性

```sql
-- 查询余额汇总表和流水表是否一致
SELECT
    mb.user_id,
    mb.balance,
    IFNULL(SUM(bt.amount), 0) as transaction_balance
FROM member_balance mb
LEFT JOIN balance_transactions bt ON mb.user_id = bt.user_id
GROUP BY mb.user_id
HAVING mb.balance != transaction_balance;
```

应返回空结果。

---

### 2. 查看充值订单

```sql
SELECT * FROM recharge_orders ORDER BY created_at DESC LIMIT 10;
```

---

### 3. 查看余额流水

```sql
SELECT * FROM balance_transactions ORDER BY created_at DESC LIMIT 10;
```

---

## 并发测试

### 测试支付回调幂等性

```bash
#!/bin/bash

ORDER_NO="RC20260315120000123456"
PAY_AMOUNT=10000

# 并发发送10次相同的回调
for i in {1..10}; do
  curl -s -X POST http://localhost:8080/api/callback/wechat-pay \
    -H "Content-Type: application/xml" \
    -d "<xml>
      <return_code><![CDATA[SUCCESS]]></return_code>
      <result_code><![CDATA[SUCCESS]]></result_code>
      <total_fee>$PAY_AMOUNT</total_fee>
      <transaction_id><![CDATA[WX20260315120000]]></transaction_id>
      <out_trade_no><![CDATA[$ORDER_NO]]></out_trade_no>
    </xml>" &
done

wait

# 验证：余额只增加一次，流水表只有一条记录
mysql -u root -p membership_dev -e "
  SELECT user_id, balance FROM member_balance WHERE user_id = 1;
  SELECT COUNT(*) as count FROM balance_transactions WHERE transaction_no = '$ORDER_NO';
"
```

预期结果：
- 余额只增加一次
- 流水表只有一条记录

---

## 注意事项

1. **金额单位**：所有金额使用「分」，避免浮点数
2. **幂等性**：支付回调可能重复，必须保证幂等
3. **事务性**：订单+余额+积分必须在同一事务中
4. **签名验证**：生产环境必须验证微信签名
5. **超时处理**：订单30分钟后自动过期

---

## 故障排查

### 1. 回调处理失败

查看日志：
```bash
tail -f logs/app_dev.log | grep "payment callback"
```

### 2. 余额未增加

检查：
- 订单状态是否更新为已支付
- 余额流水表是否有记录
- 是否有并发冲突（乐观锁失败）

### 3. 重复回调

验证：
- 分布式锁是否生效
- 流水表唯一约束是否生效
- 订单状态检查是否生效
