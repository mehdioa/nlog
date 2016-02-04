package nlog

import (
	//	"os"
	"testing"
)

type M map[string]interface{}

var testObject = M{
	"foo": "bar",
	"bah": M{
		"int":      1,
		"float":    -100.23,
		"date":     "06-01-01T15:04:05-0700",
		"bool":     true,
		"nullable": nil,
	},
}

func BenchmarkJsonSimple(b *testing.B) {
	j := NewJsonFormatter(false)
	l := NewLogger(DebugLevel, j)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("debug", Data{"key": 1, "key2": "string", "key3": false})
		l.Info("info", Data{"key": 1, "key2": "string", "key3": false})
		l.Warn("warn", Data{"key": 1, "key2": "string", "key3": false})
		l.Error("error", Data{"key": 1, "key2": "string", "key3": false})
	}
	b.StopTimer()
}

func BenchmarkJsonComplex(b *testing.B) {
	j := NewJsonFormatter(false)
	l := NewLogger(DebugLevel, j)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("debug", Data{"key": 1, "obj": testObject})
		l.Info("info", Data{"key": 1, "obj": testObject})
		l.Warn("warn", Data{"key": 1, "obj": testObject})
		l.Error("error", Data{"key": 1, "obj": testObject})
	}
	b.StopTimer()
}
func BenchmarkTextSimple(b *testing.B) {
	j := NewTextFormatter(true, true)
	l := NewLogger(DebugLevel, j)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("debug", Data{"key": 1, "key2": "string", "key3": false})
		l.Info("info", Data{"key": 1, "key2": "string", "key3": false})
		l.Warn("warn", Data{"key": 1, "key2": "string", "key3": false})
		l.Error("error", Data{"key": 1, "key2": "string", "key3": false})
	}
	b.StopTimer()
}

func BenchmarkTextComplex(b *testing.B) {
	j := NewTextFormatter(true, true)
	l := NewLogger(DebugLevel, j)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("debug", Data{"key": 1, "obj": testObject})
		l.Info("info", Data{"key": 1, "obj": testObject})
		l.Warn("warn", Data{"key": 1, "obj": testObject})
		l.Error("error", Data{"key": 1, "obj": testObject})
	}
	b.StopTimer()
}
func BenchmarkTextNoColorSimple(b *testing.B) {
	j := NewTextFormatter(false, false)
	l := NewLogger(DebugLevel, j)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("debug", Data{"key": 1, "key2": "string", "key3": false})
		l.Info("info", Data{"key": 1, "key2": "string", "key3": false})
		l.Warn("warn", Data{"key": 1, "key2": "string", "key3": false})
		l.Error("error", Data{"key": 1, "key2": "string", "key3": false})
	}
	b.StopTimer()
}

func BenchmarkTextNoColorComplex(b *testing.B) {
	j := NewTextFormatter(false, false)
	l := NewLogger(DebugLevel, j)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("debug", Data{"key": 1, "obj": testObject})
		l.Info("info", Data{"key": 1, "obj": testObject})
		l.Warn("warn", Data{"key": 1, "obj": testObject})
		l.Error("error", Data{"key": 1, "obj": testObject})
	}
	b.StopTimer()
}
