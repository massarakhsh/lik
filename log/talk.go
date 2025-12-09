package log

import (
	"sync"
	"time"
)

const periodClear = time.Second * 300

type Talk struct {
	initialized bool
	gate        sync.Mutex
	errAll      map[string]*talkErr
	lastClear   time.Time
}

type talkErr struct {
	key      string
	lastOn   time.Time
	nextFree time.Time
}

func (talk *Talk) ErrorSay(key string, dura time.Duration, say string, parms ...any) bool {
	if !talk.ErrorOn(key, dura) {
		return false
	}
	SayInfo(say, parms...)
	return true
}

func (talk *Talk) ErrorOn(key string, dura time.Duration) bool {
	talk.gate.Lock()
	defer talk.gate.Unlock()
	talk.control()

	er := talk.errAll[key]
	if er == nil {
		er = &talkErr{key: key}
		talk.errAll[key] = er
	}
	if er.isOn() {
		return false
	}
	er.setOn(dura)

	return true
}

func (talk *Talk) ErrorOff(key string) {
	talk.gate.Lock()
	defer talk.gate.Unlock()
	talk.control()

	if er := talk.errAll[key]; er != nil {
		er.nextFree = er.lastOn
	}
}

func (talk *Talk) control() {
	if !talk.initialized {
		talk.errAll = make(map[string]*talkErr)
		talk.lastClear = time.Now()
		talk.initialized = true
	} else if time.Since(talk.lastClear) > periodClear {
		talk.clear()
		talk.lastClear = time.Now()
	}
}

func (talk *Talk) clear() {
	for _, er := range talk.errAll {
		if !er.isOn() && time.Since(er.lastOn) > periodClear {
			delete(talk.errAll, er.key)
		}
	}
}

func (er *talkErr) isOn() bool {
	return time.Now().Before(er.nextFree)
}

func (er *talkErr) setOn(dura time.Duration) {
	er.lastOn = time.Now()
	er.nextFree = er.lastOn.Add(dura)
}
