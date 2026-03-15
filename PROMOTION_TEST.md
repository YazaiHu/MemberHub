# 特价商品管理系统测试指南

本文档提供特价商品管理系统的完整测试方案，包括功能测试、时间范围验证和业务逻辑测试。

## 目录
- [环境准备](#环境准备)
- [功能测试](#功能测试)
- [业务逻辑测试](#业务逻辑测试)
- [数据验证](#数据验证)

---

## 环境准备

### 1. 启动服务

```bash
# 启动依赖服务
docker-compose -f deployments/docker/docker-compose.yml up -d mysql redis

# 运行数据库迁移
migrate -path migrations -database "mysql://root:root@tcp(localhost:3306)/membership" up

# 启动API服务
go run cmd/api/main.go
```

### 2. 获取测试Token

```bash
# 管理员登录
curl -X POST http://localhost:8080/api/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'

# 保存管理员token
export ADMIN_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 小程序用户登录
curl -X POST http://localhost:8080/api/miniapp/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "code": "TEST_WECHAT_CODE"
  }'

# 保存用户token
export USER_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 功能测试

### 1. 管理端 - 创建特价商品

管理员创建新的特价商品。

```bash
curl -X POST http://localhost:8080/api/admin/promotions \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "store_id": 1,
    "name": "苹果iPhone 15特价",
    "description": "最新款iPhone，限量特价促销",
    "original_price": 699900,
    "promotion_price": 599900,
    "stock": 50,
    "image_url": "https://example.com/iphone15.jpg",
    "start_time": "2026-03-15 00:00:00",
    "end_time": "2026-03-22 23:59:59",
    "status": 1,
    "sort_order": 100
  }'
```

**预期响应:**
```json
{
  "code": 0,
  "message": "promotion product created successfully",
  "data": {
    "id": 1,
    "store_id": 1,
    "name": "苹果iPhone 15特价",
    "description": "最新款iPhone，限量特价促销",
    "original_price": 699900,
    "promotion_price": 599900,
    "stock": 50,
    "image_url": "https://example.com/iphone15.jpg",
    "start_time": "2026-03-15T00:00:00Z",
    "end_time": "2026-03-22T23:59:59Z",
    "status": 1,
    "sort_order": 100,
    "created_at": "2026-03-15T10:00:00Z",
    "updated_at": "2026-03-15T10:00:00Z"
  }
}
```

**验证点:**
- ✅ 商品ID自动生成
- ✅ 时间正确解析
- ✅ 创建时间和更新时间自动设置
- ✅ 折扣价低于原价

**数据库验证:**
```sql
SELECT * FROM promotion_products WHERE id = 1;
```

---

### 2. 参数验证测试

**测试1: 折扣价高于原价（应失败）**
```bash
curl -X POST http://localhost:8080/api/admin/promotions \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "store_id": 1,
    "name": "错误商品",
    "original_price": 10000,
    "promotion_price": 15000,
    "stock": 10,
    "start_time": "2026-03-15 00:00:00",
    "end_time": "2026-03-22 23:59:59",
    "status": 1
  }'
```

**预期响应:**
```json
{
  "code": 40001,
  "message": "promotion_price must be less than original_price"
}
```

**测试2: 结束时间早于开始时间（应失败）**
```bash
curl -X POST http://localhost:8080/api/admin/promotions \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "store_id": 1,
    "name": "错误商品",
    "original_price": 10000,
    "promotion_price": 8000,
    "stock": 10,
    "start_time": "2026-03-22 00:00:00",
    "end_time": "2026-03-15 23:59:59",
    "status": 1
  }'
```

**预期响应:**
```json
{
  "code": 40001,
  "message": "end_time must be after start_time"
}
```

**测试3: 缺少必填字段（应失败）**
```bash
curl -X POST http://localhost:8080/api/admin/promotions \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试商品",
    "original_price": 10000
  }'
