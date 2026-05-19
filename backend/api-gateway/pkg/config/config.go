package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig      `mapstructure:"server" yaml:"server"`
	JWT         JWTConfig         `mapstructure:"jwt" yaml:"jwt"`
	Services    ServicesConfig    `mapstructure:"services" yaml:"services"`
	RateLimiter RateLimiterConfig `mapstructure:"rate_limiter" yaml:"rate_limiter"`
	CORS        CORSConfig        `mapstructure:"cors" yaml:"cors"`
}

type ServerConfig struct {
	Host         string `mapstructure:"host" yaml:"host"`
	Port         int    `mapstructure:"port" yaml:"port"`
	Mode         string `mapstructure:"mode" yaml:"mode"`
	ReadTimeout  int    `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout" yaml:"write_timeout"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret" yaml:"secret"`
}

type ServicesConfig struct {
	UserServiceURL    string `mapstructure:"user_service_url" yaml:"user_service_url"`
	CoreServiceURL    string `mapstructure:"core_service_url" yaml:"core_service_url"`
	BookingServiceURL string `mapstructure:"booking_service_url" yaml:"booking_service_url"`
}

type RateLimiterConfig struct {
	Max  int `mapstructure:"max" yaml:"max"`
	Time int `mapstructure:"time" yaml:"time"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins" yaml:"allowed_origins"`
}

var AppCfg *Config

func setDefaults() {
	// Server
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)

	// JWT
	viper.SetDefault("jwt.secret", "your_jwt_secret_here")

	// Services
	viper.SetDefault("services.user_service_url", "http://localhost:8001")
	viper.SetDefault("services.core_service_url", "http://localhost:8002")
	viper.SetDefault("services.booking_service_url", "http://localhost:8003")

	// Rate Limiter
	viper.SetDefault("rate_limiter.max", 100)
	viper.SetDefault("rate_limiter.time", 1)

	// CORS
	viper.SetDefault("cors.allowed_origins", []string{
		"http://localhost:3000",
	})
}

func Load(path string) (*Config, error) {
	// Set defaults first
	setDefaults()

	// Config file
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Server env overrides
	_ = viper.BindEnv("server.host", "SERVER_HOST")
	_ = viper.BindEnv("server.port", "SERVER_PORT")
	_ = viper.BindEnv("server.mode", "SERVER_MODE")

	// JWT env overrides
	_ = viper.BindEnv("jwt.secret", "JWT_SECRET")

	// Service URLs env overrides
	_ = viper.BindEnv("services.user_service_url", "USER_SERVICE_URL")
	_ = viper.BindEnv("services.core_service_url", "CORE_SERVICE_URL")
	_ = viper.BindEnv("services.booking_service_url", "BOOKING_SERVICE_URL")

	// Rate limiter env overrides
	_ = viper.BindEnv("rate_limiter.max", "RATE_LIMITER_MAX")
	_ = viper.BindEnv("rate_limiter.time", "RATE_LIMITER_TIME")

	// Read config file
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