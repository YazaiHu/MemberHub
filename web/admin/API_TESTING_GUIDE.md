# 前后端API对接测试指南

## 1. 后端API服务状态

✅ **后端API已成功启动并运行**
- 服务地址: http://localhost:8080
- 数据库: MySQL (membership_dev)
- Redis: 已连接
- RabbitMQ: 已连接

## 2. 测试账号

### 管理员账号
- 用户名: `admin`
- 密码: `admin123`
- 角色: super_admin（超级管理员）

### 测试用户
系统中已有测试用户数据，可用于测试会员相关功能。

## 3. API端点验证

### ✅ 已验证的API

#### 3.1 管理员认证
```bash
# 登录
curl -X POST http://localhost:8080/api/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 返回示例:
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGci....",
    "admin_id": 1,
    "admin_info": {...}
  }
}
```

#### 3.2 门店管理
```bash
# 获取门店列表
curl http://localhost:8080/api/admin/stores?page=1&page_size=10 \
  -H "Authorization: Bearer YOUR_TOKEN"

# 返回示例:
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [...],
    "total": 2,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

### 📋 所有可用的管理后台API端点

#### 会员管理 (/api/admin/members)
- `GET /api/admin/members` - 会员列表
- `GET /api/admin/members/:id` - 会员详情
- `GET /api/admin/members/:id/asset` - 会员资产
- `POST /api/admin/members/:id/disable` - 禁用会员
- `POST /api/admin/members/:id/enable` - 启用会员
- `POST /api/admin/members/:id/adjust-points` - 调整积分
- `POST /api/admin/members/:id/adjust-balance` - 调整余额

#### 门店管理 (/api/admin/stores)
- `GET /api/admin/stores` - 门店列表
- `GET /api/admin/stores/:id` - 门店详情
- `POST /api/admin/stores` - 创建门店
- `PUT /api/admin/stores` - 更新门店
- `DELETE /api/admin/stores/:id` - 删除门店

#### 积分管理 (/api/admin/points)
- `GET /api/admin/points/rules` - 兑换规则列表
- `POST /api/admin/points/rules` - 创建兑换规则
- `PUT /api/admin/points/rules/:id` - 更新兑换规则
- `GET /api/admin/points/exchange/records` - 兑换记录列表
- `POST /api/admin/points/exchange/verify` - 核销兑换码

#### 充值管理 (/api/admin/recharge)
- `GET /api/admin/recharge/promotions` - 充值活动列表
- `POST /api/admin/recharge/promotions` - 创建充值活动
- `PUT /api/admin/recharge/promotions/:id` - 更新充值活动
- `GET /api/admin/recharge/orders` - 充值订单列表

#### 优惠券管理 (/api/admin/coupons)
- `GET /api/admin/coupons/templates` - 优惠券模板列表
- `GET /api/admin/coupons/templates/:id` - 模板详情
- `POST /api/admin/coupons/templates` - 创建模板
- `PUT /api/admin/coupons/templates` - 更新模板
- `POST /api/admin/coupons/push` - 创建推送任务
- `GET /api/admin/coupons/push-tasks` - 推送任务列表
- `GET /api/admin/coupons/push-tasks/:id` - 推送任务详情

#### 促销管理 (/api/admin/promotions)
- `GET /api/admin/promotions` - 特价商品列表
- `GET /api/admin/promotions/:id` - 商品详情
- `POST /api/admin/promotions` - 创建特价商品
- `PUT /api/admin/promotions` - 更新特价商品
- `DELETE /api/admin/promotions/:id` - 删除特价商品

## 4. 前端对接步骤

### Step 1: 启动后端服务（已完成）
```bash
cd /Users/niumc/Downloads/code/gostudy/review_demo
make run
```

### Step 2: 启动前端服务
```bash
cd web/admin
npm install
npm run dev
```

访问: http://localhost:3000

### Step 3: 测试登录流程

1. 打开浏览器访问 http://localhost:3000
2. 输入用户名: `admin`
3. 输入密码: `admin123`
4. 点击登录

**预期结果:**
- ✅ 成功跳转到仪表板页面
- ✅ 左侧显示完整的菜单列表
- ✅ Token保存在localStorage中

### Step 4: 测试各个功能模块

#### 4.1 仪表板
- 访问路径: `/dashboard`
- 功能: 查看统计数据、快捷操作

#### 4.2 会员管理
- 访问路径: `/members`
- 测试功能:
  - 查看会员列表
  - 搜索会员（手机号、昵称）
  - 查看会员详情
  - 调整会员积分
  - 调整会员余额

#### 4.3 积分管理
- 访问路径: `/points`
- 测试功能:
  - 查看兑换规则列表
  - 新增/编辑兑换规则
  - 查看兑换记录
  - 核销兑换码

#### 4.4 充值管理
- 访问路径: `/recharge`
- 测试功能:
  - 查看充值活动
  - 新增/编辑充值活动
  - 查看充值订单

#### 4.5 优惠券管理
- 访问路径: `/coupons`
- 测试功能:
  - 查看优惠券模板
  - 创建优惠券模板
  - 批量推送优惠券
  - 查看推送任务

#### 4.6 门店管理
- 访问路径: `/stores`
- 测试功能:
  - 查看门店列表
  - 新增/编辑门店
  - 配置经纬度
  - 删除门店

#### 4.7 促销管理
- 访问路径: `/promotions`
- 测试功能:
  - 查看特价商品
  - 新增/编辑特价商品
  - 设置促销时间

## 5. API请求格式说明

### 5.1 统一请求头
```javascript
headers: {
  'Content-Type': 'application/json',
  'Authorization': 'Bearer YOUR_JWT_TOKEN'
}
```

### 5.2 统一响应格式
```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 具体业务数据
  }
}
```

### 5.3 错误响应
```json
{
  "code": 400,
  "message": "错误信息"
}
```

### 5.4 分页参数
```javascript
// 请求参数
{
  "page": 1,
  "page_size": 10,
  // ...其他筛选参数
}

