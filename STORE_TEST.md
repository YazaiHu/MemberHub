# 门店管理系统测试指南

本文档提供门店管理系统的完整测试方案，包括功能测试和附近门店查询测试。

## 目录
- [环境准备](#环境准备)
- [功能测试](#功能测试)
- [附近门店测试](#附近门店测试)
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

### 2. 查看测试数据

数据库迁移文件已自动插入3个测试门店：

```sql
-- 查看现有门店
SELECT id, name, address, latitude, longitude, status
FROM stores;

-- 结果示例：
-- +----+--------------+------------------+-----------+------------+--------+
-- | id | name         | address          | latitude  | longitude  | status |
-- +----+--------------+------------------+-----------+------------+--------+
-- |  1 | 北京朝阳店   | 北京市朝阳区xxx  | 39.921489 | 116.443108 |      1 |
-- |  2 | 上海浦东店   | 上海市浦东新区xxx| 31.222200 | 121.544270 |      1 |
-- |  3 | 深圳南山店   | 深圳市南山区xxx  | 22.539836 | 113.946582 |      1 |
-- +----+--------------+------------------+-----------+------------+--------+
```

### 3. 获取测试Token

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

### 1. 管理端 - 门店列表

查看所有门店（带分页和状态筛选）。

```bash
# 查看所有门店
curl -X GET "http://localhost:8080/api/admin/stores?page=1&page_size=20" \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# 只查看启用的门店
curl -X GET "http://localhost:8080/api/admin/stores?status=1&page=1&page_size=20" \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# 只查看停用的门店
curl -X GET "http://localhost:8080/api/admin/stores?status=0&page=1&page_size=20" \
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
        "name": "北京朝阳店",
        "address": "北京市朝阳区建国路88号",
        "phone": "010-12345678",
        "latitude": 39.921489,
        "longitude": 116.443108,
        "opening_hours": "09:00-22:00",
        "description": "北京朝阳区旗舰店",
        "status": 1,
        "sort_order": 100
      }
    ],
    "total": 3,
    "page": 1,
    "page_size": 20
  }
}
```

**验证点:**
- ✅ 返回分页数据
- ✅ 按 `sort_order DESC, id DESC` 排序
- ✅ 状态筛选生效

---

### 2. 管理端 - 门店详情

查看单个门店的详细信息。

```bash
curl -X GET "http://localhost:8080/api/admin/stores/1" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**预期响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "北京朝阳店",
    "address": "北京市朝阳区建国路88号",
    "phone": "010-12345678",
    "latitude": 39.921489,
    "longitude": 116.443108,
    "opening_hours": "09:00-22:00",
    "description": "北京朝阳区旗舰店",
    "status": 1,
    "sort_order": 100,
    "created_at": "2026-03-15T10:00:00Z",
    "updated_at": "2026-03-15T10:00:00Z"
  }
}
```

**错误测试 - 门店不存在:**
```bash
curl -X GET "http://localhost:8080/api/admin/stores/999" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**预期响应:**
```json
{
  "code": 40004,
  "message": "resource not found"
}
```

---

### 3. 管理端 - 创建门店

管理员创建新门店。

```bash
curl -X POST http://localhost:8080/api/admin/stores \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "广州天河店",
    "address": "广州市天河区天河路123号",
    "phone": "020-88888888",
    "latitude": 23.136692,
    "longitude": 113.324520,
    "opening_hours": "10:00-22:00",
    "description": "广州天河区新店",
    "status": 1,
    "sort_order": 90
  }'
```

**预期响应:**
```json
{
  "code": 0,
  "message": "store created successfully",
  "data": {
    "id": 4,
    "name": "广州天河店",
    "address": "广州市天河区天河路123号",
    "phone": "020-88888888",
    "latitude": 23.136692,
    "longitude": 113.324520,
    "opening_hours": "10:00-22:00",
    "description": "广州天河区新店",
    "status": 1,
    "sort_order": 90
  }
}
```

**验证点:**
- ✅ 门店ID自动生成
- ✅ 创建时间和更新时间自动设置
- ✅ 所有字段正确保存

**数据库验证:**
```sql
SELECT * FROM stores WHERE id = 4;
```

**参数验证测试 - 缺少必填字段:**
```bash
curl -X POST http://localhost:8080/api/admin/stores \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试店",
    "address": "测试地址"
  }'
```

**预期响应:**
```json
{
  "code": 40001,
  "message": "invalid params: latitude is required"
}
```

---

### 4. 管理端 - 更新门店

管理员更新门店信息。

```bash
curl -X PUT http://localhost:8080/api/admin/stores \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 4,
    "name": "广州天河店（旗舰店）",
    "address": "广州市天河区天河路123号",
    "phone": "020-99999999",
    "latitude": 23.136692,
    "longitude": 113.324520,
    "opening_hours": "09:00-23:00",
    "description": "广州天河区旗舰店，面积1000平米",
    "status": 1,
    "sort_order": 95
  }'
