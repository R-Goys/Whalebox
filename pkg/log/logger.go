package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func InitLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{"stdout", "record.log"}

	zapLogger, err := config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	logger = zapLogger.Sugar()
}

func Info(args ...interface{}) {
	logger.Info(args)
}

func Debug(args ...interface{}) {
	logger.Debug(args)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

func Error(args ...interface{}) {
	logger.Error(args)
}
