// logger
package nlog

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type Logger struct {
	// The logs are `io.Copy`'d to this in a mutex. It's common to set this to a
	// file, or leave it default which is `os.Stderr`. You can also set this to
	// something more adventorous, such as logging to Kafka.
	out io.Writer
	// All log entries pass through the formatter before logged to Out. The
	// included formatters are `TextFormatter` and `JSONFormatter` for which
	// TextFormatter is the default. In development (when a TTY is attached) it
	// logs with colors, but to a file it wouldn't. You can easily implement your
	// own that implements the `Formatter` interface, see the `README` or included
	// formatters for examples.
	formatter Formatter
	// The logging level the logger should log at. This is typically (and defaults
	// to) `logrus.Info`, which allows Info(), Warn(), Error() and Fatal() to be
	// logged. `logrus.Debug` is useful in
	level Level
	// Used to sync writing to the log.
	mu sync.Mutex
}

const keyString = "Node"

func NewLogger(level Level, f Formatter) *Logger {
	return &Logger{
		out:       os.Stdout,
		formatter: f,
		level:     level,
	}
}

func (logger *Logger) log(m *message) {
	buf := pool.Get()
	defer pool.Put(buf)

	logger.formatter.Format(m, buf)

	logger.mu.Lock()
	defer logger.mu.Unlock()

	io.Copy(logger.out, buf)
}

func (logger *Logger) New(key string, data Data) *Node {
	return &Node{logger: logger, Node: nil, Data: data, key: key}
}

func (logger *Logger) SetOut(out io.Writer) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.out = out
}

func (logger *Logger) Debug(msg string, data Data) {
	if logger.level >= DebugLevel {
		logger.log(&message{Message: &msg, Level: DebugLevel, Data: data, Node: nil})
	}
}
func (logger *Logger) Info(msg string, data Data) {
	if logger.level >= InfoLevel {
		logger.log(&message{Message: &msg, Level: InfoLevel, Data: data, Node: nil})
	}
}
func (logger *Logger) Warn(msg string, data Data) {
	if logger.level >= WarnLevel {
		logger.log(&message{Message: &msg, Level: WarnLevel, Data: data, Node: nil})
	}
}
func (logger *Logger) Error(msg string, data Data) {
	if logger.level >= ErrorLevel {
		logger.log(&message{Message: &msg, Level: ErrorLevel, Data: data, Node: nil})
	}
}
func (logger *Logger) Debugf(f string, args ...interface{}) {
	if logger.level >= DebugLevel {
		m := fmt.Sprintf(f, args...)
		logger.log(&message{Message: &m, Level: DebugLevel, Data: nil, Node: nil})
	}
}
func (logger *Logger) Infof(f string, args ...interface{}) {
	if logger.level >= InfoLevel {
		m := fmt.Sprintf(f, args...)
		logger.log(&message{Message: &m, Level: InfoLevel, Data: nil, Node: nil})
	}
}
func (logger *Logger) Warnf(f string, args ...interface{}) {
	if logger.level >= WarnLevel {
		m := fmt.Sprintf(f, args...)
		logger.log(&message{Message: &m, Level: WarnLevel, Data: nil, Node: nil})
	}
}
func (logger *Logger) Errorf(f string, args ...interface{}) {
	if logger.level >= ErrorLevel {
		m := fmt.Sprintf(f, args...)
		logger.log(&message{Message: &m, Level: ErrorLevel, Data: nil, Node: nil})
	}
}
