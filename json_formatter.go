// json_formatter
package nlog

import (
	"bytes"
	"encoding/json"
	"time"
)

type jsonFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
	fmt             func(msg *message, buf *bytes.Buffer)
}

func NewJsonFormatter(show_caller bool) *jsonFormatter {
	j := &jsonFormatter{TimestampFormat: DefaultTimestampFormat}
	if show_caller {
		j.fmt = func(msg *message, buf *bytes.Buffer) {
			_msg := &_sc_message{Time: time.Now().Format(j.TimestampFormat), Message: msg.Message, Level: levelStringLower[msg.Level], Caller: caller(5), Node: msg.Node, Data: msg.Data}
			s, _ := json.Marshal(_msg)
			buf.Write(s)
			buf.WriteByte('\n')
		}

	} else {
		j.fmt = func(msg *message, buf *bytes.Buffer) {
			_msg := &_message{Time: time.Now().Format(j.TimestampFormat), Message: msg.Message, Level: levelStringLower[msg.Level], Node: msg.Node, Data: msg.Data}
			s, _ := json.Marshal(_msg)
			buf.Write(s)
			buf.WriteByte('\n')
		}
	}
	return j
}

func (f *jsonFormatter) Format(msg *message, buf *bytes.Buffer) {
	f.fmt(msg, buf)
}
