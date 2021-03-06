package util

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

func debugStream() io.Writer {
	if os.Getenv("DEBUG") == "true" {
		return os.Stderr
	}
	return ioutil.Discard
}

var dbg = log.New(debugStream(), "[DEBUG] ", log.Lshortfile)

var logger = log.New(os.Stderr, "", 0)
