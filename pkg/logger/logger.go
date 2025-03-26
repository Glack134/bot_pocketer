package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger интерфейс для работы с логгером
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	Sync() error
}

// concreteLogger - реализация Logger на основе zap
type concreteLogger struct {
	zap *zap.SugaredLogger
}

// New создает новый экземпляр логгера
func New(level string) Logger {
	config := zap.NewProductionConfig()
	setLogLevel(&config, level)

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	zapLogger, err := config.Build()
	if err != nil {
		// Fallback на стандартный логгер
		zap.L().Error("Failed to create zap logger", zap.Error(err))
		zapLogger = zap.NewNop()
	}

	return &concreteLogger{
		zap: zapLogger.Sugar(),
	}
}

func setLogLevel(config *zap.Config, level string) {
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
}

// Реализация методов интерфейса
func (l *concreteLogger) Debug(args ...interface{}) {
	l.zap.Debug(args...)
}

func (l *concreteLogger) Info(args ...interface{}) {
	l.zap.Info(args...)
}

func (l *concreteLogger) Warn(args ...interface{}) {
	l.zap.Warn(args...)
}

func (l *concreteLogger) Error(args ...interface{}) {
	l.zap.Error(args...)
}

func (l *concreteLogger) Fatal(args ...interface{}) {
	l.zap.Fatal(args...)
}

func (l *concreteLogger) Debugf(format string, args ...interface{}) {
	l.zap.Debugf(format, args...)
}

func (l *concreteLogger) Infof(format string, args ...interface{}) {
	l.zap.Infof(format, args...)
}

func (l *concreteLogger) Warnf(format string, args ...interface{}) {
	l.zap.Warnf(format, args...)
}

func (l *concreteLogger) Errorf(format string, args ...interface{}) {
	l.zap.Errorf(format, args...)
}

func (l *concreteLogger) Fatalf(format string, args ...interface{}) {
	l.zap.Fatalf(format, args...)
}

func (l *concreteLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.zap.Debugw(msg, keysAndValues...)
}

func (l *concreteLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.zap.Infow(msg, keysAndValues...)
}

func (l *concreteLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.zap.Warnw(msg, keysAndValues...)
}

func (l *concreteLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.zap.Errorw(msg, keysAndValues...)
}

func (l *concreteLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.zap.Fatalw(msg, keysAndValues...)
}

func (l *concreteLogger) Sync() error {
	return l.zap.Sync()
}
