package logging

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
)

func TestNewJSONEncoder(t *testing.T) {
	err := zap.RegisterEncoder("filter", NewJSONEncoder)
	if err != nil {
		t.Fail()
		fmt.Println(err)
		return
	}
}
