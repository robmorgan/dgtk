package log

import (
	"fmt"
	"os"
)

var debug = os.Getenv("DEBUG") == "true"
var PrintLogLevel bool

func Debug(format string, i ...interface{}) {
	if debug {
		Log("DEBUG", format, i...)
	}
}

func Error(format string, i ...interface{}) {
	Log("ERROR", format, i...)
}

func Fatal(format string, i ...interface{}) {
	Error(format, i...)
	os.Exit(1)
}

func Info(format string, i ...interface{}) {
	Log("INFO ", format, i...)
}

func Log(tag, format string, i ...interface{}) {
	if PrintLogLevel {
		format = tag + " " + format
	}
	fmt.Printf(format+"\n", i...)
}