```

**预期响应:**
```json
{
  "code": 0,
  "message": "store updated successfully"
}
```

**验证点:**
- ✅ 所有字段更新成功
- ✅ `updated_at` 时间戳更新

**数据库验证:**
```sql
SELECT name, phone, opening_hours, sort_order, updated_at
FROM stores WHERE id = 4;
```

---

### 5. 管理端 - 删除门店（软删除）

删除门店（实际是将状态改为停用）。

```bash
curl -X DELETE "http://localhost:8080/api/admin/stores/4" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**预期响应:**
```json
{
  "code": 0,
  "message": "store deleted successfully"
}
```

**验证点:**
- ✅ 门店记录仍存在
- ✅ `status` 字段变为 0（停用）

**数据库验证:**
```sql
-- 门店应该仍然存在，但状态为0
SELECT id, name, status FROM stores WHERE id = 4;

-- 结果应该是：
-- +----+------------------+--------+
-- | id | name             | status |
-- +----+------------------+--------+
-- |  4 | 广州天河店(旗舰店) |      0 |
-- +----+------------------+--------+
```

---

### 6. 小程序端 - 门店列表

用户查看所有启用的门店。

```bash
curl -X GET "http://localhost:8080/api/miniapp/stores" \
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
      "name": "北京朝阳店",
      "address": "北京市朝阳区建国路88号",
      "phone": "010-12345678",
      "latitude": 39.921489,
      "longitude": 116.443108,
      "opening_hours": "09:00-22:00",
      "status": 1
    },
    {
      "id": 2,
      "name": "上海浦东店",
      "address": "上海市浦东新区世纪大道200号",
      "phone": "021-87654321",
      "latitude": 31.222200,
      "longitude": 121.544270,
      "opening_hours": "09:00-22:00",
      "status": 1
    },
    {
      "id": 3,
      "name": "深圳南山店",
      "address": "深圳市南山区科技园南区",
      "phone": "0755-12345678",
      "latitude": 22.539836,
      "longitude": 113.946582,
      "opening_hours": "09:00-22:00",
      "status": 1
    }
  ]
}
```

**验证点:**
- ✅ 只返回 `status=1` 的门店
- ✅ 按 `sort_order DESC, id DESC` 排序
- ✅ 已删除（停用）的门店不显示

---

## 附近门店测试

### 1. 附近门店查询（北京用户）

用户在北京，查询附近门店。

```bash
curl -X POST http://localhost:8080/api/miniapp/stores/nearby \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": 39.920000,
    "longitude": 116.440000,
    "radius": 10
  }'
```

**预期响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "store": {
        "id": 1,
        "name": "北京朝阳店",
        "address": "北京市朝阳区建国路88号",
        "phone": "010-12345678",
        "latitude": 39.921489,
        "longitude": 116.443108,
        "opening_hours": "09:00-22:00",
        "status": 1
      },
      "distance": 0.32
    }
  ]
}
```

**验证点:**
- ✅ 只返回指定半径内的门店
- ✅ 按距离从近到远排序
- ✅ 距离计算正确（单位：公里，保留两位小数）
- ✅ 其他城市的门店不在结果中

---

### 2. 附近门店查询（上海用户）

用户在上海，查询附近门店。

```bash
curl -X POST http://localhost:8080/api/miniapp/stores/nearby \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": 31.220000,
    "longitude": 121.540000,
    "radius": 5
  }'
```

**预期响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "store": {
        "id": 2,
        "name": "上海浦东店",
        "address": "上海市浦东新区世纪大道200号",
        "latitude": 31.222200,
        "longitude": 121.544270,
        "status": 1
      },
      "distance": 0.45
    }
  ]
}
```

---

### 3. 附近门店查询（无门店区域）

用户在偏远地区，附近无门店。

```bash
curl -X POST http://localhost:8080/api/miniapp/stores/nearby \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": 45.750000,
    "longitude": 126.650000,
    "radius": 10
  }'
```

