package yadi

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
