// node
package nlog

import (
	"fmt"
	"io"
	"os"
	//	"time"
)

type Data map[string]interface{}

type message struct {
	msg *string
	//	caller string
	level Level
	//	time   time.Time
}

type node struct {
	key    string
	data   Data
	parent *node
	logger *Logger
}

func (n *node) NewNode(key string, data Data) *node {
	return &node{key: key, data: data, parent: n, logger: n.logger}
}

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines
func log(m *message, n *node) {
	//	m := &message{msg: msg, level: level, time: time.Now()}
	//	m.time = time.Now()
	buf := pool.Get()
	defer pool.Put(buf)

	err := n.logger.formatter.Format(n, m, buf)
	if err != nil {
		n.logger.mu.Lock()
		defer n.logger.mu.Unlock()
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
	}

	n.logger.mu.Lock()
	defer n.logger.mu.Unlock()

	_, err = io.Copy(n.logger.out, buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if m.level <= PanicLevel {
		panic(&err)
	}
}

func (n *node) Debug(msg string, data ...interface{}) {
	if n.logger.level >= DebugLevel {
		log(&message{msg: &msg, level: DebugLevel}, n)
	}
}
func (n *node) Info(msg string, data ...interface{}) {
	if n.logger.level >= InfoLevel {
		log(&message{msg: &msg, level: InfoLevel}, n)
	}
}
func (n *node) Warn(msg string, data ...interface{}) {
	if n.logger.level >= WarnLevel {
		log(&message{msg: &msg, level: WarnLevel}, n)
	}
}
func (n *node) Error(msg string, data ...interface{}) {
	if n.logger.level >= ErrorLevel {
		log(&message{msg: &msg, level: ErrorLevel}, n)
	}
}
