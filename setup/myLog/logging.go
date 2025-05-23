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

func (myLogger *MyLogger) output(skip int, message string) {
	_ = myLogger.logger.Output(skip+2, message)
	if myLogger.die {
		os.Exit(0)
	}
}

func (myLogger *MyLogger) SkipLog(skip int, message string) {
	myLogger.output(skip+1, message)
}

func (myLogger *MyLogger) Log(message string) {
	myLogger.SkipLog(1, message)
}

func (myLogger *MyLogger) Logf(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	myLogger.SkipLog(1, message)
}
