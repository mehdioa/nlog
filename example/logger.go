// logger
package main

import (
	"errors"

	log "github.com/omidnikta/nlog"
)

func main() {
	f1 := log.NewTextFormatter(true, true)
	l1 := log.NewLogger(log.DebugLevel, f1)
	check(l1)

	f2 := log.NewJsonFormatter(true)
	l2 := log.NewLogger(log.DebugLevel, f2)
	check(l2)

	f3 := log.NewTextFormatter(false, false)
	l3 := log.NewLogger(log.DebugLevel, f3)
	check(l3)
}

func check(logger *log.Logger) {
	m := logger.New("Parent", log.Data{"first": 1, "Second": "second"})
	m.Debug("Hello", nil)
	n := m.NewNode("Child", log.Data{"alef": nil, "yek": errors.New("my error")})
	n.Debug("Hello", nil)
	n.Info("Hello", nil)
	n.Warn("Hello", nil)
	n.Error("Hello", nil)
	n.Debugf("Hello %d", 32)
	n.Infof("Hello %d", 33)
	n.Warnf("Hello %d", 34)
	n.Errorf("Hello %d", 35)
	logger.Debug("Hello", nil)
	logger.Info("Hello", nil)
	logger.Warn("Hello", nil)
	logger.Error("Hello", log.Data{"alef": nil, "yek": errors.New("my error")})
	logger.Debugf("Hello %d %s", 2, "aaioeh")
	logger.Infof("Hello %d %d", 2, 32)
	logger.Warnf("Hello %d %s", 43, "i am")
	logger.Errorf("Hello %d", 56)
}
