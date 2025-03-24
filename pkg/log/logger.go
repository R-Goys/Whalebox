package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:   "msg",
		EncodeLevel:  zapcore.CapitalColorLevelEncoder, // 彩色日志
		EncodeTime:   zapcore.ISO8601TimeEncoder,       // 时间格式
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	config.Encoding = "console" // 切换为 console 格式

	logger, _ = config.Build()
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
