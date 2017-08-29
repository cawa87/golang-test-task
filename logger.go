package main

import (
	"log"
	"os"
)

//Logger Logger
type Logger struct {
	file         *os.File
	debugEnabled bool
	traceEnabled bool
}

//NewLogger NewLogger
func NewLogger(logFileName string, verbose, debugEnabled, traceEnabled bool) *Logger {
	var logr Logger
	logr.debugEnabled = debugEnabled
	logr.traceEnabled = traceEnabled
	//FIXME:verbose to console and to file if file present
	if verbose {
		log.SetOutput(os.Stdout)
	} else if logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777); err != nil {
		log.Fatal("opening log file error. file:", logFileName, " err:", err)
	} else {
		log.SetOutput(logFile)
		logr.file = logFile
	}
	return &logr
}

func (logger *Logger) flush() {
	logger.file.Sync()
}

//Trace Trace
func (logger *Logger) Trace(arg ...interface{}) {
	if logger.traceEnabled {
		log.Println(arg...)
		logger.flush()
	}
}

//Debug Debug
func (logger *Logger) Debug(arg ...interface{}) {
	if logger.debugEnabled {
		log.Println(arg...)
		logger.flush()
	}
}

//Info Info
func (logger *Logger) Info(arg ...interface{}) {
	log.Println(arg...)
	logger.flush()
	//glog.Infoln(arg...)
}

//Infof Infof
func (logger *Logger) Infof(message string, arg ...interface{}) {
	log.Printf(message+"\n", arg...)
	logger.flush()
	//glog.Infoln(arg...)
}

//Warning Warning
func (logger *Logger) Warning(arg ...interface{}) {
	//glog.Warningln(arg...)
	log.Println(arg...)
	logger.flush()
}

//Error Error
func (logger *Logger) Error(arg ...interface{}) {
	//glog.Errorln(arg...)
	log.Println(arg...)
	logger.flush()
}

//Fatal Fatal
func (logger *Logger) Fatal(arg ...interface{}) {
	//glog.Fatal(arg...)
	log.Println(arg...)
	logger.flush()
}
