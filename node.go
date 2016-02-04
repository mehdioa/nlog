// node
package nlog

import (
	"fmt"
)

type Data map[string]interface{}

type message struct {
	Level
	Time    string
	Message *string
	Data
	Node *Node
}

type Node struct {
	key string
	Data
	*Node
	logger *Logger
}

type _message struct {
	Level   string
	Time    string
	Message *string
	Data
	Node *Node
}
type _sc_message struct {
	Level   string
	Time    string
	Message *string
	Caller  string
	Data
	Node *Node
}

func (n *Node) NewNode(key string, data Data) *Node {
	return &Node{key: key, Data: data, Node: n, logger: n.logger}
}

func (n *Node) Debug(msg string, data Data) {
	if n.logger.level >= DebugLevel {
		n.logger.log(&message{Message: &msg, Level: DebugLevel, Node: n, Data: data})
	}
}
func (n *Node) Info(msg string, data Data) {
	if n.logger.level >= InfoLevel {
		n.logger.log(&message{Message: &msg, Level: InfoLevel, Node: n, Data: data})
	}
}
func (n *Node) Warn(msg string, data Data) {
	if n.logger.level >= WarnLevel {
		n.logger.log(&message{Message: &msg, Level: WarnLevel, Node: n, Data: data})
	}
}
func (n *Node) Error(msg string, data Data) {
	if n.logger.level >= ErrorLevel {
		n.logger.log(&message{Message: &msg, Level: ErrorLevel, Node: n, Data: data})
	}
}
func (n *Node) Debugf(f string, args ...interface{}) {
	if n.logger.level >= DebugLevel {
		msg := fmt.Sprintf(f, args...)
		n.logger.log(&message{Message: &msg, Level: DebugLevel, Node: n, Data: nil})
	}
}
func (n *Node) Infof(f string, args ...interface{}) {
	if n.logger.level >= InfoLevel {
		msg := fmt.Sprintf(f, args...)
		n.logger.log(&message{Message: &msg, Level: InfoLevel, Node: n, Data: nil})
	}
}
func (n *Node) Warnf(f string, args ...interface{}) {
	if n.logger.level >= WarnLevel {
		msg := fmt.Sprintf(f, args...)
		n.logger.log(&message{Message: &msg, Level: WarnLevel, Node: n, Data: nil})
	}
}
func (n *Node) Errorf(f string, args ...interface{}) {
	if n.logger.level >= ErrorLevel {
		msg := fmt.Sprintf(f, args...)
		n.logger.log(&message{Message: &msg, Level: ErrorLevel, Node: n, Data: nil})
	}
}
