package logging

import (
	"time"

	"go.uber.org/zap"

	"go.uber.org/zap/zapcore"
)

const (
	DefaultTimeFormat = "2006-01-02 15:04:05 MST"
)

func encodeTimeLayout(t time.Time, layout string, enc zapcore.PrimitiveArrayEncoder) {
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}

	enc.AppendString(t.Format(layout))
}

func SimpleTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	encodeTimeLayout(t, DefaultTimeFormat, enc)
}

type EncoderConfig struct {
	TimeKey       string `json:"timeKey" yaml:"timeKey"`
	LevelKey      string `json:"levelKey" yaml:"levelKey"`
	NameKey       string `json:"nameKey" yaml:"nameKey"`
	CallerKey     string `json:"callerKey" yaml:"callerKey"`
	MessageKey    string `json:"messageKey" yaml:"messageKey"`
	StacktraceKey string `json:"stacktraceKey" yaml:"stacktraceKey"`
	LineEnding    string `json:"lineEnding" yaml:"lineEnding"`
	TimeFormat    string `json:"timeFormat,omitempty" yaml:"timeFormat,omitempty"`
}

func (ec *EncoderConfig) ToZapEncoderConfig() zapcore.EncoderConfig {
	if ec.TimeFormat == "" {
		ec.TimeFormat = DefaultTimeFormat
	}
	theTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		encodeTimeLayout(t, ec.TimeFormat, enc)
	}
	return zapcore.EncoderConfig{
		TimeKey:        ec.TimeKey,
		LevelKey:       ec.LevelKey,
		NameKey:        ec.NameKey,
		CallerKey:      ec.CallerKey,
		MessageKey:     ec.MessageKey,
		StacktraceKey:  ec.StacktraceKey,
		LineEnding:     ec.LineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     theTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

type Config struct {
	// Level is the minimum enabled logging level. Note that this is a dynamic
	// level, so calling Config.Level.SetLevel will atomically change the log
	// level of all loggers descended from this config.
	Level zapcore.Level `json:"level" yaml:"level"`
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	Development bool `json:"development" yaml:"development"`
	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`
	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`
	// Encoding sets the logger's encoding. Valid values are "json" and
	// "console", as well as any third-party encodings registered via
	// RegisterEncoder.
	Encoding string `json:"encoding" yaml:"encoding"`
	// EncoderConfig sets options for the chosen encoder.
	EncoderConfig EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
	// OutputPaths is a list of URLs or file paths to write logging output to.
	// See Open for details.
	OutputPaths []string `json:"outputPaths" yaml:"outputPaths"`
	// ErrorOutputPaths is a list of URLs to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects internal errors; for sample code that
	// sends error-level logs to a different location from info- and debug-level
	// logs, see the package-level AdvancedConfiguration example.
	ErrorOutputPaths []string `json:"errorOutputPaths" yaml:"errorOutputPaths"`
}

func (c *Config) ToZapConfig() zap.Config {
	return zap.Config{
		Level:             zap.NewAtomicLevelAt(c.Level),
		Development:       c.Development,
		DisableCaller:     c.DisableCaller,
		DisableStacktrace: c.DisableStacktrace,
		Encoding:          c.Encoding,
		EncoderConfig:     c.EncoderConfig.ToZapEncoderConfig(),
		OutputPaths:       c.OutputPaths,
		ErrorOutputPaths:  c.ErrorOutputPaths,
	}
}
