package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	Email     EmailConfig     `yaml:"email"`
	Redis     RedisConfig     `yaml:"redis"`
	ImageKit  ImageKitConfig  `yaml:"imagekit"`
	JWT       JWTConfig       `yaml:"jwt"`
	Log       LogConfig       `yaml:"log"`
	Kafka     KafkaConfig     `yaml:"kafka"`
	App       AppConfig       `yaml:"app"`
}

type ServerConfig struct {
	Port        int `yaml:"port"`
	Mode        string `yaml:"mode"`
	ReadTimeout int `yaml:"read_timeout"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
	MaxOpenConnections    int    `json:"maxOpenConnections"`
	MaxLifeTimeConnection int    `json:"maxLifeTimeConnection"`
	MaxIdleConnections    int    `json:"maxIdleConnections"`
	MaxIdleTime           int    `json:"maxIdleTime"`
}

type EmailConfig struct {
	Token     string `yaml:"token"`
	APIKey    string `yaml:"api_key"`
	FromEmail string `yaml:"from_email"`
	AppName   string `yaml:"app_name"`
	Enabled   bool   `yaml:"enabled"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type ImageKitConfig struct {
	ID          string `yaml:"id"`
	PrivateKey  string `yaml:"private_key"`
	URLEndpoint string `yaml:"url_endpoint"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpiryHour  int    `yaml:"expiry_hour"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type KafkaConfig struct {
	Enabled     bool     `yaml:"enabled"`
	Brokers     []string `yaml:"brokers"`
	Topic       string   `yaml:"topic"`
	ServiceName string   `yaml:"service_name"`
}

type AppConfig struct {
	AppName        string `yaml:"app_name"`
	AppEnv        string `yaml:"app_env"`
	SignatureKey  string `yaml:"signature_key"`
	RateLimiterMax int    `yaml:"rate_limiter_max"`
	RateLimiterTime int   `yaml:"rate_limiter_time"`
}

var AppCfg *Config

func setDefaults(){
	viper.SetDefault("server.port", 8001)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 60)

	// Database
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.username", "drivemaster")
	viper.SetDefault("database.password", "drivemaster123")
	viper.SetDefault("database.name", "drivemaster")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.maxOpenConnections", 10)
	viper.SetDefault("database.maxLifeTimeConnection", 10)
	viper.SetDefault("database.maxIdleConnections", 10)
	viper.SetDefault("database.maxIdleTime", 10)

	// Email
	viper.SetDefault("email.token", "")
	viper.SetDefault("email.api_key", "")
	viper.SetDefault("email.from_email", "")
	viper.SetDefault("email.app_name", "user_service")
	viper.SetDefault("email.enabled", false)

	// Redis
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// ImageKit
	viper.SetDefault("imagekit.id", "")
	viper.SetDefault("imagekit.private_key", "")
	viper.SetDefault("imagekit.url_endpoint", "")

	// JWT
	viper.SetDefault("jwt.secret", "")
	viper.SetDefault("jwt.expiry_hour", 24)

	// Log
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	// Kafka
	viper.SetDefault("kafka.enabled", false)
	viper.SetDefault("kafka.brokers", []string{"localhost:29092"})
	viper.SetDefault("kafka.topic", "service-logs")
	viper.SetDefault("kafka.service_name", "user-service")

	// App
	viper.SetDefault("app.app_name", "user_service")
	viper.SetDefault("app.app_env", "development")
	viper.SetDefault("app.signature_key", "")
	viper.SetDefault("app.rate_limiter_max", 100)
	viper.SetDefault("app.rate_limiter_time", 1)
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	setDefaults()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Get() *Config {
	return AppCfg
}

func Set(cfg *Config) {
	AppCfg = cfg
}
