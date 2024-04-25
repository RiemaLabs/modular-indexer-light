package logs

import (
	"log"
	"os"
)

func newLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix, log.LstdFlags|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
}

var (
	Info  = newLogger("INFO ")
	Warn  = newLogger("WARN ")
	Error = newLogger("ERROR ")
)
