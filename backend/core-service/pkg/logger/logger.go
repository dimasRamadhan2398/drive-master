package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.SugaredLogger
	config       *Config
)

// Config holds logger configuration
type Config struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	OutputPath string `mapstructure:"output_path"`
}

// LogField represents a key-value pair for logging
type LogField struct {
	Key   string
	Value interface{}
}

// NewLogField creates a new LogField
func NewLogField(key string, value interface{}) LogField {
	return LogField{Key: key, Value: value}
}


// Init initializes the logger with the given config
func Init(cfg *Config) error {
	config = cfg

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

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var output zapcore.WriteSyncer
	if cfg.OutputPath != "" {
		file, err := os.OpenFile(cfg.OutputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		output = zapcore.AddSync(file)
	} else {
		output = zapcore.AddSync(os.Stdout)
	}

	core := zapcore.NewCore(encoder, output, level)
	logger := zap.New(core)
	globalLogger = logger.Sugar()

	return nil
}

// Sync flushes any buffered log entries
func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}

// GetLogger returns the global logger instance
func GetLogger() *zap.SugaredLogger {
	if globalLogger == nil {
		// Initialize with defaults if not initialized
		Init(&Config{Level: "info", Format: "console"})
	}
	return globalLogger
}

// Info logs an info message
func Info(msg string, fields ...LogField) {
	GetLogger().Infow(msg, fieldsToMap(fields)...)
}

// Debug logs a debug message
func Debug(msg string, fields ...LogField) {
	GetLogger().Debugw(msg, fieldsToMap(fields)...)
}

// Warn logs a warning message
func Warn(msg string, fields ...LogField) {
	GetLogger().Warnw(msg, fieldsToMap(fields)...)
}

// Error logs an error message
func Error(msg string, fields ...LogField) {
	GetLogger().Errorw(msg, fieldsToMap(fields)...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...LogField) {
	GetLogger().Fatalw(msg, fieldsToMap(fields)...)
}

func fieldsToMap(fields []LogField) []interface{} {
	result := make([]interface{}, len(fields)*2)
	for i, f := range fields {
		result[i*2] = f.Key
		result[i*2+1] = f.Value
	}
	return result
}

// WithField returns a logger with an added field
func WithField(key string, value interface{}) *zap.SugaredLogger {
	return GetLogger().With(key, value)
}

// WithTime returns a logger with a time field
func WithTime(key string, t time.Time) *zap.SugaredLogger {
	return GetLogger().With(key, t.Format(time.RFC3339))
}