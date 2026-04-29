package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the notification service
type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Database    DatabaseConfig    `yaml:"database"`
	Email       EmailConfig       `yaml:"email"`
	WhatsApp    WhatsAppConfig    `yaml:"whatsapp"`
	Redis       RedisConfig       `yaml:"redis"`
	Kafka       KafkaConfig       `yaml:"kafka"`
	Log         LogConfig         `yaml:"log"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port        int    `yaml:"port"`
	Mode        string `yaml:"mode"`
	ReadTimeout int    `yaml:"read_timeout"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

// EmailConfig holds email service configuration
type EmailConfig struct {
	Provider    string `yaml:"provider"` // mailtrap, sendgrid, smtp
	APIKey      string `yaml:"api_key"`
	FromEmail   string `yaml:"from_email"`
	FromName    string `yaml:"from_name"`
	Enabled     bool   `yaml:"enabled"`
}

// WhatsAppConfig holds WhatsApp service configuration
type WhatsAppConfig struct {
	APIURL       string `yaml:"api_url"`
	APIToken     string `yaml:"api_token"`
	PhoneNumber  string `yaml:"phone_number"`
	Enabled      bool   `yaml:"enabled"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Enabled     bool     `yaml:"enabled"`
	Brokers     []string `yaml:"brokers"`
	GroupID     string   `yaml:"group_id"`
	Topics      []string `yaml:"topics"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

var AppCfg *Config

// Load reads configuration from a YAML file
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

// Get returns the loaded configuration
func Get() *Config {
	return AppCfg
}