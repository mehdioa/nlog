// text_formatter
package nlog

import (
	"bytes"
	"fmt"
	"time"
)

type textFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat    string
	fmt                func(*message, *bytes.Buffer)
	fmtMsg             func(*message, *bytes.Buffer, Level, *string)
	fmtFirst           [lastIndexLevel]string
	fmtFirstShowCaller [lastIndexLevel]string
	fmtNode            [lastIndexLevel]string
	fmtData            [lastIndexLevel]string
}

func (f *textFormatter) Format(msg *message, buf *bytes.Buffer) {
	f.fmt(msg, buf)
}

func NewTextFormatter(show_caller, enable_color bool) *textFormatter {
	t := &textFormatter{TimestampFormat: DefaultTimestampFormat}
	defaultFG := ""
	levColFG := ""
	nodeFG := ""
	enable_color = enable_color && isTerminal
	if enable_color {
		defaultFG = "\x1b[39m"
		nodeFG = "\x1b[96m"
	}
	for i := FatalLevel; i < lastIndexLevel; i++ {
		lc := levelColor[i]
		if enable_color {
			levColFG = fmt.Sprintf("\x1b[%dm", lc)
		}
		t.fmtFirstShowCaller[i] = fmt.Sprintf("%s%s%s[%s] %s %scaller%s=%s", levColFG, "%s", defaultFG, "%s", "%-44s", levColFG, defaultFG, "%s")
		t.fmtFirst[i] = fmt.Sprintf("%s%s%s[%s] %s", levColFG, "%s", defaultFG, "%s", "%-44s")
		t.fmtNode[i] = fmt.Sprintf(" %s%s%s ", nodeFG, "%s", defaultFG)
		t.fmtData[i] = fmt.Sprintf("%s%s%s=%s ", levColFG, "%s", defaultFG, "%+v")
	}
	//	formattedTime := time.Now().Format(t.TimestampFormat)

	if show_caller {
		t.fmtMsg = func(msg *message, buf *bytes.Buffer, l Level, ls *string) {
			fmt.Fprintf(buf, t.fmtFirstShowCaller[l], *ls, time.Now().Format(t.TimestampFormat), *msg.message, caller(6))
		}
	} else {
		t.fmtMsg = func(msg *message, buf *bytes.Buffer, l Level, ls *string) {
			fmt.Fprintf(buf, t.fmtFirst[l], *ls, time.Now().Format(t.TimestampFormat), *msg.message)
		}
	}

	t.fmt = func(msg *message, buf *bytes.Buffer) {
		ls := levelString[msg.level]
		l := msg.level
		if msg != nil {
			t.fmtMsg(msg, buf, l, &ls)
			nd := &Node{key: keyString, node: msg.node, data: msg.data}
			for nd != nil {
				if nd.data != nil && len(nd.data) > 0 {
					fmt.Fprintf(buf, t.fmtNode[l], nd.key)
					for k, v := range nd.data {
						fmt.Fprintf(buf, t.fmtData[l], k, v)
					}
				}
				nd = nd.node
			}
			buf.WriteByte('\n')
		}
	}
	return t
}
