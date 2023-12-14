package lik

import (
	"sync"
	"time"
)

type ItGate struct {
	gateSync sync.Mutex
	ListElm  []ItGateElm
}

type ItGateElm struct {
	startAt  time.Time
	duration time.Duration
	limit    int64
	count    int64
}

func (it *ItGate) AddLimit(duration time.Duration, limit int64) {
	it.gateSync.Lock()
	defer it.gateSync.Unlock()

	elm := ItGateElm{startAt: time.Now(), duration: duration, limit: limit}
	it.ListElm = append(it.ListElm, elm)
}

func (it *ItGate) Probe() bool {
	it.gateSync.Lock()
	defer it.gateSync.Unlock()

	ok := true
	for n := 0; n < len(it.ListElm); n++ {
		elm := &it.ListElm[n]
		if time.Since(elm.startAt) >= elm.duration {
			elm.startAt = time.Now()
			elm.count = 0
		} else if elm.count >= elm.limit {
			ok = false
		} else {
			elm.count++
		}

	}
	return ok
}