```

**预期响应:**
```json
{
  "code": 40001,
  "message": "invalid params: ..."
}
```

---

### 3. 管理端 - 特价商品列表

查看所有特价商品（带分页和筛选）。

```bash
# 查看所有特价商品
curl -X GET "http://localhost:8080/api/admin/promotions?page=1&page_size=20" \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# 按门店筛选
curl -X GET "http://localhost:8080/api/admin/promotions?store_id=1&page=1&page_size=20" \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# 按状态筛选
curl -X GET "http://localhost:8080/api/admin/promotions?status=1&page=1&page_size=20" \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# 组合筛选
curl -X GET "http://localhost:8080/api/admin/promotions?store_id=1&status=1&page=1&page_size=20" \
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
        "store_id": 1,
        "name": "苹果iPhone 15特价",
        "original_price": 699900,
        "promotion_price": 599900,
        "stock": 50,
        "start_time": "2026-03-15T00:00:00Z",
        "end_time": "2026-03-22T23:59:59Z",
        "status": 1
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

**验证点:**
- ✅ 返回分页数据
- ✅ 按 `sort_order DESC, start_time DESC, id DESC` 排序
- ✅ 门店和状态筛选生效

---

### 4. 管理端 - 特价商品详情

查看单个商品的详细信息。

```bash
curl -X GET "http://localhost:8080/api/admin/promotions/1" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**预期响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "store_id": 1,
    "name": "苹果iPhone 15特价",
    "description": "最新款iPhone，限量特价促销",
    "original_price": 699900,
    "promotion_price": 599900,
    "stock": 50,
    "image_url": "https://example.com/iphone15.jpg",
    "start_time": "2026-03-15T00:00:00Z",
    "end_time": "2026-03-22T23:59:59Z",
    "status": 1,
    "sort_order": 100,
    "created_at": "2026-03-15T10:00:00Z",
    "updated_at": "2026-03-15T10:00:00Z"
  }
}
```

---

### 5. 管理端 - 更新特价商品

更新商品信息（如调整价格、增加库存）。

```bash
curl -X PUT http://localhost:8080/api/admin/promotions \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "store_id": 1,
    "name": "苹果iPhone 15特价（加量不加价）",
    "description": "最新款iPhone，限量特价促销，库存增加",
    "original_price": 699900,
    "promotion_price": 589900,
    "stock": 100,
    "image_url": "https://example.com/iphone15.jpg",
    "start_time": "2026-03-15 00:00:00",
    "end_time": "2026-03-25 23:59:59",
    "status": 1,
    "sort_order": 100
  }'
```

**预期响应:**
```json
{
  "code": 0,
  "message": "promotion product updated successfully"
}
```

**验证点:**
- ✅ 所有字段更新成功
- ✅ `updated_at` 时间戳更新

**数据库验证:**
```sql
SELECT name, promotion_price, stock, end_time, updated_at
FROM promotion_products WHERE id = 1;
```

---

### 6. 管理端 - 删除特价商品（软删除）

删除特价商品（实际是将状态改为停用）。

```bash
curl -X DELETE "http://localhost:8080/api/admin/promotions/1" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**预期响应:**
```json
{
  "code": 0,
  "message": "promotion product deleted successfully"
}
```

**验证点:**
- ✅ 商品记录仍存在
- ✅ `status` 字段变为 0（停用）

**数据库验证:**
```sql
SELECT id, name, status FROM promotion_products WHERE id = 1;

-- 结果应该是：
-- +----+----------------------+--------+
-- | id | name                 | status |
-- +----+----------------------+--------+
-- |  1 | 苹果iPhone 15特价... |      0 |
-- +----+----------------------+--------+
```

---

### 7. 小程序端 - 本周特价商品

用户查看本周特价商品（7天内的特价）。

**准备测试数据:**
```bash
# 创建本周特价（今天开始）
curl -X POST http://localhost:8080/api/admin/promotions \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "store_id": 1,
    "name": "本周特价-iPad",
    "original_price": 399900,
    "promotion_price": 349900,
    "stock": 30,
    "start_time": "2026-03-15 00:00:00",
    "end_time": "2026-03-20 23:59:59",
    "status": 1,
    "sort_order": 90
  }'

# 创建下周特价（8天后开始）
curl -X POST http://localhost:8080/api/admin/promotions \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "store_id": 1,
    "name": "下周特价-MacBook",
    "original_price": 999900,
    "promotion_price": 899900,
    "stock": 20,
    "start_time": "2026-03-25 00:00:00",
    "end_time": "2026-03-30 23:59:59",
    "status": 1,
    "sort_order": 80
  }'
