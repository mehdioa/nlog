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

}

func check(logger *log.Logger) {
	m := logger.New(log.Data{"first": 1, "Second": "second"})
	m.Debug("Hello")
	n := m.NewNode("parent", log.Data{"alef": nil, "yek": errors.New("my error")})
	n.Debug("Hello")
	n.Info("Hello")
	n.Warn("Hello")
	n.Error("Hello")
	logger.Debug("Hello", nil)
	logger.Info("Hello", nil)
	logger.Warn("Hello", nil)
	logger.Error("Hello", log.Data{"alef": nil, "yek": errors.New("my error")})
}
