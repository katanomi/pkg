package logging

import (
	"go.uber.org/zap"
)

func init() {
	_ = RegisterFilterEncoder("hide", FilterInfo{
		Type:   HideFilter,
		Fields: []zap.Field{},
	})
	_ = RegisterFilterEncoder("relegation", FilterInfo{
		Type:   RelegationFilter,
		Fields: []zap.Field{},
	})
}