```

**查询本周特价:**
```bash
# 不指定门店（返回所有门店）
curl -X GET "http://localhost:8080/api/miniapp/promotions/weekly" \
  -H "Authorization: Bearer $USER_TOKEN"

# 指定门店
curl -X GET "http://localhost:8080/api/miniapp/promotions/weekly?store_id=1" \
  -H "Authorization: Bearer $USER_TOKEN"
```

**预期响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 2,
      "store_id": 1,
      "name": "本周特价-iPad",
      "original_price": 399900,
      "promotion_price": 349900,
      "stock": 30,
      "start_time": "2026-03-15T00:00:00Z",
      "end_time": "2026-03-20T23:59:59Z",
      "status": 1
    }
  ]
}
```

**验证点:**
- ✅ 只返回7天内的特价（不包括8天后的）
- ✅ 只返回启用且有库存的商品
- ✅ 按 `sort_order DESC, start_time ASC` 排序

---

### 8. 小程序端 - 当前有效特价商品

用户查看当前正在进行的特价商品。

```bash
curl -X GET "http://localhost:8080/api/miniapp/promotions?store_id=1" \
  -H "Authorization: Bearer $USER_TOKEN"
```

**预期响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 2,
      "store_id": 1,
      "name": "本周特价-iPad",
      "original_price": 399900,
      "promotion_price": 349900,
      "stock": 30,
      "start_time": "2026-03-15T00:00:00Z",
      "end_time": "2026-03-20T23:59:59Z",
      "status": 1
    }
  ]
}
```

**验证点:**
- ✅ 只返回当前时间在 `start_time` 和 `end_time` 之间的商品
- ✅ 只返回启用且有库存的商品
- ✅ 未开始的特价不在结果中
- ✅ 已过期的特价不在结果中

---

## 业务逻辑测试

### 1. 时间范围测试

验证特价商品的时间范围逻辑。

**测试数据:**
```sql
-- 插入3个测试商品
INSERT INTO promotion_products (store_id, name, original_price, promotion_price, stock, start_time, end_time, status, created_at, updated_at)
VALUES
  -- 未开始（明天开始）
  (1, '未开始特价', 10000, 8000, 10, DATE_ADD(NOW(), INTERVAL 1 DAY), DATE_ADD(NOW(), INTERVAL 3 DAY), 1, NOW(), NOW()),
  -- 进行中（昨天开始，明天结束）
  (1, '进行中特价', 20000, 16000, 20, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_ADD(NOW(), INTERVAL 1 DAY), 1, NOW(), NOW()),
  -- 已过期（3天前开始，昨天结束）
  (1, '已过期特价', 30000, 25000, 30, DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), 1, NOW(), NOW());
```

**测试当前有效特价:**
```bash
curl -X GET "http://localhost:8080/api/miniapp/promotions" \
  -H "Authorization: Bearer $USER_TOKEN"
```

**预期结果:**
- ✅ 只返回"进行中特价"
- ✅ "未开始特价"不在结果中
- ✅ "已过期特价"不在结果中

---

### 2. 库存为零测试

验证库存为0的商品不显示。

```sql
-- 将商品库存设为0
UPDATE promotion_products SET stock = 0 WHERE id = 2;
```

**查询有效特价:**
```bash
curl -X GET "http://localhost:8080/api/miniapp/promotions" \
  -H "Authorization: Bearer $USER_TOKEN"
```

**预期结果:**
- ✅ 库存为0的商品不在结果中

**恢复库存:**
```sql
UPDATE promotion_products SET stock = 30 WHERE id = 2;
```

---

### 3. 状态过滤测试

验证停用的商品不显示在用户端。

```sql
-- 停用商品
UPDATE promotion_products SET status = 0 WHERE id = 2;
```

**查询有效特价:**
```bash
curl -X GET "http://localhost:8080/api/miniapp/promotions" \
  -H "Authorization: Bearer $USER_TOKEN"
