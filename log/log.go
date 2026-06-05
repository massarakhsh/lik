package log

import (
	"fmt"
	"time"

	"github.com/massarakhsh/lik"
)

var LogJson = false

func SetJson(json bool) {
	LogJson = json
}

func SayInfo(format string, parms ...interface{}) {
	at := time.Now().Format("2006-01-02T15:04:05.000")
	text := fmt.Sprintf(format, parms...)
	if LogJson {
		//fmt.Printf("%s: %s\n", at, text)
		set := lik.BuildSet("at", at, "text", text)
		fmt.Printf("%s\r\n", set.Serialize())
	} else {
		fmt.Printf("%s: %s\r\n", at, text)
	}
}
