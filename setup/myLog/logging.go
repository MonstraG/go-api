package myLog

import (
	"fmt"
	"log"
	"os"
)

var flags = log.Ldate | log.Ltime | log.Lshortfile | log.LUTC

var errorLogger = log.New(os.Stdout, "log: ", flags)
var infoLogger = log.New(os.Stdout, "err: ", flags)
var fatalLogger = log.New(os.Stdout, "FATAL: ", flags)

func Log(skip int, message string) {
	_ = infoLogger.Output(skip+2, message)
}

func Logf(skip int, format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	Log(skip+1, message)
}

func Error(skip int, message string) {
	_ = errorLogger.Output(skip+2, message)
}

func Errorf(skip int, format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	Errorf(skip+1, message)
}

func Fatal(skip int, message string) {
	_ = fatalLogger.Output(skip+2, message)
	os.Exit(1)
}
