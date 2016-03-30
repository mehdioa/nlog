// +build darwin dragonfly freebsd netbsd openbsd

// constants
package nlog

// Taken from x/ssh/terminal/util_bsd.go
import "syscall"

const ioctlReadTermios = syscall.TIOCGETA
