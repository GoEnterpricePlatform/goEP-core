package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type HttpLogger struct {
	log *zap.Logger
}

func NewHttpLogger(appEnv string) *HttpLogger {
	var zapLevel zapcore.Level
	var encoderCfg zapcore.EncoderConfig

	if appEnv == "prod" {
		encoderCfg = zap.NewProductionEncoderConfig()
		zapLevel = zapcore.InfoLevel
	} else {
		zapLevel = zapcore.DebugLevel
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 - 15:04:05")

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(ConsoleWriter{}),
		zapLevel,
	)

	return &HttpLogger{log: zap.New(core)}
}

type ConsoleWriter struct{}

func (cw ConsoleWriter) Write(p []byte) (int, error) {
	return fmt.Print(string(p))
}
