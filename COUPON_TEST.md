# 优惠券系统测试指南

本文档提供优惠券系统的完整测试方案，包括功能测试、并发测试和数据一致性验证。

## 目录
- [环境准备](#环境准备)
- [功能测试](#功能测试)
- [并发测试](#并发测试)
- [数据一致性验证](#数据一致性验证)
- [常见问题排查](#常见问题排查)

---

## 环境准备

### 1. 启动服务

```bash
# 启动依赖服务
docker-compose -f deployments/docker/docker-compose.yml up -d mysql redis rabbitmq

# 运行数据库迁移
migrate -path migrations -database "mysql://root:root@tcp(localhost:3306)/membership" up

# 启动API服务
go run cmd/api/main.go
```

### 2. 准备测试数据

```sql
-- 插入测试优惠券模板
INSERT INTO coupon_templates (
    name, type, discount_type, discount_value, min_amount,
    total_quantity, remaining_quantity, per_user_limit, valid_days,
    status, sort_order, created_at, updated_at
) VALUES
(
    '满100减20优惠券', 1, 1, 2000, 10000,
    1000, 1000, 2, 30,
    1, 100, NOW(), NOW()
),
(
    '8折优惠券', 2, 2, 800, 5000,
    500, 500, 1, 15,
    1, 90, NOW(), NOW()
),
(
    '10元代金券', 3, 1, 1000, 0,
    2000, 2000, 5, 60,
    1, 80, NOW(), NOW()
);

-- 查看插入的模板ID
SELECT id, name, type, total_quantity, remaining_quantity, per_user_limit FROM coupon_templates;
```

### 3. 获取测试Token

```bash
# 小程序用户登录
curl -X POST http://localhost:8080/api/miniapp/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "code": "TEST_WECHAT_CODE"
  }'

# 保存返回的token
export USER_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 管理员登录
curl -X POST http://localhost:8080/api/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'

# 保存管理员token
export ADMIN_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 功能测试

### 1. 小程序端 - 优惠券模板列表

查看当前可领取的优惠券模板。

```bash
curl -X GET "http://localhost:8080/api/miniapp/coupons/templates" \
  -H "Authorization: Bearer $USER_TOKEN"
```

**预期响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "满100减20优惠券",
      "type": 1,
      "discount_type": 1,
      "discount_value": 2000,
      "min_amount": 10000,
      "total_quantity": 1000,
      "remaining_quantity": 1000,
      "per_user_limit": 2,
      "valid_days": 30,
      "status": 1,
      "sort_order": 100
    }
  ]
}
```

**验证点:**
- ✅ 只返回 `status=1` 且 `remaining_quantity > 0` 的模板
- ✅ 只返回在时间范围内的模板
- ✅ 按 `sort_order DESC, id DESC` 排序

---

### 2. 小程序端 - 领取优惠券

用户领取优惠券。

```bash
curl -X POST http://localhost:8080/api/miniapp/coupons/receive \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": 1
  }'
```

**预期响应:**
```json
{
  "code": 0,
  "message": "coupon received successfully",
  "data": {
    "id": 1,
    "coupon_code": "CPN20260315AB1234",
    "user_id": 1,
    "template_id": 1,
    "name": "满100减20优惠券",
    "type": 1,
    "discount_type": 1,
    "discount_value": 2000,
    "min_amount": 10000,
    "status": 1,
    "expired_at": "2026-04-14T10:30:00Z",
    "received_at": "2026-03-15T10:30:00Z"
  }
}
```

**验证点:**
- ✅ 优惠券码唯一（格式：CPN+日期+随机字符）
- ✅ `expired_at` = 当前时间 + `valid_days` 天
- ✅ `status` = 1（未使用）
- ✅ 数据库 `coupon_templates.remaining_quantity` 减1

**数据库验证:**
```sql
-- 检查模板库存
SELECT id, name, remaining_quantity FROM coupon_templates WHERE id = 1;
-- remaining_quantity 应该为 999

-- 检查用户优惠券
SELECT id, coupon_code, user_id, template_id, status, expired_at
FROM user_coupons WHERE user_id = 1 AND template_id = 1;
```

---

### 3. 小程序端 - 领取限制测试

测试用户领取数量限制。

```bash
# 第1次领取（成功）
curl -X POST http://localhost:8080/api/miniapp/coupons/receive \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"template_id": 1}'

# 第2次领取（成功，per_user_limit=2）
curl -X POST http://localhost:8080/api/miniapp/coupons/receive \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"template_id": 1}'

# 第3次领取（失败，超过限制）
curl -X POST http://localhost:8080/api/miniapp/coupons/receive \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"template_id": 1}'
```

**预期第3次响应:**
```json
{
  "code": 40011,
  "message": "coupon limit exceeded"
}
```

---

### 4. 小程序端 - 库存耗尽测试

测试库存不足的情况。

```bash
# 先更新库存为0
mysql> UPDATE coupon_templates SET remaining_quantity = 0 WHERE id = 1;

# 尝试领取
curl -X POST http://localhost:8080/api/miniapp/coupons/receive \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"template_id": 1}'
```

**预期响应:**
```json
{
  "code": 40012,
  "message": "coupon sold out"
}
```

---

### 5. 小程序端 - 查询可用优惠券

查询用户当前可用的优惠券（未使用且未过期）。

```bash
curl -X GET "http://localhost:8080/api/miniapp/coupons/available?page=1&page_size=20" \
  -H "Authorization: Bearer $USER_TOKEN"
```

**预期响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "coupon_code": "CPN20260315AB1234",
        "name": "满100减20优惠券",
        "status": 1,
        "expired_at": "2026-04-14T10:30:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

**验证点:**
- ✅ 只返回 `status=1` 的优惠券
- ✅ 按 `status ASC, expired_at ASC` 排序（即将过期的排前面）

---

### 6. 小程序端 - 查询所有优惠券

查询用户所有优惠券（包括已使用、已过期）。

```bash
curl -X GET "http://localhost:8080/api/miniapp/coupons/all?page=1&page_size=20" \
  -H "Authorization: Bearer $USER_TOKEN"
```

**预期响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "coupon_code": "CPN20260315AB1234",
        "status": 1,
        "expired_at": "2026-04-14T10:30:00Z"
      },
      {
        "id": 2,
        "coupon_code": "CPN20260310XY5678",
        "status": 2,
        "used_at": "2026-03-12T15:20:00Z"
      }
    ],
    "total": 2,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 7. 管理端 - 优惠券模板列表

管理员查看所有优惠券模板。

```bash
curl -X GET "http://localhost:8080/api/admin/coupons/templates?page=1&page_size=20" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**预期响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "满100减20优惠券",
        "type": 1,
        "total_quantity": 1000,
        "remaining_quantity": 998,
        "status": 1
      }
    ],
    "total": 3,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 8. 管理端 - 创建优惠券模板

管理员创建新的优惠券模板。

```bash
curl -X POST http://localhost:8080/api/admin/coupons/templates \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "新春特惠券",
    "type": 1,
    "discount_type": 1,
    "discount_value": 5000,
    "min_amount": 20000,
    "total_quantity": 500,
    "per_user_limit": 1,
    "valid_days": 7,
    "description": "新春活动专属优惠券",
    "status": 1,
    "sort_order": 100
  }'
```

**预期响应:**
```json
{
  "code": 0,
  "message": "template created successfully",
  "data": {
    "id": 4,
    "name": "新春特惠券",
    "total_quantity": 500,
    "remaining_quantity": 500
  }
}
```

**验证点:**
- ✅ `remaining_quantity` 自动初始化为 `total_quantity`
- ✅ 创建时间和更新时间自动设置

---

### 9. 管理端 - 更新优惠券模板

管理员更新现有模板。

```bash
curl -X PUT http://localhost:8080/api/admin/coupons/templates \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 4,
    "name": "新春特惠券（升级版）",
    "type": 1,
    "discount_type": 1,
    "discount_value": 6000,
    "min_amount": 20000,
    "total_quantity": 500,
    "remaining_quantity": 500,
    "per_user_limit": 2,
    "valid_days": 10,
    "status": 1,
    "sort_order": 100
  }'
```

**预期响应:**
```json
{
  "code": 0,
  "message": "template updated successfully"
}
```

---

### 10. 管理端 - 创建推送任务

管理员创建批量推送任务。

**推送给指定用户:**
```bash
curl -X POST http://localhost:8080/api/admin/coupons/push \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": 4,
    "title": "新春优惠券发放",
    "target_type": 2,
    "target_user_ids": [1, 2, 3, 4, 5],
    "send_wechat_msg": true
  }'
