/*
Copyright 2022 The AlaudaDevops Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logger

import (
	"context"

	"github.com/AlaudaDevops/pkg/command/io"
	"knative.dev/pkg/logging"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// WithLogger set a logger instance into a context
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	// reusing knative method
	return logging.WithLogger(ctx, logger)
}

// GetLogger get a logger instance form a context
func GetLogger(ctx context.Context) (logger *zap.SugaredLogger) {
	if ctx == nil {
		return nil
	}
	// reusing method from knative
	return logging.FromContext(ctx)
}

// NewLoggerFromContext similar to `GetLogger`, but return a default logger if there is no
// logger instance in the context
func NewLoggerFromContext(ctx context.Context) (logger *zap.SugaredLogger) {
	if logger = GetLogger(ctx); logger == nil {
		streams := io.MustGetIOStreams(ctx)
		logger = NewLogger(zapcore.AddSync(streams.ErrOut), zapcore.DebugLevel)
	}
	return
}

// NewLogger construct a logger
func NewLogger(writer zapcore.WriteSyncer, level zapcore.LevelEnabler, opts ...zap.Option) *zap.SugaredLogger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey: "msg",
		LevelKey:   "level",
		NameKey:    "logger",
		//// by commenting out this attribute
		//// will remove the level attribute the logger output result
		//// which is a better experience for CLI commands
		//// TODO: may revisit this and improve the overall CLI experience again
		// EncodeLevel:    zapcore.LowercaseLevelEncoder,
		// EncodeLevel:    EmojiLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), writer, level)
	return zap.New(core, opts...).Sugar()
}

// EmojiLevelEncoder prints an emoji instead of the log level
// ❌ for Panic, Error and Fatal levels
// 🐛 for Debug
// ❗️ for Warning
// 📢 for Info and everything else
func EmojiLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	value := "==> "
	switch l {
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.ErrorLevel, zapcore.FatalLevel:
		value += "❌"
	case zap.DebugLevel:
		value += "🐛"
	case zap.WarnLevel:
		value += "❗️"
	case zap.InfoLevel:
		fallthrough
	default:
		value += "📢"
	}
	enc.AppendString(value)
}
