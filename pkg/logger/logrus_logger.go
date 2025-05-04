package logger

// import (
// 	"os"
// 	"runtime"
// 	"strings"
// 	"time"

// 	"github.com/sirupsen/logrus"
// )

// type Logger struct {
// 	*logrus.Logger
// 	serviceName string
// }

// func NewLogger(serviceName string) *Logger {
// 	// Create a new logger instance
// 	logger := logrus.New()

// 	logger.SetOutput(os.Stdout)
// 	logger.SetFormatter(&CustomFormatter{
// 		ServiceName: serviceName,
// 		Formatter: &logrus.TextFormatter{
// 			ForceColors:     true,
// 			FullTimestamp:   true,
// 			TimestampFormat: "2006-01-02 15:04:05",
// 		},
// 	})

// 	// Set log level based on environment
// 	if os.Getenv("ENV") == "production" {
// 		logger.SetFormatter(&logrus.JSONFormatter{
// 			TimestampFormat: time.RFC3339Nano,
// 		})
// 		logger.SetLevel(logrus.InfoLevel)
// 	} else {
// 		logger.SetLevel(logrus.DebugLevel)
// 	}
// 	return &Logger{
// 		Logger:      logger,
// 		serviceName: serviceName,
// 	}
// }

// // CustomFormatter adds service name and colors to log output
// type CustomFormatter struct {
// 	ServiceName string
// 	logrus.Formatter
// }

// // Format implements logrus.Formatter interface
// func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
// 	// Add service name to all log entries
// 	if entry.Data == nil {
// 		entry.Data = make(logrus.Fields)
// 	}
// 	entry.Data["service"] = f.ServiceName

// 	// Add caller information for debug and error levels
// 	if entry.Level <= logrus.DebugLevel || entry.Level == logrus.ErrorLevel {
// 		if pc, file, line, ok := runtime.Caller(8); ok {
// 			funcName := runtime.FuncForPC(pc).Name()
// 			entry.Data["file"] = file[strings.LastIndex(file, "/")+1:]
// 			entry.Data["line"] = line
// 			entry.Data["func"] = funcName[strings.LastIndex(funcName, ".")+1:]
// 		}
// 	}

// 	// Colorize based on log level
// 	switch entry.Level {
// 	case logrus.DebugLevel:
// 		entry.Message = "\033[36m" + entry.Message + "\033[0m" // Cyan
// 	case logrus.InfoLevel:
// 		entry.Message = "\033[32m" + entry.Message + "\033[0m" // Green
// 	case logrus.WarnLevel:
// 		entry.Message = "\033[33m" + entry.Message + "\033[0m" // Yellow
// 	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
// 		entry.Message = "\033[31m" + entry.Message + "\033[0m" // Red
// 	}

// 	return f.Formatter.Format(entry)
// }

// // WithError adds error context to the logger
// func (l *Logger) WithError(err error) *logrus.Entry {
// 	return l.WithFields(logrus.Fields{
// 		"error": err.Error(),
// 	})
// }

// // WithComponent adds component name to the logger
// func (l *Logger) WithComponent(component string) *logrus.Entry {
// 	return l.WithFields(logrus.Fields{
// 		"component": component,
// 	})
// }