// 响应数据
{
  "list": [...],
  "total": 100,
  "page": 1,
  "page_size": 10,
  "total_pages": 10
}
```

## 6. 前端代码已完成的功能

### ✅ 基础设施
- Axios请求封装（request.js）
- 自动添加Token
- 统一错误处理
- 401自动跳转登录
- Pinia状态管理
- 路由守卫

### ✅ API模块
- `/api/auth.js` - 认证相关
- `/api/member.js` - 会员管理
- `/api/points.js` - 积分管理（新创建）
- `/api/coupon.js` - 优惠券管理
- `/api/store.js` - 门店管理

### ✅ 页面组件
- 登录页 ✅
- 仪表板 ✅
- 会员管理 ✅（已更新API调用）
- 积分管理 ✅（已更新API调用）
- 充值管理 ✅
- 优惠券管理 ✅
- 门店管理 ✅
- 促销管理 ✅

## 7. 常见问题排查

### Q1: 登录失败
**检查项:**
- 后端服务是否运行 (`ps aux | grep "cmd/api"`)
- 数据库连接是否正常
- 用户名密码是否正确（admin/admin123）
- 浏览器控制台是否有错误

### Q2: 接口401错误
**原因:** Token过期或无效
**解决:** 重新登录获取新Token

### Q3: 接口404错误
**原因:** API路径不正确
**检查:** 确认API路径是否匹配后端路由配置

### Q4: 数据为空
**原因:** 数据库中可能没有测试数据
**解决:**
```bash
# 可以通过管理后台UI创建测试数据
# 或者导入测试数据SQL
mysql -h127.0.0.1 -uroot -proot membership_dev < migrations/999_test_data.sql
```

### Q5: CORS跨域错误
**解决:** 后端已配置CORS中间件，前端Vite已配置代理

## 8. 下一步工作

由于前端页面较多，部分页面的API调用还保留了 `// TODO: 调用API` 注释，但基础架构已经完备，您可以按照以下模式更新：

### 更新模式示例
```javascript
// 旧代码（模拟数据）
const loadData = async () => {
  // TODO: 调用API
  data.value = [/* 模拟数据 */]
}

// 新代码（真实API）
const loadData = async () => {
  try {
    const result = await getDataAPI(params)
    data.value = result.list || []
    total.value = result.total || 0
  } catch (error) {
    ElMessage.error('加载失败')
  }
}
```

## 9. 完整测试清单

- [ ] 登录成功并获取Token
- [ ] 仪表板数据展示
- [ ] 会员列表查询
- [ ] 会员搜索功能
- [ ] 会员详情查看
- [ ] 调整会员积分
- [ ] 调整会员余额
- [ ] 门店列表查询
- [ ] 创建新门店
- [ ] 编辑门店信息
- [ ] 删除门店
- [ ] 积分规则列表
- [ ] 创建兑换规则
- [ ] 兑换记录查询
- [ ] 充值活动列表
- [ ] 创建充值活动
- [ ] 充值订单查询
- [ ] 优惠券模板列表
- [ ] 创建优惠券模板
- [ ] 批量推送优惠券
- [ ] 推送任务查询
- [ ] 特价商品列表
- [ ] 创建特价商品

## 10. 性能验证

### 响应时间基准
- 登录: < 500ms
- 列表查询: < 300ms
- 详情查询: < 200ms
- 创建/更新操作: < 500ms

### 并发测试
```bash
# 使用ab进行简单的并发测试
ab -n 100 -c 10 -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/admin/stores
```

---

**文档创建时间:** 2024-03-15
**API版本:** v1.0
**状态:** ✅ 后端API已就绪，前端对接进行中
