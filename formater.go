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
	Format(*node, *message, *bytes.Buffer) error
}

type textFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
	fmt             func(*node, *message, *bytes.Buffer) error
	fmtParent       func(nd *node, lc int, buf *bytes.Buffer, key string) error
}

func (f *textFormatter) Format(nd *node, msg *message, buf *bytes.Buffer) error {
	return f.fmt(nd, msg, buf)
}

func NewTEXTFormatter() *textFormatter {
	t := &textFormatter{TimestampFormat: DefaultTimestampFormat}
	if isTerminal {
		t.fmt = func(nd *node, msg *message, buf *bytes.Buffer) (err error) {
			if msg != nil {
				lc := levelColor[msg.Level]
				if nd.logger.showCaller {
					_, err = fmt.Fprintf(buf, "\x1b[%dm%s\x1b[0m[%s] %-44s \x1b[%dmcaller\x1b[0m=%s", lc, levelString[msg.Level], time.Now().Format(t.TimestampFormat), *msg.Message, lc, caller(5))
				} else {
					_, err = fmt.Fprintf(buf, "\x1b[%dm%s\x1b[0m[%s] %-44s", lc, levelString[msg.Level], time.Now().Format(t.TimestampFormat), *msg.Message)
				}
				if err != nil {
					return err
				}
				if len(nd.Data) > 0 || nd.Node != nil {
					err = t.fmtParent(nd, lc, buf, "data")
					if err != nil {
						return err
					}
				}
				err = buf.WriteByte('\n')
				if err != nil {
					return err
				}
			}
			return nil
		}

		t.fmtParent = func(nd *node, lc int, buf *bytes.Buffer, key string) (err error) {
			err = nil
			if nd != nil {
				_, err = fmt.Fprintf(buf, " \x1b[%dm%s\x1b[0m={", lc, key)
				if err != nil {
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
					if err != nil {
						return
					}
				}
				err = t.fmtParent(nd.Node, lc, buf, nd.key)
				if err != nil {
					return
				}
				_, err = buf.WriteString("}")
			}
			return
		}

	} else {
		t.fmt = func(nd *node, msg *message, buf *bytes.Buffer) (err error) {
			err = nil
			if msg != nil {
				lc := levelColor[msg.Level]
				if nd.logger.showCaller {
					_, err = fmt.Fprintf(buf, "%s[%s] %-44s caller=%s", levelString[msg.Level], time.Now().Format(t.TimestampFormat), *msg.Message, caller(5))
				} else {
					_, err = fmt.Fprintf(buf, "%s[%s] %-44s", levelString[msg.Level], time.Now().Format(t.TimestampFormat), *msg.Message)
				}
				if err != nil {
					return
				}
				if len(nd.Data) > 0 || nd.Node != nil {
					err = t.fmtParent(nd, lc, buf, "data")
					if err != nil {
						return
					}
				}
				err = buf.WriteByte('\n')
			}
			return
		}

		t.fmtParent = func(nd *node, lc int, buf *bytes.Buffer, key string) (err error) {
			err = nil
			if nd != nil {
				_, err = fmt.Fprintf(buf, " %s={", key)
				if err != nil {
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
				err = t.fmtParent(nd.Node, lc, buf, nd.key)
				if err != nil {
					return
				}
				_, err = buf.WriteString("}")
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

func (f *JSONFormatter) Format(nd *node, msg *message, buf *bytes.Buffer) (err error) {
	_msg := &_message{Time: time.Now().Format(f.TimestampFormat), Message: msg.Message, Level: levelString[msg.Level], Node: nd}
	s, err := json.Marshal(_msg)
	_, err = buf.Write(s)
	err = buf.WriteByte('\n')
	return
}

func (f *JSONFormatter) fmtParent(nd *node, buf *bytes.Buffer, key string) (err error) {
	err = nil
	if nd != nil {
		_, err = fmt.Fprintf(buf, `, "%s":{`, key)
		if err != nil {
			return
		}
		l := len(nd.Data)
		i := 0
		for k, v := range nd.Data {
			i = i + 1
			switch v := v.(type) {
			case string:
				_, err = fmt.Fprintf(buf, `"%s":"%+v"`, k, v)
			case error:
				_, err = fmt.Fprintf(buf, `"%s":"%s"`, k, v.Error())
			case nil:
				_, err = fmt.Fprintf(buf, `"%s":"nil"`, k)
			default:
				_, err = fmt.Fprintf(buf, `"%s":%+v`, k, v)
			}
			if err != nil {
				return
			}
			if i < l {
				_, err = buf.WriteString(", ")
			}
			if err != nil {
				return
			}
		}
		err = f.fmtParent(nd.Node, buf, nd.key)
		if err != nil {
			return
		}
		_, err = buf.WriteString("}")
	}
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
