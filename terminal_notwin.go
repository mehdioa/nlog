// terminal_notwin
// Based on ssh/terminal:
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux darwin freebsd openbsd netbsd dragonfly

package nlog

import (
	"syscall"
	"unsafe"
)

func checkIsTerminal() bool {
	fd := syscall.Stderr
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}
