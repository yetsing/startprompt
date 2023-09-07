package startprompt

import (
	"fmt"
	"log"
	"os"
)

var debugLogger *log.Logger

func DebugLog(format string, a ...any) {
	if debugLogger != nil {
		msg := fmt.Sprintf(format, a...)
		_ = debugLogger.Output(2, msg)
	}
}

func enableDebugLog() {
	f, err := os.OpenFile("startprompt.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	debugLogger = log.New(f, "DEBUG ", log.Ldate|log.Ltime|log.Lshortfile)
}

func disableDebugLog() {
	debugLogger = nil
}