```

**预期响应:**
```json
{
  "code": 0,
  "message": "push task created successfully",
  "data": {
    "id": 1,
    "task_no": "CPT20260315001",
    "template_id": 4,
    "title": "新春优惠券发放",
    "target_type": 2,
    "total_count": 5,
    "status": 1,
    "created_at": "2026-03-15T10:30:00Z"
  }
}
```

**验证点:**
- ✅ `task_no` 自动生成（格式：CPT+日期+序号）
- ✅ `total_count` = `target_user_ids` 数组长度
- ✅ `status` = 1（待执行）
- ✅ `target_users` 字段存储JSON格式的用户ID列表

**数据库验证:**
```sql
SELECT id, task_no, template_id, target_type, target_users, total_count, status
FROM coupon_push_tasks WHERE id = 1;

-- target_users 应该为: "[1,2,3,4,5]"
```

---

### 11. 管理端 - 推送任务列表

查看所有推送任务。

```bash
curl -X GET "http://localhost:8080/api/admin/coupons/push-tasks?page=1&page_size=20" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**预期响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "task_no": "CPT20260315001",
        "title": "新春优惠券发放",
        "total_count": 5,
        "success_count": 0,
        "fail_count": 0,
        "status": 1
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

---

## 并发测试

### 1. 并发领取优惠券测试

