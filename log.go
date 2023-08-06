package startprompt

import (
	"fmt"
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
)

type caller struct {
	filename string
	line     int
	name     string
}

func getCaller() *caller {
	pc, filename, line, ok := runtime.Caller(3)
	if ok {
		return &caller{
			filename: filename,
			line:     line,
			name:     runtime.FuncForPC(pc).Name(),
		}
	} else {
		return &caller{}
	}

}

func DebugLog(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	//c := getCaller()
	//log.WithFields(log.Fields{
	//	"method":   c.name,
	//	"filename": fmt.Sprintf("%s:%d", c.filename, c.line),
	//}).Debug(msg)
	log.Debug(msg)
}

func enableDebugLog() {
	log.SetLevel(log.DebugLevel)
}

func disableDebugLog() {
	log.SetLevel(log.PanicLevel)
}

func init() {
	file, err := os.OpenFile("startprompt.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}
