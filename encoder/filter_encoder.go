package encoder

import (
	"errors"
	"strings"

	"go.uber.org/zap"

	"sync"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var (
	_theFilter     = map[string][]string{}
	encoderSlice   = []string{}
	registerLocker = sync.Mutex{}
)

func SetFilter(encoderName string, loggerList []string) (err error) {
	if _, ok := _theFilter[encoderName]; !ok {
		return errors.New("encoder not exist")
	}
	_theFilter[encoderName] = loggerList
	return err
}

func RegisterFilterEncoder(name string, filter []string) error {
	registerLocker.Lock()
	defer registerLocker.Unlock()
	_theFilter[name] = filter
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

func (enc *FilterEncoder) IsPassedFilter(ent zapcore.Entry, theFilter []string) bool {
	for _, v := range theFilter {
		if strings.Contains(ent.LoggerName, v) {
			return true
		}
	}
	return false
}

func (enc *FilterEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	var ok bool
	var filter []string
	if filter, ok = _theFilter[enc.Name]; !ok {
		return nil, errors.New("encoder not exist")
	}
	if len(filter) == 0 {
		return enc.jsonEncoder.EncodeEntry(ent, fields)
	}
	if !enc.IsPassedFilter(ent, filter) {
		enc.buf.Reset()
		return enc.buf, nil
	}
	return enc.jsonEncoder.EncodeEntry(ent, fields)
}