测试多个用户同时领取同一模板的优惠券，验证库存扣减正确性。

#### 准备测试脚本

创建 `scripts/test_concurrent_receive.sh`:

```bash
#!/bin/bash

TEMPLATE_ID=1
CONCURRENT_NUM=100

# 重置库存
mysql -uroot -proot membership -e "UPDATE coupon_templates SET remaining_quantity = 50 WHERE id = $TEMPLATE_ID"

# 并发领取
for i in $(seq 1 $CONCURRENT_NUM); do
  {
    curl -X POST http://localhost:8080/api/miniapp/coupons/receive \
      -H "Authorization: Bearer $USER_TOKEN" \
      -H "Content-Type: application/json" \
      -d "{\"template_id\": $TEMPLATE_ID}" \
      2>/dev/null
  } &
done

wait

echo "测试完成！"

# 检查结果
mysql -uroot -proot membership -e "
  SELECT
    (SELECT remaining_quantity FROM coupon_templates WHERE id = $TEMPLATE_ID) as remaining_stock,
    (SELECT COUNT(*) FROM user_coupons WHERE template_id = $TEMPLATE_ID) as issued_count
"
```

#### 执行测试

```bash
chmod +x scripts/test_concurrent_receive.sh
./scripts/test_concurrent_receive.sh
```

#### 验证结果

```sql
-- 检查库存
SELECT id, name, total_quantity, remaining_quantity
FROM coupon_templates WHERE id = 1;

-- 检查实际发放数量
SELECT COUNT(*) as issued_count
FROM user_coupons WHERE template_id = 1;

-- 验证等式
-- remaining_quantity + issued_count = total_quantity
```

**预期结果:**
- ✅ `remaining_quantity` + 发放数量 = 初始库存（50）
- ✅ 无超发现象
- ✅ 所有失败请求返回 "库存不足" 或 "超过限领数量"

---

### 2. 并发重复领取测试

测试同一用户并发领取同一模板，验证分布式锁防重效果。

```bash
#!/bin/bash

TEMPLATE_ID=1
USER_TOKEN="your_user_token"

# 同一用户并发发起10次领取请求
for i in $(seq 1 10); do
  {
    curl -X POST http://localhost:8080/api/miniapp/coupons/receive \
      -H "Authorization: Bearer $USER_TOKEN" \
      -H "Content-Type: application/json" \
      -d "{\"template_id\": $TEMPLATE_ID}"
  } &
done

wait

# 检查该用户实际领取数量
mysql -uroot -proot membership -e "
  SELECT COUNT(*) as received_count
  FROM user_coupons
  WHERE user_id = 1 AND template_id = $TEMPLATE_ID
"
```

**预期结果:**
- ✅ 该用户实际领取数量 ≤ `per_user_limit`
- ✅ 无重复领取

---

## 数据一致性验证

### 1. 库存一致性检查

验证模板库存与实际发放数量一致。

```sql
SELECT
    ct.id,
    ct.name,
    ct.total_quantity,
    ct.remaining_quantity,
    COUNT(uc.id) as issued_count,
    (ct.total_quantity - ct.remaining_quantity) as should_issued,
    CASE
        WHEN (ct.total_quantity - ct.remaining_quantity) = COUNT(uc.id)
        THEN 'OK'
        ELSE 'INCONSISTENT'
    END as status
FROM coupon_templates ct
LEFT JOIN user_coupons uc ON ct.id = uc.template_id
GROUP BY ct.id;
```

**预期结果:**
- ✅ `status` 列全部为 'OK'
- ✅ `total_quantity - remaining_quantity = issued_count`

---

