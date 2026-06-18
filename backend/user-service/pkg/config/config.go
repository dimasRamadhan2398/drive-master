package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Database    DatabaseConfig    `yaml:"database"`
	Email       EmailConfig       `yaml:"email"`
	Redis       RedisConfig       `yaml:"redis"`
	ImageKit    ImageKitConfig    `yaml:"imagekit"`
	JWT         JWTConfig         `yaml:"jwt"`
	Log         LogConfig         `yaml:"log"`
	Kafka       KafkaConfig       `yaml:"kafka"`
	App         AppConfig         `yaml:"app"`
	CoreService CoreServiceConfig `yaml:"core_service"`
}

type ServerConfig struct {
	Port        int    `mapstructure:"port" yaml:"port"`
	Mode        string `mapstructure:"mode" yaml:"mode"`
	ReadTimeout int    `mapstructure:"read_timeout" yaml:"read_timeout"`
}

type DatabaseConfig struct {
	Host                  string `mapstructure:"host" yaml:"host"`
	Port                  int    `mapstructure:"port" yaml:"port"`
	Username              string `mapstructure:"username" yaml:"username"`
	Password              string `mapstructure:"password" yaml:"password"`
	Name                  string `mapstructure:"name" yaml:"name"`
	SSLMode               string `mapstructure:"sslmode" yaml:"sslmode"`
	MaxOpenConnections    int    `mapstructure:"max_open_connections" yaml:"maxOpenConnections"`
	MaxLifeTimeConnection int    `mapstructure:"max_life_time_connection" yaml:"maxLifeTimeConnection"`
	MaxIdleConnections    int    `mapstructure:"max_idle_connections" yaml:"maxIdleConnections"`
	MaxIdleTime           int    `mapstructure:"max_idle_time" yaml:"maxIdleTime"`
}

