# 数据库迁移文件

## 使用说明

### 方式1：使用 golang-migrate 工具

```bash
# 安装 migrate 工具
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 执行迁移
migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/membership?charset=utf8mb4&parseTime=True&loc=Local" up

# 回滚迁移
migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/membership?charset=utf8mb4&parseTime=True&loc=Local" down
```

### 方式2：直接执行 SQL 文件

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS membership CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 执行迁移文件
mysql -u root -p membership < migrations/001_create_users_and_admins.up.sql
mysql -u root -p membership < migrations/002_create_stores.up.sql
mysql -u root -p membership < migrations/003_create_points.up.sql
mysql -u root -p membership < migrations/004_create_recharge.up.sql
mysql -u root -p membership < migrations/005_create_coupons.up.sql
mysql -u root -p membership < migrations/006_create_promotions_and_logs.up.sql
```

## 数据库表结构

### 1. 用户与权限模块
- `users` - 用户表（微信小程序用户）
- `admins` - 管理员表
- `roles` - 角色表
- `permissions` - 权限表
- `role_permissions` - 角色权限关联表

### 2. 门店模块
- `stores` - 门店表

### 3. 积分模块
- `member_points` - 会员积分汇总表
- `points_transactions` - 积分流水表（数据源头）
- `points_exchange_rules` - 积分兑换规则表
- `points_exchange_records` - 积分兑换记录表

### 4. 充值模块
- `member_balance` - 会员余额汇总表
- `balance_transactions` - 余额流水表（数据源头）
- `recharge_promotions` - 充值活动表
- `recharge_orders` - 充值订单表

### 5. 优惠券模块
- `coupon_templates` - 优惠券模板表
- `user_coupons` - 用户优惠券表
- `coupon_push_tasks` - 优惠券推送任务表

### 6. 促销模块
- `promotion_products` - 特价商品表

### 7. 系统模块
- `operation_logs` - 操作日志表

## 默认数据

### 管理员账号
- 用户名：`admin`
- 密码：`admin123`
- 角色：超级管理员

### 默认角色
- `super_admin` - 超级管理员
- `store_admin` - 店铺管理员
- `member` - 普通会员

### 示例数据
- 2个示例门店
- 3个充值活动
- 3个优惠券模板
- 2个特价商品

## 关键设计说明

### 1. 乐观锁
在高并发场景使用 `version` 字段防止并发冲突：
- `member_points.version`
- `member_balance.version`
- `points_exchange_rules.version`

### 2. 流水表作为数据源
汇总表数据应该从流水表计算得出，保证数据一致性：
- `member_points` 从 `points_transactions` 计算
- `member_balance` 从 `balance_transactions` 计算

### 3. 唯一索引防重
使用唯一索引防止重复处理：
- `points_transactions.transaction_no`
- `balance_transactions.transaction_no`
- `recharge_orders.order_no`
- `user_coupons.coupon_code`

### 4. 索引优化
关键查询字段都建立了索引：
- 外键字段
- 状态字段
- 时间范围查询字段
- 排序字段

## 注意事项

1. **金额单位**：所有金额字段使用「分」作为单位存储，避免浮点数精度问题
2. **时间字段**：使用 DATETIME 类型，配合应用层时区处理
3. **JSON字段**：适用门店、目标用户等需要灵活存储的场景
4. **字符集**：使用 utf8mb4 支持 emoji 等特殊字符
5. **引擎**：使用 InnoDB 支持事务和外键

## 数据对账

定期执行以下查询检查数据一致性：

```sql
-- 积分对账
SELECT mp.user_id, mp.available_points,
       IFNULL(SUM(pt.points), 0) as transaction_points
FROM member_points mp
LEFT JOIN points_transactions pt ON mp.user_id = pt.user_id AND pt.status = 1
GROUP BY mp.user_id
HAVING mp.available_points != transaction_points;

-- 余额对账
SELECT mb.user_id, mb.balance,
       IFNULL(SUM(bt.amount), 0) as transaction_balance
FROM member_balance mb
LEFT JOIN balance_transactions bt ON mb.user_id = bt.user_id
GROUP BY mb.user_id
HAVING mb.balance != transaction_balance;
```