### 2. 优惠券码唯一性检查

验证所有优惠券码唯一。

```sql
-- 查找重复的优惠券码
SELECT coupon_code, COUNT(*) as count
FROM user_coupons
GROUP BY coupon_code
HAVING count > 1;
```

**预期结果:**
- ✅ 返回空结果（无重复）

---

### 3. 过期优惠券标记

验证过期优惠券自动标记功能（定时任务）。

```sql
-- 手动将某些优惠券的过期时间设为过去
UPDATE user_coupons
SET expired_at = '2026-03-01 00:00:00'
WHERE id IN (1, 2, 3) AND status = 1;

-- 调用过期标记方法（通常由定时任务调用）
-- 这里需要在代码中暴露一个管理接口或手动执行
-- UPDATE user_coupons SET status = 3 WHERE status = 1 AND expired_at <= NOW();

-- 检查结果
SELECT id, coupon_code, status, expired_at
FROM user_coupons
WHERE id IN (1, 2, 3);
```

**预期结果:**
- ✅ `status` 应该从 1（未使用）变为 3（已过期）

---

## 常见问题排查

### 问题1: 领取优惠券失败 - "coupon sold out"

**可能原因:**
1. 模板库存为0
2. Redis缓存库存与数据库不一致

**排查步骤:**
```sql
-- 检查模板库存
SELECT id, name, remaining_quantity, status
FROM coupon_templates WHERE id = ?;

-- 检查Redis缓存（如果使用了缓存）
redis-cli GET "coupon:stock:1"
```

**解决方案:**
- 如果库存确实为0，增加库存或创建新模板
- 如果Redis缓存不一致，清除缓存让其重新加载

---

### 问题2: 并发领取出现超发

**可能原因:**
1. 分布式锁未生效
2. 数据库原子操作失败

**排查步骤:**
```bash
# 检查Redis连接
redis-cli PING

# 查看错误日志
tail -f logs/app.log | grep "coupon"
```

**解决方案:**
- 确保Redis正常运行
- 检查分布式锁的TTL设置（建议5-10秒）
- 验证数据库更新使用了原子操作（UPDATE ... WHERE remaining_quantity > 0）

---

### 问题3: 用户超过限领数量仍能领取

**可能原因:**
1. 限领检查逻辑失效
2. 并发请求绕过了检查

**排查步骤:**
```sql
-- 检查用户实际领取数量
SELECT user_id, template_id, COUNT(*) as received_count
FROM user_coupons
GROUP BY user_id, template_id
HAVING received_count > (
    SELECT per_user_limit FROM coupon_templates WHERE id = template_id
);
```

**解决方案:**
- 在分布式锁内再次检查用户领取数量
- 在数据库层面添加唯一索引（如果业务允许）

---

### 问题4: 推送任务创建失败

**可能原因:**
1. 模板不存在或已停用
2. 指定用户列表为空（target_type=2时）

**排查步骤:**
```sql
-- 检查模板状态
SELECT id, name, status FROM coupon_templates WHERE id = ?;
```

**解决方案:**
- 确保模板 `status=1`（启用状态）
- `target_type=2` 时必须提供 `target_user_ids`

---

## 性能基准测试

### 使用 wrk 进行压测

```bash
# 安装wrk
brew install wrk  # macOS
# 或
sudo apt-get install wrk  # Linux

# 创建测试脚本 receive.lua
cat > receive.lua << 'EOF'
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Bearer YOUR_TOKEN"
wrk.body = '{"template_id": 1}'
EOF

# 执行压测（10线程，100并发，持续30秒）
wrk -t10 -c100 -d30s -s receive.lua http://localhost:8080/api/miniapp/coupons/receive
```

**性能目标:**
- QPS > 500（领取接口）
- 平均响应时间 < 100ms
- 99%请求响应时间 < 500ms
- 无数据不一致

---

## 总结

优惠券系统的关键测试点：

1. **功能完整性**
   - ✅ 模板管理
   - ✅ 领取限制
   - ✅ 库存控制
   - ✅ 推送任务

2. **并发安全性**
   - ✅ 分布式锁防重
   - ✅ 原子减库存
   - ✅ 无超发现象

3. **数据一致性**
   - ✅ 库存 = 总量 - 发放数
   - ✅ 优惠券码唯一
   - ✅ 用户领取数 ≤ 限领数

4. **性能表现**
   - ✅ 高并发下响应快速
   - ✅ 数据库压力可控
   - ✅ Redis缓存命中率高

通过以上测试，可以确保优惠券系统在生产环境下稳定可靠运行。