type EmailConfig struct {
	Token     string `mapstructure:"token" yaml:"token"`
	APIKey    string `mapstructure:"api_key" yaml:"api_key"`
	FromEmail string `mapstructure:"from_email" yaml:"from_email"`
	FromName  string `mapstructure:"from_name" yaml:"from_name"`
	AppName   string `mapstructure:"app_name" yaml:"app_name"`
	Enabled   bool   `mapstructure:"enabled" yaml:"enabled"`
	Host      string `mapstructure:"smtp_host" yaml:"smtp_host"`
	Port      int    `mapstructure:"smtp_port" yaml:"smtp_port"`
	User      string `mapstructure:"smtp_username" yaml:"smtp_username"`
	Password  string `mapstructure:"smtp_password" yaml:"smtp_password"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Password string `mapstructure:"password" yaml:"password"`
	DB       int    `mapstructure:"db" yaml:"db"`
}

type ImageKitConfig struct {
	ID          string `mapstructure:"id" yaml:"id"`
	PrivateKey  string `mapstructure:"private_key" yaml:"private_key"`
	URLEndpoint string `mapstructure:"url_endpoint" yaml:"url_endpoint"`
	PublicKey   string `mapstructure:"public_key" yaml:"public_key"`
}

type JWTConfig struct {
	Secret                 string `mapstructure:"secret" yaml:"secret"`
	ExpiryHour             int    `mapstructure:"expiry_hour" yaml:"expiry_hour"`
	RefreshTokenExpiryDays int    `mapstructure:"refresh_token_expiry_days" yaml:"refresh_token_expiry_days"`
}

type LogConfig struct {
	Level  string `mapstructure:"level" yaml:"level"`
	Format string `mapstructure:"format" yaml:"format"`
}

type KafkaConfig struct {
	Enabled     bool     `mapstructure:"enabled" yaml:"enabled"`
	Brokers     []string `mapstructure:"brokers" yaml:"brokers"`
	Topic       string   `mapstructure:"topic" yaml:"topic"`
	ServiceName string   `mapstructure:"service_name" yaml:"service_name"`
}

type AppConfig struct {
	AppName         string `mapstructure:"app_name" yaml:"app_name"`
	AppEnv          string `mapstructure:"app_env" yaml:"app_env"`
	SignatureKey    string `mapstructure:"signature_key" yaml:"signature_key"`
	RateLimiterMax  int    `mapstructure:"rate_limiter_max" yaml:"rate_limiter_max"`
	RateLimiterTime int    `mapstructure:"rate_limiter_time" yaml:"rate_limiter_time"`
}

type CoreServiceConfig struct {
	BaseURL string `mapstructure:"base_url" yaml:"base_url"`
}

var AppCfg *Config

func setDefaults() {
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
	viper.SetDefault("jwt.refresh_token_expiry_days", 7)

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

	// Core Service
	viper.SetDefault("core_service.base_url", "http://localhost:8002")

	viper.SetDefault("email.smtp_host", "sandbox.smtp.mailtrap.io")
	viper.SetDefault("email.smtp_port", 587)
	viper.SetDefault("email.smtp_user", "49bbd2a554bd3e")
	viper.SetDefault("email.smtp_password", "4f4bc1b03a1d70")
	viper.SetDefault("email.from_name", "Drive Master Indonesia")
}

func Load(path string) (*Config, error) {
	// IMPORTANT: Set defaults FIRST before reading config file
	setDefaults()

	// Set config file path
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// Bind environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Allow env var overrides for server config
	_ = viper.BindEnv("server.port", "PORT")
	_ = viper.BindEnv("server.mode", "SERVER_MODE")

	// Database env overrides
	_ = viper.BindEnv("database.host", "POSTGRES_HOST")
	_ = viper.BindEnv("database.port", "POSTGRES_PORT")
	_ = viper.BindEnv("database.username", "POSTGRES_USER")
	_ = viper.BindEnv("database.password", "POSTGRES_PASSWORD")
	_ = viper.BindEnv("database.name", "POSTGRES_DB")

	// Redis env overrides
	_ = viper.BindEnv("redis.host", "REDIS_HOST")
	_ = viper.BindEnv("redis.port", "REDIS_PORT")
	_ = viper.BindEnv("redis.password", "REDIS_PASSWORD")

	// Kafka env overrides
	_ = viper.BindEnv("kafka.brokers", "KAFKA_BROKERS")

	// ImageKit env overrides
	_ = viper.BindEnv("imagekit.id", "IMAGEKIT_ID")
	_ = viper.BindEnv("imagekit.private_key", "IMAGEKIT_PRIVATE_KEY")
	_ = viper.BindEnv("imagekit.url_endpoint", "IMAGEKIT_URL_ENDPOINT")

	// JWT env overrides
	_ = viper.BindEnv("jwt.secret", "JWT_SECRET")
	_ = viper.BindEnv("jwt.expiry_hour", "JWT_EXPIRY_HOUR")
	_ = viper.BindEnv("jwt.refresh_token_expiry_days", "JWT_REFRESH_TOKEN_EXPIRY_DAYS")

	// Email env overrides
	_ = viper.BindEnv("email.token", "EMAIL_TOKEN")
	_ = viper.BindEnv("email.api_key", "EMAIL_API_KEY")
	_ = viper.BindEnv("email.from_email", "EMAIL_FROM_EMAIL")
	_ = viper.BindEnv("email.from_name", "EMAIL_FROM_NAME")
	_ = viper.BindEnv("email.enabled", "EMAIL_ENABLED")

	// SMTP env overrides
	_ = viper.BindEnv("email.smtp_host", "EMAIL_SMTP_HOST")
	_ = viper.BindEnv("email.smtp_port", "EMAIL_SMTP_PORT")
	_ = viper.BindEnv("email.smtp_username", "EMAIL_SMTP_USERNAME")
	_ = viper.BindEnv("email.smtp_password", "EMAIL_SMTP_PASSWORD")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

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
