/*
Copyright 2022 The Katanomi Authors.

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

	"github.com/katanomi/pkg/command/io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKey struct{}

// WithLogger set a logger instance into a context
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, loggerKey{}, logger)
}

// GetLogger get a logger instance form a context
func GetLogger(ctx context.Context) (logger *zap.SugaredLogger) {
	if ctx == nil {
		return nil
	}
	val := ctx.Value(loggerKey{})
	if val == nil {
		return nil
	}
	if l, ok := val.(*zap.SugaredLogger); ok {
		return l
	}
	return nil
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
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), writer, level)
	return zap.New(core, opts...).Sugar()
}
