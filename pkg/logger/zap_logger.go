package logger

import (
	"os"
	"runtime"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zapLogger   *zap.Logger
	serviceName string
}

var (
	instance *Logger
	once     sync.Once
)

// NewLogger initializes the logger only once (singleton)
func NewLogger(serviceName string) *Logger {
	once.Do(func() {
		instance = &Logger{
			zapLogger:   initZapLogger(serviceName),
			serviceName: serviceName,
		}
	})
	return instance
}

// initZapLogger configures a zap logger based on environment
func initZapLogger(serviceName string) *zap.Logger {
	var cfg zap.Config
	env := os.Getenv("ENV")

	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	cfg.OutputPaths = []string{"stdout"}
	cfg.InitialFields = map[string]interface{}{
		"service": serviceName,
	}

	logger, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
	return logger
}

// WithComponent adds a component field to the logger
func (l *Logger) WithComponent(component string) *zap.Logger {
	return l.zapLogger.With(zap.String("component", component))
}

// WithError logs error with additional error context

func (l *Logger) String(key string, val string) zap.Field {
	return zap.String(key, val)
}
func (l *Logger) WithError(err error) *zap.Logger {
	return l.zapLogger.With(zap.String("error", err.Error()))
}

// Debug logs a debug message with optional fields
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, fields...)
}

// Info logs an info message with optional fields
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

// Panic logs a panic message
func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.zapLogger.Panic(msg, fields...)
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zapLogger.Fatal(msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, fields...)
}

// ErrorField creates a zap field for an error
func (l *Logger) ErrorFields(err error) zap.Field {
	return zap.Error(err)
}

// Sync flushes buffered log entries
func (l *Logger) Sync() {
	_ = l.zapLogger.Sync()
}

// WithFields adds multiple fields to the logger (like logrus.WithFields)
func (l *Logger) WithFields(fields map[string]interface{}) *zap.Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return l.zapLogger.With(zapFields...)
}

// AddCallerInfo adds file, line, and function name manually (useful for deeper call stacks)
func AddCallerInfo(skip int) zap.Field {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return zap.Skip()
	}
	funcName := runtime.FuncForPC(pc).Name()
	return zap.String("caller", formatCaller(file, line, funcName))
}

func formatCaller(file string, line int, funcName string) string {
	fileName := file[strings.LastIndex(file, "/")+1:]
	funcShort := funcName[strings.LastIndex(funcName, ".")+1:]
	return fileName + ":" + funcShort + ":" + string(rune(line))
}
