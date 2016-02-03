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

	showCaller bool
}

const keyString = "Node"

func NewLogger() *Logger {
	return &Logger{
		out:        os.Stderr,
		formatter:  new(textFormatter),
		level:      InfoLevel,
		showCaller: true,
	}
}

func (logger *Logger) New(key string, data Data) *node {
	return &node{logger: logger, Node: nil, Data: data, key: key}
}

func (logger *Logger) SetFormatter(formatter Formatter) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.formatter = formatter
}

func (logger *Logger) SetOut(out io.Writer) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.out = out
}

func (logger *Logger) SetLevel(level Level) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.level = level
}

func (logger *Logger) SetShowCaller(b bool) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.showCaller = b
}

func (logger *Logger) Debug(msg string, data Data) {
	if logger.level >= DebugLevel {
		log(&message{Message: &msg, Level: DebugLevel, Data: data, Node: nil, logger: logger})
	}
}
func (logger *Logger) Info(msg string, data Data) {
	if logger.level >= InfoLevel {
		log(&message{Message: &msg, Level: InfoLevel, Data: data, Node: nil, logger: logger})
	}
}
func (logger *Logger) Warn(msg string, data Data) {
	if logger.level >= WarnLevel {
		log(&message{Message: &msg, Level: WarnLevel, Data: data, Node: nil, logger: logger})
	}
}
func (logger *Logger) Error(msg string, data Data) {
	if logger.level >= ErrorLevel {
		log(&message{Message: &msg, Level: ErrorLevel, Data: data, Node: nil, logger: logger})
	}
}
func (logger *Logger) Debugf(f string, args ...interface{}) {
	if logger.level >= DebugLevel {
		m := fmt.Sprintf(f, args...)
		log(&message{Message: &m, Level: DebugLevel, Data: nil, Node: nil, logger: logger})
	}
}
func (logger *Logger) Infof(f string, args ...interface{}) {
	if logger.level >= InfoLevel {
		m := fmt.Sprintf(f, args...)
		log(&message{Message: &m, Level: InfoLevel, Data: nil, Node: nil, logger: logger})
	}
}
func (logger *Logger) Warnf(f string, args ...interface{}) {
	if logger.level >= WarnLevel {
		m := fmt.Sprintf(f, args...)
		log(&message{Message: &m, Level: WarnLevel, Data: nil, Node: nil, logger: logger})
	}
}
func (logger *Logger) Errorf(f string, args ...interface{}) {
	if logger.level >= ErrorLevel {
		m := fmt.Sprintf(f, args...)
		log(&message{Message: &m, Level: ErrorLevel, Data: nil, Node: nil, logger: logger})
	}
}
