package logger

import (
	"encoding/json"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestNewZap(t *testing.T) {
	tests := []struct {
		name string
		want *zap.Logger
	}{
		{"ok", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogger(""); !reflect.DeepEqual(got, tt.want) {
				if got == nil {
					t.Errorf("not got log install, got: %v", got)
				}
			}
		})
	}
}

func ExampleLogs() {
	raw := json.RawMessage(`{"a":"b", "ar": [1,2]}`)
	Debug("DebugMsg", zap.Any("data", &raw))
	Info("InfoMsg", zap.Any("data", &raw))
	Warn("WarnMsg", zap.Any("data", &raw))
	Error("ErrorMsg", zap.Any("data", &raw))

	// output:
}

