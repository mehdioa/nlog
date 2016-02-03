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
	TimestampFormat    string
	fmt                func(*message, *bytes.Buffer) error
	fmtFirst           [lastIndexLevel]string
	fmtFirstShowCaller [lastIndexLevel]string
	fmtSecond          [lastIndexLevel]string
	fmtThird1          [lastIndexLevel]string
	fmtThird2          [lastIndexLevel]string
}

func (f *textFormatter) Format(msg *message, buf *bytes.Buffer) error {
	return f.fmt(msg, buf)
}

func NewTEXTFormatter() *textFormatter {
	t := &textFormatter{TimestampFormat: DefaultTimestampFormat}
	suffix := ""
	prefix := ""
	if isTerminal {
		suffix = "\x1b[0m"
	}
	for i := PanicLevel; i < lastIndexLevel; i++ {
		lc := levelColor[i]
		if isTerminal {
			prefix = fmt.Sprintf("\x1b[%dm", lc)
		}
		t.fmtFirstShowCaller[i] = fmt.Sprintf("%s%s%s[%s] %s %scaller%s=%s", prefix, "%s", suffix, "%s", "%-44s", prefix, suffix, "%s")
		t.fmtFirst[i] = fmt.Sprintf("%s%s%s[%s] %s", prefix, "%s", suffix, "%s", "%-44s")
		t.fmtSecond[i] = fmt.Sprintf(" %s%s%s={", prefix, "%s", suffix)
		t.fmtThird1[i] = fmt.Sprintf("%s%s%s=%s", prefix, "%s", suffix, "%+v")
		t.fmtThird2[i] = fmt.Sprintf(" %s%s%s=%s", prefix, "%s", suffix, "%+v")
	}
	formattedTime := time.Now().Format(t.TimestampFormat)
	t.fmt = func(msg *message, buf *bytes.Buffer) (err error) {
		ls := levelString[msg.Level]
		l := msg.Level
		if msg != nil {
			//			lc := levelColor[msg.Level]
			if msg.logger.showCaller {
				_, err = fmt.Fprintf(buf, t.fmtFirstShowCaller[l], ls, formattedTime, *msg.Message, caller(5))
				//					_, err = fmt.Fprintf(buf, "\x1b[%dm%s\x1b[0m[%s] %-44s \x1b[%dmcaller\x1b[0m=%s", lc, ls, formattedTime, *msg.Message, lc, caller(5))
			} else {
				_, err = fmt.Fprintf(buf, t.fmtFirst[l], ls, formattedTime, *msg.Message)
				//					_, err = fmt.Fprintf(buf, "\x1b[%dm%s\x1b[0m[%s] %-44s", lc, ls, formattedTime, *msg.Message)
			}
			if err != nil {
				return
			}
			if msg.Data != nil && len(msg.Data) > 0 {
				if _, err = fmt.Fprintf(buf, t.fmtSecond[l], keyString); err != nil {
					return
				}
				first := true
				for k, v := range msg.Data {
					if first {
						_, err = fmt.Fprintf(buf, t.fmtThird1[l], k, v)
						//							_, err = fmt.Fprintf(buf, "\x1b[%dm%s\x1b[0m=%+v", lc, k, v)
						first = false
					} else {
						_, err = fmt.Fprintf(buf, t.fmtThird2[l], k, v)
						//							_, err = fmt.Fprintf(buf, " \x1b[%dm%s\x1b[0m=%+v", lc, k, v)
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
				if _, err = fmt.Fprintf(buf, t.fmtSecond[l], nd.key); err != nil {
					return
				}
				first := true
				for k, v := range nd.Data {
					if first {
						_, err = fmt.Fprintf(buf, t.fmtThird1[l], k, v)
						first = false
					} else {
						_, err = fmt.Fprintf(buf, t.fmtThird2[l], k, v)
						//							_, err = fmt.Fprintf(buf, " \x1b[%dm%s\x1b[0m=%+v", lc, k, v)
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
