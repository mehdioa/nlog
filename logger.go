// logger
package nlog

import (
	//	"fmt"
	"io"
	"os"
	"sync"
)

type Logger struct {
	// The logs are `io.Copy`'d to this in a mutex. It's common to set this to a
	// file, or leave it default which is `os.Stderr`. You can also set this to
	// something more adventorous, such as logging to Kafka.
	out io.Writer
	// Hooks for the logger instance. These allow firing events based on logging
	// levels and log entries. For example, to send errors to an error tracking
	// service, log to StatsD or dump the core on fatal errors.
	hooks LevelHooks
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

func NewLogger() *Logger {
	return &Logger{
		out:        os.Stderr,
		formatter:  new(textFormatter),
		hooks:      make(LevelHooks),
		level:      InfoLevel,
		showCaller: true,
	}
}

func (logger *Logger) New(data Data) *node {
	return &node{logger: logger, parent: nil, data: data, key: "Data"}
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

func (logger *Logger) Debug(msg string, ctx ...interface{}) {
	if logger.level >= DebugLevel {
		log(&message{msg: &msg, level: DebugLevel}, &node{logger: logger})
	}
}
func (logger *Logger) Info(msg string, ctx ...interface{}) {
	if logger.level >= InfoLevel {
		log(&message{msg: &msg, level: InfoLevel}, &node{logger: logger})
	}
}
func (logger *Logger) Warn(msg string, ctx ...interface{}) {
	if logger.level >= WarnLevel {
		log(&message{msg: &msg, level: WarnLevel}, &node{logger: logger})
	}
}
func (logger *Logger) Error(msg string, ctx ...interface{}) {
	if logger.level >= ErrorLevel {
		log(&message{msg: &msg, level: ErrorLevel}, &node{logger: logger})
	}
}