**预期响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": []
}
```

---

### 4. 默认半径测试

不指定半径，使用默认值（10公里）。

```bash
curl -X POST http://localhost:8080/api/miniapp/stores/nearby \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": 39.920000,
    "longitude": 116.440000
  }'
```

**验证点:**
- ✅ `radius` 参数可选
- ✅ 默认半径为 10 公里

---

## 数据验证

### 1. 距离计算验证

验证 Haversine 公式计算的距离是否准确。

**测试数据:**
- 北京朝阳店: (39.921489, 116.443108)
- 用户位置: (39.920000, 116.440000)

**手动计算（在线工具验证）:**
使用在线经纬度距离计算器，输入上述坐标，结果应约为 **0.32 公里**。

**API响应验证:**
```bash
curl -X POST http://localhost:8080/api/miniapp/stores/nearby \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": 39.920000,
    "longitude": 116.440000,
    "radius": 10
  }' | jq '.data[0].distance'

# 输出应该约为 0.32
```

---

### 2. 门店状态过滤验证

确保停用的门店不会出现在用户端。

```sql
-- 停用一个门店
UPDATE stores SET status = 0 WHERE id = 1;
```

然后请求小程序端门店列表：

```bash
curl -X GET "http://localhost:8080/api/miniapp/stores" \
  -H "Authorization: Bearer $USER_TOKEN"
```

**验证点:**
- ✅ 停用的门店（id=1）不在结果中
- ✅ 其他启用门店正常显示

**恢复门店状态:**
```sql
UPDATE stores SET status = 1 WHERE id = 1;
```

---

### 3. 排序验证

验证门店列表按正确顺序排列。

```sql
-- 修改排序权重
UPDATE stores SET sort_order = 100 WHERE id = 1;
UPDATE stores SET sort_order = 90 WHERE id = 2;
UPDATE stores SET sort_order = 80 WHERE id = 3;
```

请求门店列表：

```bash
curl -X GET "http://localhost:8080/api/miniapp/stores" \
  -H "Authorization: Bearer $USER_TOKEN" | jq '.data[].id'
```

**预期输出:**
```
1
2
3
```

（按 `sort_order DESC, id DESC` 排序）

---

## 常见问题排查

### 问题1: 附近门店查询返回空数组

**可能原因:**
1. 半径太小
2. 附近确实无门店
3. 所有门店都是停用状态

**排查步骤:**
```sql
-- 检查门店状态
SELECT id, name, status FROM stores;

-- 检查经纬度是否合理
SELECT id, name, latitude, longitude FROM stores;
```

**解决方案:**
- 增大查询半径
- 确保门店状态为启用
- 检查经纬度坐标是否正确

---

### 问题2: 距离计算不准确

**可能原因:**
1. 经纬度数据录入错误
2. Haversine 公式实现有误

**排查步骤:**
```go
// 检查 calculateDistance 函数
distance := calculateDistance(39.921489, 116.443108, 39.920000, 116.440000)
fmt.Printf("Distance: %.2f km\n", distance)
// 应该输出约 0.32 km
```

**解决方案:**
- 使用在线工具验证经纬度
- 检查公式实现（已使用标准 Haversine 公式）

---

### 问题3: 创建门店失败 - 参数验证错误

**可能原因:**
1. 缺少必填字段
2. 经纬度格式错误
3. 状态值不合法

**正确的请求示例:**
```json
{
  "name": "门店名称",           // 必填
  "address": "详细地址",        // 必填
  "phone": "联系电话",          // 可选
  "latitude": 39.921489,       // 必填，float64
  "longitude": 116.443108,     // 必填，float64
  "opening_hours": "09:00-22:00", // 可选
  "description": "门店描述",    // 可选
  "status": 1,                 // 必填，0或1
  "sort_order": 100            // 可选，默认0
}
```

---

## 总结

门店管理系统的关键测试点：

1. **基础CRUD**
   - ✅ 创建门店
   - ✅ 更新门店
   - ✅ 删除门店（软删除）
   - ✅ 查询门店列表
   - ✅ 查询门店详情

2. **状态管理**
   - ✅ 停用门店不显示在用户端
   - ✅ 管理端可查看所有状态门店
   - ✅ 状态筛选功能

3. **附近门店**
   - ✅ Haversine 公式计算距离
   - ✅ 按距离排序
   - ✅ 半径过滤
   - ✅ 只返回启用门店

4. **排序功能**
   - ✅ 按 sort_order 降序
   - ✅ 相同权重按 id 降序

通过以上测试，可以确保门店管理系统功能完整且数据准确。
