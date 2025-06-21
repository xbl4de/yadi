package log

import (
	"log"
	"os"
)

const yadiPrefix = "yadi"

var _log = log.New(os.Stdout, yadiPrefix, log.LstdFlags)

var verboseLog = func() *log.Logger {
	if os.Getenv("YADI_DEBUG") != "" {
		return log.New(os.Stdout, yadiPrefix, log.LstdFlags)
	} else {
		devNull, err := os.Open(os.DevNull)
		if err != nil {
			panic(err)
		}
		return log.New(devNull, yadiPrefix, log.LstdFlags)
	}
}()

func Log(fmt string, args ...interface{}) {
	_log.Printf(fmt, args...)
}

func Verbose(fmt string, args ...interface{}) {
	verboseLog.Printf(fmt, args...)
}
