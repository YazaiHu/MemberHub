package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var GlobalConfig *Config

// Config 全局配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	WeChat   WeChatConfig   `mapstructure:"wechat"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"` // debug, release
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"name"`
	Charset         string        `mapstructure:"charset"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	LogLevel        string        `mapstructure:"log_level"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// RabbitMQConfig RabbitMQ配置
type RabbitMQConfig struct {
	URL             string `mapstructure:"url"`
	PrefetchCount   int    `mapstructure:"prefetch_count"`
	ReconnectDelay  int    `mapstructure:"reconnect_delay"`
	PublishRetry    int    `mapstructure:"publish_retry"`
	ConsumeRetry    int    `mapstructure:"consume_retry"`
	HeartbeatSecond int    `mapstructure:"heartbeat_second"`
}

// WeChatConfig 微信配置
type WeChatConfig struct {
	AppID          string `mapstructure:"app_id"`
	AppSecret      string `mapstructure:"app_secret"`
	MchID          string `mapstructure:"mch_id"`          // 商户号
	APIKey         string `mapstructure:"api_key"`         // API密钥
	NotifyURL      string `mapstructure:"notify_url"`      // 支付回调URL
	TemplateID     string `mapstructure:"template_id"`     // 模板消息ID
	MessageRateSec int    `mapstructure:"message_rate_sec"` // 消息发送速率（条/秒）
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret          string        `mapstructure:"secret"`
	ExpireTime      time.Duration `mapstructure:"expire_time"`
	AdminExpireTime time.Duration `mapstructure:"admin_expire_time"`
	Issuer          string        `mapstructure:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`       // debug, info, warn, error
	Format     string `mapstructure:"format"`      // json, console
	OutputPath string `mapstructure:"output_path"` // 日志文件路径
	MaxSize    int    `mapstructure:"max_size"`    // 日志文件最大大小（MB）
	MaxBackups int    `mapstructure:"max_backups"` // 保留旧文件的最大个数
	MaxAge     int    `mapstructure:"max_age"`     // 保留旧文件的最大天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	// 解析配置
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	GlobalConfig = &config
	return &config, nil
}

// setDefaults 设置默认值
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "release")
	viper.SetDefault("server.read_timeout", 30*time.Second)
	viper.SetDefault("server.write_timeout", 30*time.Second)

	// Database defaults
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", 1*time.Hour)
	viper.SetDefault("database.log_level", "warn")

	// Redis defaults
	viper.SetDefault("redis.pool_size", 50)
	viper.SetDefault("redis.min_idle_conns", 10)
	viper.SetDefault("redis.dial_timeout", 5*time.Second)
	viper.SetDefault("redis.read_timeout", 3*time.Second)
	viper.SetDefault("redis.write_timeout", 3*time.Second)

	// RabbitMQ defaults
	viper.SetDefault("rabbitmq.prefetch_count", 10)
	viper.SetDefault("rabbitmq.reconnect_delay", 5)
	viper.SetDefault("rabbitmq.publish_retry", 3)
	viper.SetDefault("rabbitmq.consume_retry", 3)
	viper.SetDefault("rabbitmq.heartbeat_second", 60)

	// WeChat defaults
	viper.SetDefault("wechat.message_rate_sec", 20)

	// JWT defaults
	viper.SetDefault("jwt.expire_time", 7*24*time.Hour)      // 用户token 7天
	viper.SetDefault("jwt.admin_expire_time", 24*time.Hour)  // 管理员token 1天
	viper.SetDefault("jwt.issuer", "membership-system")

	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output_path", "logs/app.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 30)
	viper.SetDefault("log.max_age", 30)
	viper.SetDefault("log.compress", true)
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return GlobalConfig
}
