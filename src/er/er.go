package er

import (
	"fmt"
	"runtime"
	"strings"
)

type _Er struct {
	cause   error
	message string
	file    string
	line    int
	details []interface{}
}

func path_to_src(filename string) string {
	i := strings.Index(filename, "src/")
	if i > 0 {
		return filename[i+4:]
	}
	return filename
}

func (e _Er) Error() string {
	line := fmt.Sprintf("%s:%d", path_to_src(e.file), e.line)
	if e.message != "" {
		line += " " + e.message
	}
	for i := 0; i+1 < len(e.details); i += 2 {
		line += fmt.Sprintf(" %v=%+v", e.details[i], e.details[i+1])
	}
	if e.cause != nil {
		line += " == " + e.cause.Error()
	}
	return line
}

func Er(cause error, message string, details ...interface{}) error {
	_, file, line, _ := runtime.Caller(1)
	return _Er{cause, message, file, line, details}
}
