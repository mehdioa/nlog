// constants
package nlog

import (
	"syscall"
	"unsafe"
)

var pool = NewBufferPool()

// Level type
type Level uint8

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel

	lastIndexLevel
)

var levelString = map[Level]string{
	DebugLevel: "DEBU",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERRO",
	PanicLevel: "PANI",
}
var levelColor = map[Level]int{
	DebugLevel: 37,
	InfoLevel:  34,
	WarnLevel:  33,
	ErrorLevel: 31,
	PanicLevel: 31,
}

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
//func (level Level) String() string {
//	switch level {
//	case DebugLevel:
//		return "debug"
//	case InfoLevel:
//		return "info"
//	case WarnLevel:
//		return "warning"
//	case ErrorLevel:
//		return "error"
//	case PanicLevel:
//		return "panic"
//	}

//	return "unknown"
//}

var isTerminal bool

var (
	headerFormat string
	msgFormat    string
	parentFormat string
	dataFormat   string
)

func init() {
	isTerminal = checkIsTerminal()
	//	if isTerminal {
	//		headerFormat = "\x1b[%dm%s\x1b[0m[%s] %-44s \x1b[%dmcaller\x1b[0m=%s "
	//		msgFormat = "\x1b[%dm%s\x1b[0m[%s] %-44s "
	//		parentFormat = "\x1b[%dm%s\x1b[0m={"
	//		dataFormat = "\x1b[%dm%s\x1b[0m=%+v "

	//	} else {
	//		headerFormat = "%s[%s] %-44s caller=%s "
	//		msgFormat = "%s[%s] %-44s "
	//		parentFormat = "%s={"
	//		dataFormat = "%s=%+v "
	//	}
}

func checkIsTerminal() bool {
	fd := syscall.Stderr
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TCGETS, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}

func SetIsTerminal(b bool) {
	isTerminal = b
}
