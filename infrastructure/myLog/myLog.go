package myLog

import (
	"fmt"
	"log"
	"os"
)

var flags = log.Ldate | log.Ltime | log.Lshortfile | log.LUTC

var Info = MyLogger{logger: log.New(os.Stdout, "log: ", flags), die: false}
var Error = MyLogger{logger: log.New(os.Stdout, "err: ", flags), die: false}
var Fatal = MyLogger{logger: log.New(os.Stdout, "FATAL: ", flags), die: true}

type MyLogger struct {
	logger *log.Logger
	die    bool
}

// lazySprintf will
func lazySprintf(format string, a ...any) string {
	if len(a) == 0 {
		return format
	}
	return fmt.Sprintf(format, a...)
}

func (myLogger *MyLogger) output(skip int, format string, a ...any) {
	message := lazySprintf(format, a...)
	_ = myLogger.logger.Output(skip+2, message)
	if myLogger.die {
		os.Exit(0)
	}
}

// SkipLog logs a message, but skips this amount of stack levels to show real code point
// only useful in helpers, otherwise just use Logf
func (myLogger *MyLogger) SkipLog(skip int, format string, a ...any) {
	myLogger.output(skip+1, format, a...)
}

// Logf logs a message
func (myLogger *MyLogger) Logf(format string, a ...any) {
	myLogger.SkipLog(1, format, a...)
}
