// logger
package main

import (
	"errors"
	"os"
	"time"

	log "github.com/omidnikta/nlog"
)

func main() {
	logger := log.NewLogger()
	logger.SetLevel(log.DebugLevel)
	f := log.NewTEXTFormatter()
	logger.SetFormatter(f)
	logger.SetOut(os.Stdout)
	logger.SetShowCaller(true)

	check(logger)
	logger.SetFormatter(&log.JSONFormatter{TimestampFormat: time.Stamp})
	check(logger)
	log.SetIsTerminal(false)
	f = log.NewTEXTFormatter()
	logger.SetFormatter(f)
	check(logger)
}

func check(logger *log.Logger) {
	m := logger.New("Data", log.Data{"first": 1, "Second": "second"})
	m.Debug("Hello", nil)
	n := m.NewNode("parent", log.Data{"alef": nil, "yek": errors.New("my error")})
	n.Debug("Hello", nil)
	n.Info("Hello", nil)
	n.Warn("Hello", nil)
	n.Error("Hello", nil)
	n.Debugf("Hello", 32)
	n.Infof("Hello", 33)
	n.Warnf("Hello", 34)
	n.Errorf("Hello", 35)
	logger.Debug("Hello", nil)
	logger.Info("Hello", nil)
	logger.Warn("Hello", nil)
	logger.Error("Hello", log.Data{"alef": nil, "yek": errors.New("my error")})
	logger.Debugf("Hello", 2, 32)
	logger.Infof("Hello", 2, 32)
	logger.Warnf("Hello", 43, "i am")
	logger.Errorf("Hello", 56)
}
