package config

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(env string) (*zap.Logger, error) {
	if env != "production" {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stderr"}

		return config.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// Production Logging - Separate files by level
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Ensure logs directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	// Open log files
	infoFile, err := os.OpenFile("logs/info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	warnFile, err := os.OpenFile("logs/warn.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	errorFile, err := os.OpenFile("logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Define level filters
	infoLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == zapcore.InfoLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == zapcore.WarnLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zapcore.ErrorLevel
	})

	// Create core with multi-sink
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoFile), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnFile), warnLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorFile), errorLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}