```

**预期结果:**
- ✅ 停用的商品不在结果中

**恢复状态:**
```sql
UPDATE promotion_products SET status = 1 WHERE id = 2;
```

---

### 4. 折扣计算测试

验证商品的折扣百分比和节省金额计算。

**测试商品:**
- 原价：6999.00元（699900分）
- 特价：5999.00元（599900分）

**预期计算结果:**
- 折扣百分比：85.7% ≈ 86折
- 节省金额：1000.00元（100000分）

```go
// 在代码中验证
product := &PromotionProduct{
    OriginalPrice:  699900,
    PromotionPrice: 599900,
}

discountPercent := product.DiscountPercent()
// 应该返回 85（85.7%向下取整）

saveAmount := product.SaveAmount()
// 应该返回 100000（分）
```

---

## 数据验证

### 1. 排序验证

验证特价商品列表按正确顺序排列。

```sql
-- 设置不同的排序权重
UPDATE promotion_products SET sort_order = 100 WHERE id = 1;
UPDATE promotion_products SET sort_order = 90 WHERE id = 2;
UPDATE promotion_products SET sort_order = 80 WHERE id = 3;
```

**请求商品列表:**
```bash
curl -X GET "http://localhost:8080/api/admin/promotions" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.data.list[].id'
```

**预期输出:**
```
1
2
3
```

（按 `sort_order DESC, start_time DESC, id DESC` 排序）

---

### 2. 时间格式验证

验证时间字段的格式和解析。

**正确的时间格式:**
```
"2026-03-15 00:00:00"
```

**错误的时间格式（应失败）:**
```bash
curl -X POST http://localhost:8080/api/admin/promotions \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "store_id": 1,
    "name": "测试商品",
    "original_price": 10000,
    "promotion_price": 8000,
    "stock": 10,
    "start_time": "2026-03-15",
    "end_time": "2026-03-22",
    "status": 1
  }'
```

**预期响应:**
```json
{
  "code": 40001,
  "message": "invalid start_time format"
}
```

---

## 常见问题排查

### 问题1: 特价商品不显示在用户端

**可能原因:**
1. 商品状态为停用
2. 库存为0
3. 时间范围不在当前时间内

**排查步骤:**
```sql
-- 检查商品状态
SELECT id, name, status, stock, start_time, end_time
FROM promotion_products
WHERE id = ?;

-- 检查当前时间
SELECT NOW();
```

---

### 问题2: 本周特价返回空数组

**可能原因:**
1. 所有特价都在7天之外
2. 所有特价都已过期
3. 所有特价都是停用状态

**排查步骤:**
```sql
-- 检查7天内的特价
SELECT id, name, start_time, end_time, status, stock
FROM promotion_products
WHERE start_time <= DATE_ADD(NOW(), INTERVAL 7 DAY)
  AND end_time >= NOW()
  AND status = 1
  AND stock > 0;
```

---

### 问题3: 创建特价商品失败

**常见错误及解决方案:**

| 错误信息 | 原因 | 解决方案 |
|---------|------|---------|
| `promotion_price must be less than original_price` | 折扣价≥原价 | 确保折扣价 < 原价 |
| `end_time must be after start_time` | 结束时间≤开始时间 | 确保结束时间 > 开始时间 |
| `invalid start_time format` | 时间格式错误 | 使用格式 `YYYY-MM-DD HH:MM:SS` |
| `store_id is required` | 缺少门店ID | 提供有效的门店ID |

---

## 总结

特价商品管理系统的关键测试点：

1. **基础CRUD**
   - ✅ 创建特价商品
   - ✅ 更新特价商品
   - ✅ 删除特价商品（软删除）
   - ✅ 查询商品列表
   - ✅ 查询商品详情

2. **参数验证**
   - ✅ 折扣价必须低于原价
   - ✅ 结束时间必须晚于开始时间
   - ✅ 时间格式验证
   - ✅ 必填字段验证

3. **时间范围逻辑**
   - ✅ 未开始的特价不显示
   - ✅ 进行中的特价正常显示
   - ✅ 已过期的特价不显示
   - ✅ 本周特价（7天内）

4. **业务规则**
   - ✅ 库存为0不显示
   - ✅ 停用状态不显示
   - ✅ 折扣计算正确
   - ✅ 按门店筛选

通过以上测试，可以确保特价商品管理系统功能完整且业务逻辑正确。
