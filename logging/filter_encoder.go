package logging

import (
	"errors"
	"sync"

	"go.uber.org/zap"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type FilterType string

const (
	RelegationFilter FilterType = "relegation"
	HideFilter       FilterType = "hide"
)

type FilterInfo struct {
	Type        FilterType      `json:"type" yaml:"type"`
	Fields      []zapcore.Field `json:"fields" yaml:"fields"`
	TargetLevel zapcore.Level   `json:"targetLevel,omitempty" yaml:"targetLevel,omitempty"`
}

var (
	FieldMap       = map[string]FilterInfo{}
	encoderSlice   = []string{}
	registerLocker = sync.Mutex{}
)

func RegisterFilterEncoder(name string, filter FilterInfo) error {
	registerLocker.Lock()
	defer registerLocker.Unlock()
	FieldMap[name] = filter
	encoderSlice = append([]string{name}, encoderSlice...)
	return zap.RegisterEncoder(name, NewFilterEncoder)
}

func NewFilterEncoder(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
	EncoderName := encoderSlice[len(encoderSlice)-1]
	encoderSlice = encoderSlice[0 : len(encoderSlice)-1]
	return newFilterEncoder(cfg, false, EncoderName), nil
}

func newFilterEncoder(cfg zapcore.EncoderConfig, spaced bool, encoderName string) *FilterEncoder {
	return &FilterEncoder{
		Name: encoderName,
		jsonEncoder: jsonEncoder{
			EncoderConfig: &cfg,
			buf:           Get(),
			spaced:        spaced,
		},
	}
}

var _filterPool = sync.Pool{New: func() interface{} {
	return &FilterEncoder{}
}}

func getFilterEncoder() *FilterEncoder {
	return _filterPool.Get().(*FilterEncoder)
}

type FilterEncoder struct {
	jsonEncoder
	Name string
}

func (enc *FilterEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes())
	return clone
}

func (enc *FilterEncoder) clone() *FilterEncoder {
	clone := getFilterEncoder()
	clone.Name = enc.Name
	clone.EncoderConfig = enc.EncoderConfig
	clone.spaced = enc.spaced
	clone.openNamespaces = enc.openNamespaces
	clone.buf = Get()
	return clone
}

func (enc *FilterEncoder) IsPassedFilter(fields []zapcore.Field, theFilter FilterInfo) bool {
	if len(theFilter.Fields) > 0 {
		for _, field := range fields {
			for _, filterField := range theFilter.Fields {
				if field.Key == filterField.Key && field.String == filterField.String && field.Type == filterField.Type &&
					field.Integer == filterField.Integer && field.Interface == filterField.Interface {
					return true
				}
			}
		}
	} else {
		return true
	}
	return false
}

func (enc *FilterEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	var theFilter FilterInfo
	var ok bool
	if theFilter, ok = FieldMap[enc.Name]; !ok {
		return nil, errors.New("not register filter fields")
	}
	switch theFilter.Type {
	case RelegationFilter:
		if enc.IsPassedFilter(fields, theFilter) {
			ent.Level = theFilter.TargetLevel
		}
	case HideFilter:
		if !enc.IsPassedFilter(fields, theFilter) {
			enc.buf.Reset()
			return enc.buf, nil
		}
	}
	return enc.jsonEncoder.EncodeEntry(ent, fields)
}
