package nlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

const DefaultTimestampFormat = "2006-01-02 15:04:05"

type Formatter interface {
	Format(*message, *bytes.Buffer) error
}

type textFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
	fmt             func(*message, *bytes.Buffer) error
}

func (f *textFormatter) Format(msg *message, buf *bytes.Buffer) error {
	return f.fmt(msg, buf)
}

func NewTEXTFormatter() *textFormatter {
	t := &textFormatter{TimestampFormat: DefaultTimestampFormat}
	formattedTime := time.Now().Format(t.TimestampFormat)
	if isTerminal {
		t.fmt = func(msg *message, buf *bytes.Buffer) (err error) {
			ls := levelString[msg.Level]
			if msg != nil {
				lc := levelColor[msg.Level]
				if msg.logger.showCaller {
					_, err = fmt.Fprintf(buf, "\x1b[%dm%s\x1b[0m[%s] %-44s \x1b[%dmcaller\x1b[0m=%s", lc, ls, formattedTime, *msg.Message, lc, caller(5))
				} else {
					_, err = fmt.Fprintf(buf, "\x1b[%dm%s\x1b[0m[%s] %-44s", lc, ls, formattedTime, *msg.Message)
				}
				if err != nil {
					return
				}
				if msg.Data != nil && len(msg.Data) > 0 {
					if _, err = fmt.Fprintf(buf, " \x1b[%dm%s\x1b[0m={", lc, keyString); err != nil {
						return
					}
					first := true
					for k, v := range msg.Data {
						if first {
							_, err = fmt.Fprintf(buf, "\x1b[%dm%s\x1b[0m=%+v", lc, k, v)
							first = false
						} else {
							_, err = fmt.Fprintf(buf, " \x1b[%dm%s\x1b[0m=%+v", lc, k, v)
						}
						if err != nil {
							return
						}
					}
					if err = buf.WriteByte('}'); err != nil {
						return
					}
				}

				nd := msg.Node
				i := 0
				for nd != nil {
					i = i + 1
					if _, err = fmt.Fprintf(buf, " \x1b[%dm%s\x1b[0m={", lc, nd.key); err != nil {
						return
					}
					first := true
					for k, v := range nd.Data {
						if first {
							_, err = fmt.Fprintf(buf, "\x1b[%dm%s\x1b[0m=%+v", lc, k, v)
							first = false
						} else {
							_, err = fmt.Fprintf(buf, " \x1b[%dm%s\x1b[0m=%+v", lc, k, v)
						}
						if err != nil {
							return
						}
					}
					nd = nd.Node
				}
				for j := 0; j < i; j++ {
					if err = buf.WriteByte('}'); err != nil {
						return
					}
				}
				err = buf.WriteByte('\n')
			}
			return
		}
	} else {
		t.fmt = func(msg *message, buf *bytes.Buffer) (err error) {
			err = nil
			if msg != nil {
				ls := levelString[msg.Level]
				if msg.logger.showCaller {
					_, err = fmt.Fprintf(buf, "%s[%s] %-44s caller=%s", ls, formattedTime, *msg.Message, caller(5))
				} else {
					_, err = fmt.Fprintf(buf, "%s[%s] %-44s", ls, formattedTime, *msg.Message)
				}
				if err != nil {
					return
				}
				if msg.Data != nil && len(msg.Data) > 0 {
					if _, err = fmt.Fprintf(buf, " %s={", keyString); err != nil {
						return
					}
					first := true
					for k, v := range msg.Data {
						if first {
							_, err = fmt.Fprintf(buf, "%s=%+v", k, v)
							first = false
						} else {
							_, err = fmt.Fprintf(buf, " %s=%+v", k, v)
						}
						if err != nil {
							return
						}
					}
					if err = buf.WriteByte('}'); err != nil {
						return
					}
				}
				nd := msg.Node
				i := 0
				for nd != nil {
					i = i + 1
					if _, err = fmt.Fprintf(buf, " %s={", nd.key); err != nil {
						return
					}
					first := true
					for k, v := range nd.Data {
						if first {
							_, err = fmt.Fprintf(buf, "%s=%+v", k, v)
							first = false
						} else {
							_, err = fmt.Fprintf(buf, " %s=%+v", k, v)
						}
						if err != nil {
							return
						}
					}
					nd = nd.Node
				}
				for j := 0; j < i; j++ {
					if err = buf.WriteByte('}'); err != nil {
						return
					}
				}
				err = buf.WriteByte('\n')
			}
			return
		}
	}
	return t
}

type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
}

func (f *JSONFormatter) Format(msg *message, buf *bytes.Buffer) (err error) {
	_msg := &_message{Time: time.Now().Format(f.TimestampFormat), Message: msg.Message, Level: levelString[msg.Level], Node: msg.Node}
	s, err := json.Marshal(_msg)
	_, err = buf.Write(s)
	err = buf.WriteByte('\n')
	return
}

func caller(depth int) (str string) {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		str = "???: ?"
	} else {
		str = fmt.Sprint(filepath.Base(file), ":", line)
	}
	return
}
