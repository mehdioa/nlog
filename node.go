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
	Node *node
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
	if m.Level <= PanicLevel {
		panic(&err)
	}
}

func (n *node) Debug(msg string, data Data) {
	if n.logger.level >= DebugLevel {
		log(&message{Message: &msg, Level: DebugLevel}, n)
	}
}
func (n *node) Info(msg string, data Data) {
	if n.logger.level >= InfoLevel {
		log(&message{Message: &msg, Level: InfoLevel}, n)
	}
}
func (n *node) Warn(msg string, data Data) {
	if n.logger.level >= WarnLevel {
		log(&message{Message: &msg, Level: WarnLevel}, n)
	}
}
func (n *node) Error(msg string, data Data) {
	if n.logger.level >= ErrorLevel {
		log(&message{Message: &msg, Level: ErrorLevel}, n)
	}
}
