package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func init() {
	logCfg := zap.NewDevelopmentEncoderConfig()
	logCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	Logger = zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(logCfg),
			os.Stdout,
			zapcore.DebugLevel,
		),
		zap.AddCaller(),
	).Sugar()
}
