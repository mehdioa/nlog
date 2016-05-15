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
	return &Node{logger: logger, node: nil, data: data, key: key}
}

func (logger *Logger) SetOut(out io.Writer) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.out = out
}

func (logger *Logger) Debug(msg string, data Data) {
	if logger.level >= DebugLevel {
		logger.log(&message{message: &msg, level: DebugLevel, data: data, node: nil})
	}
}
func (logger *Logger) Info(msg string, data Data) {
	if logger.level >= InfoLevel {
		logger.log(&message{message: &msg, level: InfoLevel, data: data, node: nil})
	}
}
func (logger *Logger) Warn(msg string, data Data) {
	if logger.level >= WarnLevel {
		logger.log(&message{message: &msg, level: WarnLevel, data: data, node: nil})
	}
}
func (logger *Logger) Error(msg string, data Data) {
	if logger.level >= ErrorLevel {
		logger.log(&message{message: &msg, level: ErrorLevel, data: data, node: nil})
	}
}
func (logger *Logger) Panic(msg string, data Data) {
	if logger.level >= PanicLevel {
		_msg := message{message: &msg, level: PanicLevel, data: data, node: nil}
		logger.log(&_msg)
		panic(_msg)
	}
}
func (logger *Logger) Fatal(msg string, data Data) {
	if logger.level >= FatalLevel {
		_msg := message{message: &msg, level: FatalLevel, data: data, node: nil}
		logger.log(&_msg)
		os.Exit(1)
	}
}

func (logger *Logger) Debugf(f string, args ...interface{}) {
	if logger.level >= DebugLevel {
		m := fmt.Sprintf(f, args...)
		logger.log(&message{message: &m, level: DebugLevel, data: nil, node: nil})
	}
}
func (logger *Logger) Infof(f string, args ...interface{}) {
	if logger.level >= InfoLevel {
		m := fmt.Sprintf(f, args...)
		logger.log(&message{message: &m, level: InfoLevel, data: nil, node: nil})
	}
}
func (logger *Logger) Warnf(f string, args ...interface{}) {
	if logger.level >= WarnLevel {
		m := fmt.Sprintf(f, args...)
		logger.log(&message{message: &m, level: WarnLevel, data: nil, node: nil})
	}
}
func (logger *Logger) Errorf(f string, args ...interface{}) {
	if logger.level >= ErrorLevel {
		m := fmt.Sprintf(f, args...)
		logger.log(&message{message: &m, level: ErrorLevel, data: nil, node: nil})
	}
}
func (logger *Logger) Panicf(f string, args ...interface{}) {
	if logger.level >= PanicLevel {
		m := fmt.Sprintf(f, args...)
		logger.log(&message{message: &m, level: PanicLevel, data: nil, node: nil})
		panic(m)
	}
}
func (logger *Logger) Fatalf(f string, args ...interface{}) {
	if logger.level >= FatalLevel {
		m := fmt.Sprintf(f, args...)
		logger.log(&message{message: &m, level: FatalLevel, data: nil, node: nil})
		os.Exit(1)
	}
}
