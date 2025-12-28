package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	logOnce sync.Once
	std     *log.Logger
)

const logPath = "log/gateway.log"

// Logger returns a singleton logger writing to both stdout and log/gateway.log.
func Logger() *log.Logger {
	logOnce.Do(func() {
		dir := filepath.Dir(logPath)
		_ = os.MkdirAll(dir, 0755)

		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			std = log.New(os.Stdout, "[gateway] ", log.LstdFlags|log.Lshortfile)
			return
		}

		mw := io.MultiWriter(os.Stdout, f)
		std = log.New(mw, "[gateway] ", log.LstdFlags|log.Lshortfile)
	})

	return std
}

func Infof(format string, args ...interface{}) {
	Logger().Printf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Logger().Printf(format, args...)
}
