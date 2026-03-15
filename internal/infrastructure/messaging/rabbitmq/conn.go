package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"

	"github.com/YazaiHu/MemberHub/internal/pkg/config"
)

var conn *amqp.Connection
var ch *amqp.Channel

// 交换机定义
const (
	ExchangePoints   = "exchange.points"
	ExchangeCoupon   = "exchange.coupon"
	ExchangeNotify   = "exchange.notify"
	ExchangeRecharge = "exchange.recharge"
)

// 队列定义
const (
	QueuePointsSettle  = "queue.points.settle"
	QueuePointsExpire  = "queue.points.expire"
	QueueCouponPush    = "queue.coupon.push"
	QueueNotifyWechat  = "queue.notify.wechat"
	QueueNotifySMS     = "queue.notify.sms"
	QueueRechargeCallback = "queue.recharge.callback"
)

// 路由键定义
const (
	RoutingPointsEarned   = "points.earned"
	RoutingPointsExpired  = "points.expired"
	RoutingCouponPush     = "coupon.push"
	RoutingNotifyWechat   = "notify.wechat"
	RoutingNotifySMS      = "notify.sms"
	RoutingRechargeCallback = "recharge.callback"
)

// InitRabbitMQ 初始化RabbitMQ连接
func InitRabbitMQ(cfg *config.RabbitMQConfig) error {
	var err error

	// 连接配置
	amqpConfig := amqp.Config{
		Heartbeat: time.Duration(cfg.HeartbeatSecond) * time.Second,
	}

	// 建立连接
	conn, err = amqp.DialConfig(cfg.URL, amqpConfig)
	if err != nil {
		return fmt.Errorf("connect to rabbitmq failed: %w", err)
	}

	// 创建Channel
	ch, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel failed: %w", err)
	}

	// 设置消费者预取数量
	if err := ch.Qos(cfg.PrefetchCount, 0, false); err != nil {
		return fmt.Errorf("set qos failed: %w", err)
	}

	// 声明交换机和队列
	if err := declareTopology(); err != nil {
		return err
	}

	log.Println("RabbitMQ connected successfully")
	return nil
}

// declareTopology 声明交换机和队列
func declareTopology() error {
	// 声明交换机
	exchanges := []struct {
		name       string
		kind       string
		durable    bool
		autoDelete bool
	}{
		{ExchangePoints, "topic", true, false},
		{ExchangeCoupon, "direct", true, false},
		{ExchangeNotify, "direct", true, false},
		{ExchangeRecharge, "direct", true, false},
	}

	for _, ex := range exchanges {
		if err := ch.ExchangeDeclare(
			ex.name,
			ex.kind,
			ex.durable,
			ex.autoDelete,
			false,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("declare exchange %s failed: %w", ex.name, err)
		}
	}

	// 声明队列
	queues := []struct {
		name       string
		durable    bool
		autoDelete bool
	}{
		{QueuePointsSettle, true, false},
		{QueuePointsExpire, true, false},
		{QueueCouponPush, true, false},
		{QueueNotifyWechat, true, false},
		{QueueNotifySMS, true, false},
		{QueueRechargeCallback, true, false},
	}

	for _, q := range queues {
		if _, err := ch.QueueDeclare(
			q.name,
			q.durable,
			q.autoDelete,
			false,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("declare queue %s failed: %w", q.name, err)
		}
	}

	// 绑定队列到交换机
	bindings := []struct {
		queue      string
		exchange   string
		routingKey string
	}{
		{QueuePointsSettle, ExchangePoints, RoutingPointsEarned},
		{QueuePointsExpire, ExchangePoints, RoutingPointsExpired},
		{QueueCouponPush, ExchangeCoupon, RoutingCouponPush},
		{QueueNotifyWechat, ExchangeNotify, RoutingNotifyWechat},
		{QueueNotifySMS, ExchangeNotify, RoutingNotifySMS},
		{QueueRechargeCallback, ExchangeRecharge, RoutingRechargeCallback},
	}

	for _, b := range bindings {
		if err := ch.QueueBind(
			b.queue,
			b.routingKey,
			b.exchange,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("bind queue %s to exchange %s failed: %w", b.queue, b.exchange, err)
		}
	}

	return nil
}

// GetChannel 获取Channel
func GetChannel() *amqp.Channel {
	return ch
}

// Close 关闭连接
func Close() error {
	if ch != nil {
		if err := ch.Close(); err != nil {
			log.Printf("close channel failed: %v", err)
		}
	}
	if conn != nil {
		if err := conn.Close(); err != nil {
			log.Printf("close connection failed: %v", err)
		}
	}
	return nil
}

// Publish 发布消息
func Publish(exchange, routingKey string, body []byte) error {
	return ch.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
}

// Consume 消费消息
func Consume(queue, consumer string) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		queue,
		consumer,
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
}
