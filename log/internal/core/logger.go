package core

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type logWriter struct {
	io.Writer
	timeFormat string
}

func (w logWriter) Write(p []byte) (n int, err error) {
	timestamp := "[" + time.Now().UTC().Format(w.timeFormat) + "] "
	p = append([]byte(timestamp), p...)
	return w.Writer.Write(p)
}

const errConfig = "logger has not been configured!"

var lDebug *log.Logger = nil
var lInfo *log.Logger = nil
var lWarning *log.Logger = nil
var lError *log.Logger = nil

func newLogger(out io.Writer, prefix string, prefixColor string) *log.Logger {
	prefix = "[" + Color(prefix, prefixColor) + "] "
	return log.New(out, prefix, 0)
}

func Configure() {
	writer := &logWriter{
		os.Stdout,
		"2006-01-02 15:04:05 ",
	}

	lDebug = newLogger(writer, "DBUG", ColorGreen)
	lInfo = newLogger(writer, "INFO", ColorCyan)
	lWarning = newLogger(writer, "WARN", ColorYellow)
	lError = newLogger(writer, "ERR!", ColorRed)
}

func dlog(logger *log.Logger, location string, format string, a ...interface{}) {
	if logger == nil {
		log.Println(errConfig)
		return
	}

	message := "[" + location + "] " + fmt.Sprintf(format, a...)
	logger.Println(message)
}

func Debug(location string, format string, a ...interface{}) {
	dlog(lDebug, location, format, a...)
}

func Info(location string, format string, a ...interface{}) {
	dlog(lInfo, location, format, a...)
}

func Warning(location string, format string, a ...interface{}) {
	dlog(lWarning, location, format, a...)
}

func Error(location string, format string, a ...interface{}) {
	dlog(lError, location, format, a...)
}

func Fatal(location string, format string, a ...interface{}) {
	dlog(lError, location, format, a...)
	os.Exit(1)
}
