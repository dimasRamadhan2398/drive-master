package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server        ServerConfig   `yaml:"server"`
	Database      DatabaseConfig `yaml:"database"`
	Redis         RedisConfig    `yaml:"redis"`
	JWT           JWTConfig      `yaml:"jwt"`
	Log           LogConfig      `yaml:"log"`
	Kafka         KafkaConfig    `yaml:"kafka"`
	Midtrans      MidtransConfig `yaml:"midtrans"`
	Doku          DokuConfig     `yaml:"doku"`
	App           AppConfig      `yaml:"app"`
}

type ServerConfig struct {
	Port        int    `mapstructure:"port" yaml:"port"`
	Mode        string `mapstructure:"mode" yaml:"mode"`
	ReadTimeout int    `mapstructure:"read_timeout" yaml:"read_timeout"`
}

type DatabaseConfig struct {
	Host               string `mapstructure:"host" yaml:"host"`
	Port               int    `mapstructure:"port" yaml:"port"`
	Username           string `mapstructure:"username" yaml:"username"`
	Password           string `mapstructure:"password" yaml:"password"`
	Name               string `mapstructure:"name" yaml:"name"`
	SSLMode            string `mapstructure:"sslmode" yaml:"sslmode"`
	MaxOpenConnections int    `mapstructure:"max_open_connections" yaml:"maxOpenConnections"`
	MaxLifeTimeConnection int `mapstructure:"max_life_time_connection" yaml:"maxLifeTimeConnection"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections" yaml:"maxIdleConnections"`
	MaxIdleTime        int    `mapstructure:"max_idle_time" yaml:"maxIdleTime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Password string `mapstructure:"password" yaml:"password"`
	DB       int    `mapstructure:"db" yaml:"db"`
}

type JWTConfig struct {
	Secret                 string `mapstructure:"secret" yaml:"secret"`
	ExpiryHour             int    `mapstructure:"expiry_hour" yaml:"expiry_hour"`
	RefreshTokenExpiryDays  int    `mapstructure:"refresh_token_expiry_days" yaml:"refresh_token_expiry_days"`
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
	AppName        string `mapstructure:"app_name" yaml:"app_name"`
	AppEnv         string `mapstructure:"app_env" yaml:"app_env"`
	SignatureKey   string `mapstructure:"signature_key" yaml:"signature_key"`
	RateLimiterMax int    `mapstructure:"rate_limiter_max" yaml:"rate_limiter_max"`
	RateLimiterTime int   `mapstructure:"rate_limiter_time" yaml:"rate_limiter_time"`
	PaymentGateway string `mapstructure:"payment_gateway" yaml:"payment_gateway"`
}


type MidtransConfig struct {
	GopayCallbackURL string `mapstructure:"gopay_callback_url" yaml:"gopay_callback_url"`
	FrontendURL string `mapstructure:"frontend_url" yaml:"frontend_url"`
	ServerKey   string `mapstructure:"server_key" yaml:"server_key"`
	ClientKey   string `mapstructure:"client_key" yaml:"client_key"`
	Environment string `mapstructure:"environment" yaml:"environment"`
	MerchantID      string `mapstructure:"merchant_id"`
	BaseURL         string `mapstructure:"base_url"`
	SnapURL         string `mapstructure:"snap_url"`
	Enabled         bool   `mapstructure:"enabled"`
	TimeoutSeconds  int    `mapstructure:"timeout_seconds"`
	NotificationURL string `mapstructure:"notification_url"`
}

type DokuConfig struct {
	ClientID        string `mapstructure:"client_id" yaml:"client_id"`
	SecretKey       string `mapstructure:"secret_key" yaml:"secret_key"`
	Environment     string `mapstructure:"environment" yaml:"environment"`
	BaseURL         string `mapstructure:"base_url" yaml:"base_url"`
	NotificationURL string `mapstructure:"notification_url" yaml:"notification_url"`
	PaymentDueDate  int    `mapstructure:"payment_due_date" yaml:"payment_due_date"`
	FrontendURL     string `mapstructure:"frontend_url" yaml:"frontend_url"`
}

var AppCfg *Config

func setDefaults() {
	viper.SetDefault("server.port", 8004)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 60)

	// Database
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.username", "admin_drive")
	viper.SetDefault("database.password", "drivemaster123")
	viper.SetDefault("database.name", "drivemaster_payment_service")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.maxOpenConnections", 10)
	viper.SetDefault("database.maxLifeTimeConnection", 10)
	viper.SetDefault("database.maxIdleConnections", 10)
	viper.SetDefault("database.maxIdleTime", 10)

	// Redis
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

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
	viper.SetDefault("kafka.topic", "payment.events")
	viper.SetDefault("kafka.service_name", "payment-service")

	// App
	viper.SetDefault("app.app_name", "payment_service")
	viper.SetDefault("app.app_env", "development")
	viper.SetDefault("app.signature_key", "")
	viper.SetDefault("app.rate_limiter_max", 100)
	viper.SetDefault("app.rate_limiter_time", 1)
	viper.SetDefault("app.payment_gateway", "doku")

	// Midtrans
	viper.SetDefault("midtrans.frontend_url", "http://localhost:3000")
	viper.SetDefault("midtrans.server_key", "")
	viper.SetDefault("midtrans.client_key", "")
	viper.SetDefault("midtrans.environment", "sandbox")
	viper.SetDefault("midtrans.merchant_id", "")
	viper.SetDefault("midtrans.base_url", "https://api.midtrans.com")
	viper.SetDefault("midtrans.snap_url", "https://app.midtrans.com/snap")
	viper.SetDefault("midtrans.enabled", true)
	viper.SetDefault("midtrans.timeout_seconds", 30)
	viper.SetDefault("midtrans.notification_url", "")

	// Doku
	viper.SetDefault("doku.client_id", "")
	viper.SetDefault("doku.secret_key", "")
	viper.SetDefault("doku.environment", "sandbox")
	viper.SetDefault("doku.base_url", "https://api-sandbox.doku.com")
	viper.SetDefault("doku.notification_url", "")
	viper.SetDefault("doku.payment_due_date", 60)
	viper.SetDefault("doku.frontend_url", "http://localhost:3001")
}

func Load(path string) (*Config, error) {
	_ = godotenv.Load()
	setDefaults()

	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Server env overrides
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
	_ = viper.BindEnv("kafka.enabled", "KAFKA_ENABLED")

	// JWT env overrides
	_ = viper.BindEnv("jwt.secret", "JWT_SECRET")

	// App env overrides
	_ = viper.BindEnv("app.app_env", "APP_ENV")
	_ = viper.BindEnv("app.payment_gateway", "PAYMENT_GATEWAY")

	_ = viper.BindEnv("midtrans.server_key", "MIDTRANS_SERVER_KEY")
	_ = viper.BindEnv("midtrans.client_key", "MIDTRANS_CLIENT_KEY")
	_ = viper.BindEnv("midtrans.environment", "MIDTRANS_ENVIRONMENT")
	_ = viper.BindEnv("midtrans.merchant_id", "MIDTRANS_MERCHANT_ID")
	_ = viper.BindEnv("midtrans.base_url", "MIDTRANS_BASE_URL")
	_ = viper.BindEnv("midtrans.snap_url", "MIDTRANS_SNAP_URL")
	_ = viper.BindEnv("midtrans.enabled", "MIDTRANS_ENABLED")
	_ = viper.BindEnv("midtrans.notification_url", "MIDTRANS_NOTIFICATION_URL")
	_ = viper.BindEnv("midtrans.frontend_url", "MIDTRANS_FRONTEND_URL")

	_ = viper.BindEnv("doku.client_id", "DOKU_CLIENT_ID")
	_ = viper.BindEnv("doku.secret_key", "DOKU_SECRET_KEY")
	_ = viper.BindEnv("doku.environment", "DOKU_ENVIRONMENT")
	_ = viper.BindEnv("doku.base_url", "DOKU_BASE_URL")
	_ = viper.BindEnv("doku.notification_url", "DOKU_NOTIFICATION_URL")
	_ = viper.BindEnv("doku.payment_due_date", "DOKU_PAYMENT_DUE_DATE")
	_ = viper.BindEnv("doku.frontend_url", "DOKU_FRONTEND_URL")

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