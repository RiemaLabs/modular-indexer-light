// Package log ...
package log

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"go.uber.org/atomic"
)

const (
	colorRed = iota + 91
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
)

const (
	LevelOff = iota
	LevelError
	LevelWarn
	LevelDebug
	LevelVerbose
)

var v = log.New(os.Stdout, "\x1b["+strconv.Itoa(colorMagenta)+"m[V]\x1b[0m ", log.Ldate|log.Ltime|log.Lmicroseconds)

var i = log.New(os.Stdout, "\x1b["+strconv.Itoa(colorGreen)+"m[I]\x1b[0m ", log.Ldate|log.Ltime|log.Lmicroseconds)

var d = log.New(os.Stdout, "\x1b["+strconv.Itoa(colorBlue)+"m[D]\x1b[0m ", log.Ldate|log.Ltime|log.Lmicroseconds)

var w = log.New(os.Stderr, "\x1b["+strconv.Itoa(colorYellow)+"m[W]\x1b[0m ", log.Ldate|log.Ltime|log.Lmicroseconds)

var e = log.New(os.Stderr, "\x1b["+strconv.Itoa(colorRed)+"m[E]\x1b[0m ", log.Ldate|log.Ltime|log.Lmicroseconds)

var lvl atomic.Int32

var ver atomic.String

var build atomic.String

// SetVerion ...
func SetVerion(v, b string) {
	ver.Store(v)
	build.Store(b)
}

// SetLevel ...
func SetLevel(level int) {
	lvl.Store(int32(level))
}

// toArgs ...
func toArgs(c int, a ...interface{}) string {
	n := len(a)
	if n == 0 {
		return ""
	}
	if n == 1 {
		return fmt.Sprintf("%v", a[0])
	}
	kvs := []string{}
	if (n % 2) != 0 {
		a = a[0 : n-1]
	}
	for i := 0; i < len(a); i = i + 2 {
		kvs = append(kvs, fmt.Sprintf("\x1b["+strconv.Itoa(c)+"m%v\x1b[0m=%v", a[i], a[i+1]))
	}
	return strings.Join(kvs, ", ")
}

// Verbose ...
func Verbose(tag string, a ...interface{}) {
	if lvl.Load() >= LevelVerbose {
		v.Println(fmt.Sprintf("[%s.%s][ %24s ]", ver.Load(), build.Load(), tag), toArgs(colorMagenta, a...))
	}
}

// Debug ...
func Debug(tag string, a ...interface{}) {
	if lvl.Load() >= LevelDebug {
		d.Println(fmt.Sprintf("[%s.%s][ %24s ]", ver.Load(), build.Load(), tag), toArgs(colorBlue, a...))
	}
}

// Warn ...
func Warn(tag string, a ...interface{}) {
	if lvl.Load() >= LevelWarn {
		w.Println(fmt.Sprintf("[%s.%s][ %24s ]", ver.Load(), build.Load(), tag), toArgs(colorYellow, a...))
	}
}

// Info ...
func Info(tag string, a ...interface{}) {
	if lvl.Load() >= LevelError {
		i.Println(fmt.Sprintf("[%s.%s][ %24s ]", ver.Load(), build.Load(), tag), toArgs(colorGreen, a...))
	}
}

// Error ...
func Error(tag string, a ...interface{}) {
	if lvl.Load() >= LevelError {
		e.Println(fmt.Sprintf("[%s.%s][ %24s ]", ver.Load(), build.Load(), tag), toArgs(colorRed, a...))
	}
}

// Panicf ...
func Panicf(err error) {
	e.Panicf("%+v", err)
}
