// node
package nlog

import (
	"fmt"
	"io"
	"os"
)

type Data map[string]interface{}

type message struct {
	Time    string
	Message *string
	//	caller string
	Level
	//	time   time.Time
	Data
	Node   *node
	logger *Logger
}

type node struct {
	key string
	Data
	Node   *node
	logger *Logger
}

type _message struct {
	Time    string
	Message *string
	Level   string
	Data
	Node *node
}

func (n *node) NewNode(key string, data Data) *node {
	return &node{key: key, Data: data, Node: n, logger: n.logger}
}

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines
func log(m *message) {
	//	m := &message{msg: msg, level: level, time: time.Now()}
	//	m.time = time.Now()
	buf := pool.Get()
	defer pool.Put(buf)

	err := m.logger.formatter.Format(m, buf)
	if err != nil {
		m.logger.mu.Lock()
		defer m.logger.mu.Unlock()
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
	}

	m.logger.mu.Lock()
	defer m.logger.mu.Unlock()

	_, err = io.Copy(m.logger.out, buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if m.Level <= PanicLevel {
		panic(&err)
	}
}

func (n *node) Debug(msg string, data Data) {
	if n.logger.level >= DebugLevel {
		log(&message{Message: &msg, Level: DebugLevel, logger: n.logger, Node: n, Data: data})
	}
}
func (n *node) Info(msg string, data Data) {
	if n.logger.level >= InfoLevel {
		log(&message{Message: &msg, Level: InfoLevel, logger: n.logger, Node: n, Data: data})
	}
}
func (n *node) Warn(msg string, data Data) {
	if n.logger.level >= WarnLevel {
		log(&message{Message: &msg, Level: WarnLevel, logger: n.logger, Node: n, Data: data})
	}
}
func (n *node) Error(msg string, data Data) {
	if n.logger.level >= ErrorLevel {
		log(&message{Message: &msg, Level: ErrorLevel, logger: n.logger, Node: n, Data: data})
	}
}
func (n *node) Debugf(args ...interface{}) {
	if n.logger.level >= DebugLevel {
		msg := fmt.Sprint(args...)
		log(&message{Message: &msg, Level: DebugLevel, logger: n.logger, Node: n, Data: nil})
	}
}
func (n *node) Infof(args ...interface{}) {
	if n.logger.level >= InfoLevel {
		msg := fmt.Sprint(args...)
		log(&message{Message: &msg, Level: InfoLevel, logger: n.logger, Node: n, Data: nil})
	}
}
func (n *node) Warnf(args ...interface{}) {
	if n.logger.level >= WarnLevel {
		msg := fmt.Sprint(args...)
		log(&message{Message: &msg, Level: WarnLevel, logger: n.logger, Node: n, Data: nil})
	}
}
func (n *node) Errorf(args ...interface{}) {
	if n.logger.level >= ErrorLevel {
		msg := fmt.Sprint(args...)
		log(&message{Message: &msg, Level: ErrorLevel, logger: n.logger, Node: n, Data: nil})
	}
}
