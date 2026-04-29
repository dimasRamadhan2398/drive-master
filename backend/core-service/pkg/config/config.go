package config

import (
	"os"

	"gopkg.in/yaml.v3"
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

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	AppCfg = &cfg
	return &cfg, nil
}

func Get() *Config {
	return AppCfg
}
