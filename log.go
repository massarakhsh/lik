package lik

import (
	"fmt"
	"os"
	"time"
)

var (
	logDebug   = 0
	logMax     = 2
	logDir     = "var/log"
	logToday   string
	protoLog   []string
	maxProto   = 1000
	indexProto = 0
)

func SetLevelErr() {
	logDebug = 0
}

func SetLevelWar() {
	logDebug = 1
}

func SetLevelInf() {
	logDebug = 2
}

func GetProtoIndex() int {
	return indexProto
}

func GetProtoLog(length int) []string {
	if size := len(protoLog); length >= size {
		return protoLog
	} else {
		return protoLog[size-length:]
	}
}

func SayError(text string) {
	sayToLog(0, text)
}
func SayWarning(text string) {
	sayToLog(1, text)
}
func SayInfo(text string) {
	sayToLog(2, text)
}

func sayToLog(lev int, text string) {
	day := time.Now().Format("2006/01/02")
	line := day + time.Now().Format(" 15:04:05 ") + text
	if lep := len(protoLog); lep >= maxProto {
		protoLog = protoLog[lep-maxProto+1:]
	}
	protoLog = append(protoLog, line)
	if lev > logDebug {
		return
	}

	fmt.Println(line)
	os.MkdirAll(logDir, os.ModePerm)
	if logToday == "" {
		if stat, err := os.Stat(logDir + "/info"); err == nil {
			if last := stat.ModTime().Format("2006/01/02"); last != day {
				logRoll()
			}
		}
	} else if day != logToday {
		logRoll()
	}
	logToday = day
	indexProto++
	for l := lev; l <= logDebug && l < 3; l++ {
		file := logDir
		if l == 0 {
			file += "/error"
		} else if l == 1 {
			file += "/warning"
		} else if l == 2 {
			file += "/info"
		}
		if f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			f.WriteString(line + "\n")
			f.Close()
		}
	}
}

func logRoll() {
	logRollFile("error")
	logRollFile("warning")
	logRollFile("info")
}

func logRollFile(name string) {
	md := 9
	path := logDir + "/" + name
	file := path + "." + IntToStr(md)
	if _, err := os.Stat(file); err == nil {
		os.Remove(file)
	}
	for nd := md - 1; nd >= 0; nd-- {
		filefrom := path
		if nd > 0 {
			filefrom += "." + IntToStr(nd)
		}
		if _, err := os.Stat(filefrom); err == nil {
			os.Rename(filefrom, file)
		}
		file = filefrom
	}
}
