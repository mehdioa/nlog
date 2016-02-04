package nlog

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
)

const DefaultTimestampFormat = "2006-01-02 15:04:05"

type Formatter interface {
	Format(*message, *bytes.Buffer)
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
