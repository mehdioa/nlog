// node
package nlog

import (
	"fmt"
	"os"
)

type Data map[string]interface{}

type message struct {
	level   Level
	time    string
	message *string
	data    Data
	node    *Node
}

type Node struct {
	key    string
	data   Data
	node   *Node
	logger *Logger
}

type _message struct {
	level   string
	time    string
	message *string
	data    Data
	node    *Node
}
type _sc_message struct {
	level   string
	time    string
	message *string
	caller  string
	data    Data
	node    *Node
}

func (n *Node) NewNode(key string, data Data) *Node {
	return &Node{key: key, data: data, node: n, logger: n.logger}
}

func (n *Node) Debug(msg string, data Data) {
	if n.logger.level >= DebugLevel {
		n.logger.log(&message{message: &msg, level: DebugLevel, node: n, data: data})
	}
}
func (n *Node) Info(msg string, data Data) {
	if n.logger.level >= InfoLevel {
		n.logger.log(&message{message: &msg, level: InfoLevel, node: n, data: data})
	}
}
func (n *Node) Warn(msg string, data Data) {
	if n.logger.level >= WarnLevel {
		n.logger.log(&message{message: &msg, level: WarnLevel, node: n, data: data})
	}
}
func (n *Node) Error(msg string, data Data) {
	if n.logger.level >= ErrorLevel {
		n.logger.log(&message{message: &msg, level: ErrorLevel, node: n, data: data})
	}
}
func (n *Node) Panic(msg string, data Data) {
	if n.logger.level >= PanicLevel {
		_msg := message{message: &msg, level: PanicLevel, node: n, data: data}
		n.logger.log(&_msg)
		panic(_msg)
	}
}
func (n *Node) Fatal(msg string, data Data) {
	if n.logger.level >= FatalLevel {
		_msg := message{message: &msg, level: FatalLevel, node: n, data: data}
		n.logger.log(&_msg)
		os.Exit(1)
	}
}
func (n *Node) Debugf(f string, args ...interface{}) {
	if n.logger.level >= DebugLevel {
		msg := fmt.Sprintf(f, args...)
		n.logger.log(&message{message: &msg, level: DebugLevel, node: n, data: nil})
	}
}
func (n *Node) Infof(f string, args ...interface{}) {
	if n.logger.level >= InfoLevel {
		msg := fmt.Sprintf(f, args...)
		n.logger.log(&message{message: &msg, level: InfoLevel, node: n, data: nil})
	}
}
func (n *Node) Warnf(f string, args ...interface{}) {
	if n.logger.level >= WarnLevel {
		msg := fmt.Sprintf(f, args...)
		n.logger.log(&message{message: &msg, level: WarnLevel, node: n, data: nil})
	}
}
func (n *Node) Errorf(f string, args ...interface{}) {
	if n.logger.level >= ErrorLevel {
		msg := fmt.Sprintf(f, args...)
		n.logger.log(&message{message: &msg, level: ErrorLevel, node: n, data: nil})
	}
}
func (n *Node) Panicf(f string, args ...interface{}) {
	if n.logger.level >= PanicLevel {
		msg := fmt.Sprintf(f, args...)
		_msg := message{message: &msg, level: PanicLevel, node: n, data: nil}
		n.logger.log(&_msg)
		panic(_msg)
	}
}
func (n *Node) Fatalf(f string, args ...interface{}) {
	if n.logger.level >= FatalLevel {
		msg := fmt.Sprintf(f, args...)
		_msg := message{message: &msg, level: FatalLevel, node: n, data: nil}
		n.logger.log(&_msg)
		os.Exit(1)
	}
}
