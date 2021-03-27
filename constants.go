// constants
package nlog

var pool = NewBufferPool()

// Level type
type Level uint8

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	FatalLevel Level = iota
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel
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

var levelString = []string{"FATA", "PANI", "ERRO", "WARN", "INFO", "DEBU"}
var levelColor = []int{31, 31, 31, 33, 34, 37}
var levelStringLower = []string{"fatal", "panic", "error", "warn", "info", "debug"}

var isTerminal bool

func init() {
	isTerminal = checkIsTerminal()
}

func StringToLevel(level string) Level {
	switch level {
	case "panic":
		return PanicLevel
	case "error":
		return ErrorLevel
	case "warn":
		return WarnLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	}
	return FatalLevel
}
