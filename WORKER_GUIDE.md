# RabbitMQ Worker 系统指南

本文档提供 RabbitMQ Worker 系统的完整说明，包括功能介绍、部署运行和测试方法。

## 目录
- [系统架构](#系统架构)
- [Worker 列表](#worker-列表)
- [部署运行](#部署运行)
- [消息队列设计](#消息队列设计)
- [测试指南](#测试指南)

---

## 系统架构

Worker 系统基于 RabbitMQ 消息队列实现异步任务处理，采用生产者-消费者模式：

```
API Service (生产者)
    ↓
RabbitMQ (消息队列)
    ↓
Worker Service (消费者)
```

**核心特性：**
- 异步处理：不阻塞主业务流程
- 削峰填谷：平滑处理突发流量
- 失败重试：自动重试机制
- 限流保护：防止下游系统过载
- 消息持久化：保证消息不丢失

---

## Worker 列表

### 1. CouponPushWorker（优惠券推送Worker）

**功能：** 批量发放优惠券给用户

**触发方式：** 管理员创建优惠券推送任务时触发

**消息队列：** `queue.coupon.push`

**处理流程：**
1. 从队列中获取推送任务ID
2. 查询推送任务详情
3. 更新任务状态为"执行中"
4. 解析目标用户列表
5. 批量发放优惠券
6. 发送微信模板消息（可选）
7. 更新任务完成状态

**关键代码：**
```go
type CouponPushMessage struct {
    TaskID int64 `json:"task_id"`
}

// 批量发放优惠券
for _, userID := range targetUserIDs {
    _, err := couponService.IssueCouponToUser(ctx, userID, templateID)
    if err != nil {
        failCount++
        continue
    }
    successCount++
}
```

**性能：**
- 每批处理100条记录打印一次进度
- 支持数千用户并发发放
- 失败自动记录，不影响其他用户

---

### 2. WechatNotifyWorker（微信通知Worker）

**功能：** 发送微信模板消息通知

**触发方式：** 系统事件触发（优惠券发放、积分变动等）

**消息队列：** `queue.notify.wechat`

**处理流程：**
1. 从队列中获取通知消息
2. 限流控制（20条/秒）
3. 调用微信API发送模板消息
4. 失败重试（最多3次）
5. 确认消息

**关键代码：**
```go
type WechatNotifyMessage struct {
    UserID     int64                  `json:"user_id"`
    TemplateID string                 `json:"template_id"`
    Data       map[string]interface{} `json:"data"`
    Page       string                 `json:"page,omitempty"`
}

// 限流器（20条/秒）
rateLimiter := time.NewTicker(50 * time.Millisecond)
<-rateLimiter.C

// 发送模板消息
sendTemplateMessage(msg)
```

**限流机制：**
- 使用 Ticker 控制发送速率
- 20条/秒（符合微信限制）
- 自动排队等待

**重试机制：**
- 失败自动重试3次
- 使用消息头记录重试次数
- 超过次数丢弃消息并记录日志

---

### 3. PointsExpireWorker（积分过期Worker）

**功能：** 定时处理过期积分

**触发方式：** 定时任务（每天凌晨2点）

**处理流程：**
1. 计算下次执行时间
2. 延迟到凌晨2点
3. 查询所有过期积分
4. 批量标记为已过期
5. 记录过期流水
6. 等待24小时后重复

**关键代码：**
```go
// 计算下次执行时间（明天凌晨2点）
now := time.Now()
next := time.Date(now.Year(), now.Month(), now.Day()+1, 2, 0, 0, 0, now.Location())
duration := next.Sub(now)

// 首次延迟
time.Sleep(duration)

// 每24小时执行一次
ticker := time.NewTicker(24 * time.Hour)
for {
    processExpiredPoints()
    <-ticker.C
}
```

**执行时间：**
- 首次运行：下一个凌晨2点
- 后续运行：每天凌晨2点
- 选择凌晨2点的原因：业务低峰期

---

## 部署运行

### 1. 环境要求

- RabbitMQ 3.8+
- MySQL 8.0+
- Redis 7+
- Go 1.25+

### 2. 配置RabbitMQ

确保 RabbitMQ 交换机和队列已创建（API服务启动时自动创建）：

**交换机：**
- `exchange.coupon` (Topic)
- `exchange.notify` (Topic)
- `exchange.points` (Topic)

**队列：**
- `queue.coupon.push` - 绑定到 `exchange.coupon`，routing key: `coupon.push`
- `queue.notify.wechat` - 绑定到 `exchange.notify`，routing key: `notify.wechat`

### 3. 启动Worker服务

**开发环境：**
```bash
# 使用配置文件
CONFIG_PATH=configs/config.dev.yaml go run cmd/worker/main.go

# 或使用Make
make worker
```

**生产环境：**
```bash
# 编译
go build -o bin/worker cmd/worker/main.go

# 运行
CONFIG_PATH=configs/config.yaml ./bin/worker
```

**Docker部署：**
```yaml
# docker-compose.yml
services:
  worker:
    build: .
    command: /app/worker
    environment:
      - CONFIG_PATH=/app/configs/config.yaml
    depends_on:
      - mysql
      - redis
      - rabbitmq
    restart: always
```

启动：
```bash
docker-compose up -d worker
```

### 4. 查看运行状态

**查看日志：**
```bash
# 本地运行
tail -f logs/worker.log

# Docker运行
docker-compose logs -f worker
```

**预期输出：**
```
2026-03-15 10:00:00 INFO Worker service starting...
2026-03-15 10:00:00 INFO Starting CouponPushWorker
2026-03-15 10:00:00 INFO CouponPushWorker started successfully
2026-03-15 10:00:00 INFO Starting WechatNotifyWorker
2026-03-15 10:00:00 INFO WechatNotifyWorker started successfully
2026-03-15 10:00:00 INFO Starting PointsExpireWorker
2026-03-15 10:00:00 INFO PointsExpireWorker scheduled next_run=2026-03-16 02:00:00
2026-03-15 10:00:00 INFO All workers started successfully
```

---

## 消息队列设计

### 1. 交换机（Exchange）

| 交换机名称 | 类型 | 用途 |
|-----------|------|------|
| `exchange.coupon` | Topic | 优惠券相关事件 |
| `exchange.notify` | Topic | 通知相关事件 |
| `exchange.points` | Topic | 积分相关事件 |

### 2. 队列（Queue）

| 队列名称 | 绑定交换机 | Routing Key | 持久化 |
|---------|-----------|-------------|--------|
| `queue.coupon.push` | exchange.coupon | coupon.push | 是 |
| `queue.notify.wechat` | exchange.notify | notify.wechat | 是 |

### 3. 消息格式

**优惠券推送消息：**
```json
{
  "task_id": 123
}
```

**微信通知消息：**
```json
{
  "user_id": 1,
  "template_id": "coupon_received",
  "data": {
    "title": "新优惠券",
    "content": "您有一张新的优惠券"
  },
  "page": "pages/coupon/list"
}
```

### 4. 消息可靠性保证

**生产者确认：**
```go
// 发布消息时启用确认模式
channel.Confirm(false)

// 等待确认
if confirmed := <-confirms; !confirmed.Ack {
    // 消息发送失败，处理重试
}
```

**消费者确认：**
```go
// 手动ACK
d.Ack(false)  // 确认消息

// 拒绝并重新入队
d.Nack(false, true)

// 拒绝并丢弃
d.Nack(false, false)
```

**消息持久化：**
```go
// 发布持久化消息
channel.Publish(
    exchange,
    routingKey,
    false,
    false,
    amqp.Publishing{
        DeliveryMode: amqp.Persistent,  // 持久化
        ContentType:  "application/json",
        Body:         body,
    },
)
```

---

## 测试指南

### 1. 测试优惠券推送Worker

**步骤1: 创建优惠券模板**
```bash
curl -X POST http://localhost:8080/api/admin/coupons/templates \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试优惠券",
    "type": 1,
    "discount_type": 1,
    "discount_value": 1000,
    "min_amount": 0,
    "total_quantity": 1000,
    "per_user_limit": 5,
    "valid_days": 30,
    "status": 1
  }'
```

**步骤2: 创建推送任务**
```bash
curl -X POST http://localhost:8080/api/admin/coupons/push \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": 1,
    "title": "测试推送",
    "target_type": 2,
    "target_user_ids": [1, 2, 3, 4, 5],
    "send_wechat_msg": true
  }'
```

**步骤3: 查看Worker日志**
```bash
tail -f logs/worker.log
```

**预期日志：**
```
2026-03-15 10:05:00 INFO Processing coupon push task task_id=1
2026-03-15 10:05:01 INFO Coupon push progress task_id=1 success_count=5 fail_count=0 total=5
2026-03-15 10:05:01 INFO Coupon push task completed task_id=1 success_count=5 fail_count=0
```

**步骤4: 验证结果**
```bash
# 查询推送任务状态
curl -X GET "http://localhost:8080/api/admin/coupons/push-tasks/1" \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# 查询用户优惠券
curl -X GET "http://localhost:8080/api/miniapp/coupons/available" \
  -H "Authorization: Bearer $USER_TOKEN"
```

**预期结果：**
- 任务状态为"已完成" (status=3)
- success_count = 5
- 用户可以看到新优惠券

---

### 2. 测试微信通知Worker

**步骤1: 发布通知消息**

在代码中调用：
```go
import "review_demo/internal/infrastructure/messaging/rabbitmq"

// 构造消息
msg := WechatNotifyMessage{
    UserID:     1,
    TemplateID: "test_template",
    Data: map[string]interface{}{
        "title":   "测试标题",
        "content": "测试内容",
    },
}

body, _ := json.Marshal(msg)

// 发布消息
rabbitmq.Publish("exchange.notify", "notify.wechat", body)
```

**步骤2: 查看Worker日志**
```bash
tail -f logs/worker.log
```

**预期日志：**
```
2026-03-15 10:10:00 INFO Sending wechat notification user_id=1 template_id=test_template
2026-03-15 10:10:00 INFO Template message sent (mock) user_id=1
2026-03-15 10:10:00 INFO Wechat notification sent successfully user_id=1
```

---

### 3. 测试积分过期Worker

**步骤1: 插入过期积分测试数据**
```sql
-- 插入一些已过期的积分（过期时间设为昨天）
INSERT INTO points_transactions (user_id, points, type, expired_at, created_at)
VALUES
  (1, 100, 1, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 30 DAY)),
  (2, 200, 1, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 30 DAY));
```

**步骤2: 等待Worker执行**

Worker会在下一个凌晨2点执行，或者修改代码立即执行：
```go
// 临时修改执行时间为当前时间
next := time.Now().Add(5 * time.Second)
```

**步骤3: 查看日志和数据**
```bash
# 查看日志
tail -f logs/worker.log

# 查询过期积分
SELECT * FROM points_transactions WHERE expired_at < NOW();
```

---

### 4. 监控RabbitMQ队列

**使用Management UI：**
```bash
# 访问 RabbitMQ 管理界面
open http://localhost:15672

# 默认账号: guest / guest
```

**使用命令行：**
```bash
# 查看队列状态
rabbitmqctl list_queues name messages_ready messages_unacknowledged

# 查看交换机
rabbitmqctl list_exchanges name type

# 查看绑定关系
rabbitmqctl list_bindings
```

---

## 性能调优

### 1. 并发消费

增加消费者数量提高处理速度：

```yaml
# docker-compose.yml
services:
  worker1:
    build: .
    command: /app/worker
  worker2:
    build: .
    command: /app/worker
  worker3:
    build: .
    command: /app/worker
```

### 2. 预取数量

设置预取数量控制并发：
```go
// 每次预取10条消息
channel.Qos(10, 0, false)
```

### 3. 批量处理

批量发送通知消息：
```go
// 批量发布消息
for i := 0; i < 1000; i++ {
    rabbitmq.Publish("exchange.notify", "notify.wechat", body)
}
```

---

## 故障排查

### 问题1: Worker无法启动

**可能原因：**
- RabbitMQ连接失败
- 配置文件错误

**排查步骤：**
```bash
# 检查RabbitMQ是否运行
docker ps | grep rabbitmq

# 检查连接
telnet localhost 5672

# 查看日志
tail -f logs/worker.log
```

---

### 问题2: 消息堆积

**可能原因：**
- Worker处理速度慢
- Worker进程崩溃

**排查步骤：**
```bash
# 查看队列长度
rabbitmqctl list_queues name messages

# 增加Worker数量
docker-compose up -d --scale worker=3
```

---

### 问题3: 消息丢失

**可能原因：**
- 消息未持久化
- 消费者未正确ACK

**解决方案：**
- 确保消息持久化
- 使用手动ACK模式
- 启用生产者确认

---

## 总结

Worker系统的关键点：

1. **异步处理**
   - ✅ 优惠券批量发放
   - ✅ 微信模板消息
   - ✅ 积分定时过期

2. **可靠性保证**
   - ✅ 消息持久化
   - ✅ 失败重试
   - ✅ 手动ACK

3. **性能优化**
   - ✅ 限流控制
   - ✅ 并发处理
   - ✅ 批量操作

4. **监控告警**
   - ✅ 日志记录
   - ✅ 队列监控
   - ✅ 错误追踪

通过Worker系统，系统可以处理大规模异步任务，保证业务的高性能和高可用性。
