package printlog

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

const NONE = -1
const ERROR = 0
const WARNING = 1
const INFO = 2
const DEBUG = 3

var TimestampFormat = "2006-01-02 15:04:05 -0700"
var Verbosity = INFO
var Writer io.Writer = os.Stderr

var _LEVEL_TO_NAME = []string{
	"none",
	"error",
	"warning",
	"info",
	"debug",
}

var _LEVEL_TO_SYSLOG = []int{
	7,
	3,
	4,
	5,
	7,
}

func Print(level int, message interface{}, details ...interface{}) {
	if Verbosity < level {
		return
	}

	line := ""
	line += "<" + strconv.Itoa(_LEVEL_TO_SYSLOG[level-NONE]) + ">"
	line += time.Now().Format(TimestampFormat)
	line += " [" + _LEVEL_TO_NAME[level-NONE] + "]"
	if s, ok := message.(string); ok {
		line += fmt.Sprintf(" " + s, details...)
	} else {
		line += fmt.Sprintf(" %v", message)
	}

	fmt.Fprintln(Writer, line)
}

func Error(message interface{}, details ...interface{}) {
	Print(ERROR, message, details...)
}

func Warning(message interface{}, details ...interface{}) {
	Print(WARNING, message, details...)
}

func Info(message interface{}, details ...interface{}) {
	Print(INFO, message, details...)
}

func Debug(message interface{}, details ...interface{}) {
	Print(DEBUG, message, details...)
}
