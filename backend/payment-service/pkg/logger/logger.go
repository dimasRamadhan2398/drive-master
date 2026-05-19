package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level  string `mapstructure:"level" yaml:"level"`
	Format string `mapstructure:"format" yaml:"format"`
}

var logger *zap.Logger

func Init(cfg *Config) error {
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	var config zap.Config
	if cfg.Format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.Level = zap.NewAtomicLevelAt(level)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	logger, err = config.Build()
	return err
}

func GetLogger() *zap.Logger {
	return logger
}

func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}

func Info(msg string, fields ...Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	logger.Error(msg, fields...)
}

func Debug(msg string, fields ...Field) {
	logger.Debug(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	logger.Fatal(msg, fields...)
}

type Field = zap.Field

func LogField(key string, value interface{}) Field {
	return zap.Any(key, value)
}

func NewLogField(key string, value interface{}) Field {
	return zap.Any(key, value)
}