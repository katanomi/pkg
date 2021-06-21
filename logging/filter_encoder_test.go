package logging

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestRegisterFilterEncoder(t *testing.T) {
	err := RegisterFilterEncoder("filter", FilterInfo{
		Type:   HideFilter,
		Fields: []zap.Field{zap.Int("a", 1)},
	})
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     SimpleTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	// 构造 Config
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      true,
		Encoding:         "filter",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}
	logger, err := config.Build()
	if err != nil {
		t.Fail()
		fmt.Println(err)
		return
	}
	logger.Info("ok", zap.Int("a", 1))
	logger.Info("ok", zap.Int("a", 2))
	logger.Info("ok", zap.Int("b", 1))
}
