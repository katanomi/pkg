package logging

import (
	"fmt"
	"testing"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var (
	infoA = BaseInfo{
		Controller: "a",
		Namespace:  "default",
		Name:       "byxiao",
		Resource:   "user",
	}

	infoB = BaseInfo{
		Controller: "b",
		Namespace:  "default",
		Name:       "byxiao",
		Resource:   "user",
	}
)

func TestNewLogger(t *testing.T) {
	err := RegisterFilterEncoder("filter", FilterInfo{
		Type:   HideFilter,
		Fields: []zap.Field{zap.String("katanomi/controller", "a")},
	})
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	encoderConfig := EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "severity",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    "\n",
	}
	theConfig := Config{
		Level:            zapcore.DebugLevel,
		Development:      true,
		Encoding:         "filter",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}
	loggerObj, err := NewLogger(theConfig)
	if err != nil {
		t.Fail()
		fmt.Println(err)
		return
	}
	loggerObj.Info("test", infoA)
	loggerObj.Info("test", infoB)
}

func TestNewLogger2(t *testing.T) {
	err := RegisterFilterEncoder("filter", FilterInfo{
		Type:        RelegationFilter,
		Fields:      []zap.Field{zap.String("katanomi/controller", "a")},
		TargetLevel: zapcore.DebugLevel,
	})
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	encoderConfig := EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "severity",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    "\n",
	}
	theConfig := Config{
		Level:            zapcore.DebugLevel,
		Development:      true,
		Encoding:         "filter",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}
	loggerObj, err := NewLogger(theConfig)
	if err != nil {
		t.Fail()
		fmt.Println(err)
		return
	}
	loggerObj.Info("test", infoA)
	loggerObj.Info("test", infoB)
}
