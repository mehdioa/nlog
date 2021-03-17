// simple_test
package nlog

import (
	"testing"
)

func TestJsonSimple(t *testing.T) {
	j := NewJsonFormatter(false)
	l := NewLogger(DebugLevel, j)

	l.Debug("debug", Data{"key": 1, "key2": "string", "key3": false})
	l.Info("info", Data{"key": 1, "key2": "string", "key3": false})
	l.Warn("warn", Data{"key": 1, "key2": "string", "key3": false})
	l.Error("error", Data{"key": 1, "key2": "string", "key3": false})
}

func TestTextSimple(t *testing.T) {
	j := NewTextFormatter(true, true)
	l := NewLogger(DebugLevel, j)

	l.Debugf("debug")
	l.Debug("debug", Data{"key": 1, "key2": "string", "key3": false})
	l.Info("info", Data{"key": 1, "key2": "string", "key3": false})
	l.Warn("warn", Data{"key": 1, "key2": "string", "key3": false})
	l.Error("error", Data{"key": 1, "key2": "string", "key3": false})
}
