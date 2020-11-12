package logger

import (
	"log"
	"os"
)

var (
	//Info logger
	Info *log.Logger
	//Debug logger
	Debug *log.Logger
	//Warn logger
	Warn *log.Logger
	//Error logger
	Error *log.Logger
)

func init() {
	flags := log.Ldate | log.Ltime | log.Lshortfile
	Info = log.New(os.Stdout, "[INFO ] ", flags)
	Warn = log.New(os.Stderr, "[WARN] ", flags)
	Error = log.New(os.Stderr, "[ERROR] ", flags)
	Debug = log.New(os.Stderr, "[DEBUG] ", flags)
}
