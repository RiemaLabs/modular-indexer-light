package logs

import (
	"io"
	"log"
	"os"
)

func newLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix, log.LstdFlags|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
}

func discard() *log.Logger {
	return log.New(io.Discard, "", 0)
}

var (
	Debug = discard()
	Info  = newLogger("INFO ")
	Warn  = newLogger("WARN ")
	Error = newLogger("ERROR ")
)

func SetDebug() {
	Debug = newLogger("DEBUG ")
}
