package log

import (
	"fmt"
	"time"

	"github.com/massarakhsh/lik"
)

type Level int

const (
	LEVEL_INFO  Level = 0
	LEVEL_EVENT Level = 1
	LEVEL_DEBUG Level = 2
)

var LogJson = false
var LogLevel Level = LEVEL_INFO

func SetLevel(lev Level, json bool) {
	LogLevel = lev
	LogJson = json
}

func SayInfo(format string, parms ...interface{}) {
	Say(LEVEL_INFO, format, parms...)
}

func SayEvent(format string, parms ...interface{}) {
	Say(LEVEL_EVENT, format, parms...)
}

func SayDebug(format string, parms ...interface{}) {
	Say(LEVEL_DEBUG, format, parms...)
}

func Say(lev Level, format string, parms ...interface{}) {
	if lev <= LogLevel {
		at := time.Now().Format("2006-01-02T15:04:05.000")
		text := fmt.Sprintf(format, parms...)
		if LogJson {
			//fmt.Printf("%s: %s\n", at, text)
			set := lik.BuildSet("at", at, "text", text)
			fmt.Printf("%s\n", set.Serialize())
		} else {
			fmt.Printf("%s: %s\n", at, text)
		}
	}
}
