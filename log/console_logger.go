package log

import (
	"fmt"
	"log"
	"time"
)

const (
	infoLvl  = "info"
	errorLvl = "error"
)

func Info(format string, args ...interface{}) {
	fmt.Println(wrapper(format, infoLvl, args...))
}

func Error(format string, args ...interface{}) {
	fmt.Println(wrapper(format, errorLvl, args...))
}

func Fatal(format string, args ...interface{}) {
	log.Fatalln(wrapper(format, infoLvl, args...))
}

func wrapper(txt string, lvl string, args ...interface{}) string {
	if len(args) > 0 {
		txt = fmt.Sprintf(txt, args...)
	}
	return fmt.Sprintf("[%s %s] %s", lvl, timeForLog(), txt)
}

func timeForLog() string {
	return time.Now().Format("2006.01.02 15:04:05")
}
