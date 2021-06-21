package logging

import "go.uber.org/zap"

const (
	ControllerKey  = "katanomi/controller"
	NamespaceKey   = "katanomi/namespace"
	NameKey        = "katanomi/name"
	ResourceKey    = "katanomi/resource"
	SubresourceKey = "katanomi/subresource"
)

type BaseInfo struct {
	Controller  string
	Namespace   string
	Name        string
	Resource    string
	SubResource string
}

type structuredLogger struct {
	*zap.Logger
	Config Config
}

func (log *structuredLogger) InfoToFields(info BaseInfo) (fieldList []zap.Field) {
	fieldList = append(fieldList, zap.String(ControllerKey, info.Controller))
	fieldList = append(fieldList, zap.String(NamespaceKey, info.Namespace))
	fieldList = append(fieldList, zap.String(NameKey, info.Name))
	fieldList = append(fieldList, zap.String(ResourceKey, info.Resource))
	if info.SubResource != "" {
		fieldList = append(fieldList, zap.String(SubresourceKey, info.SubResource))
	}
	return fieldList
}

func (log *structuredLogger) Debug(msg string, info BaseInfo, fields ...zap.Field) {
	fieldList := log.InfoToFields(info)
	fieldList = append(fieldList, fields...)
	log.Logger.Debug(msg, fieldList...)
}

func (log *structuredLogger) Info(msg string, info BaseInfo, fields ...zap.Field) {
	fieldList := log.InfoToFields(info)
	fieldList = append(fieldList, fields...)
	log.Logger.Info(msg, fieldList...)
}

func (log *structuredLogger) Warn(msg string, info BaseInfo, fields ...zap.Field) {
	fieldList := log.InfoToFields(info)
	fieldList = append(fieldList, fields...)
	log.Logger.Warn(msg, fieldList...)
}

func (log *structuredLogger) Error(msg string, info BaseInfo, fields ...zap.Field) {
	fieldList := log.InfoToFields(info)
	fieldList = append(fieldList, fields...)
	log.Logger.Error(msg, fieldList...)
}

func (log *structuredLogger) DPanic(msg string, info BaseInfo, fields ...zap.Field) {
	fieldList := log.InfoToFields(info)
	fieldList = append(fieldList, fields...)
	log.Logger.DPanic(msg, fieldList...)
}

func (log *structuredLogger) Panic(msg string, info BaseInfo, fields ...zap.Field) {
	fieldList := log.InfoToFields(info)
	fieldList = append(fieldList, fields...)
	log.Logger.Panic(msg, fieldList...)
}

func (log *structuredLogger) Fatal(msg string, info BaseInfo, fields ...zap.Field) {
	fieldList := log.InfoToFields(info)
	fieldList = append(fieldList, fields...)
	log.Logger.Fatal(msg, fieldList...)
}

func NewLogger(config Config, opts ...zap.Option) (theLogger structuredLogger, err error) {
	var zapLogger *zap.Logger
	zapLogger, err = config.ToZapConfig().Build(opts...)
	if err != nil {
		return structuredLogger{}, err
	}
	theLogger = structuredLogger{
		Config: config,
		Logger: zapLogger,
	}
	return
}
