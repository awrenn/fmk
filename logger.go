package fmk

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

var Log Logger

func init() {
	LogInit(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	return
}

func LogInit(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Log.SetTrace(traceHandle)
	Log.SetInfo(infoHandle)
	Log.SetWarning(warningHandle)
	Log.SetError(errorHandle)
}

func (logger *Logger) SetTrace(w io.Writer) {
	logger.Trace = log.New(w,
		"TRACE: ",
		log.Ldate|log.Ltime)
}

func (logger *Logger) SetInfo(w io.Writer) {
	logger.Info = log.New(w,
		"INFO: ",
		log.Ldate|log.Ltime)
}

func (logger *Logger) SetWarning(w io.Writer) {
	logger.Warning = log.New(w,
		"WARNING: ",
		log.Ldate|log.Ltime)
}

func (logger *Logger) SetError(w io.Writer) {
	logger.Error = log.New(w,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
