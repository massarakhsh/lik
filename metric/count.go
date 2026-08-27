package metric

import (
	"sync"
	"time"
)

const countSize = 15

type Count struct {
	gate      sync.Mutex
	isFreq    bool
	at        int64
	size      int
	line      [countSize]countElm
	lastValue float64
}

type countElm struct {
	value  float64
	factor int64
}

func (count *Count) Inc() {
	count.Add(1)
}

func (count *Count) Add(val float64) {
	count.SetMod(val, true)
}

func (count *Count) Set(val float64) {
	count.SetMod(val, false)
}

func (count *Count) SetMod(val float64, isFreq bool) {
	count.gate.Lock()
	defer count.gate.Unlock()

	count.prepare(isFreq)
	elm := &count.line[0]
	if count.isFreq {
		elm.value += val
	} else if elm.factor == 0 {
		elm.value = val
	} else {
		sum := elm.value * float64(elm.factor)
		sum += val
		elm.value = sum / float64(elm.factor+1)
	}
	elm.factor++
	count.lastValue = val
}

func (count *Count) Get() float64 {
	count.gate.Lock()
	defer count.gate.Unlock()

	if count.size <= 1 {
		return count.line[0].value
	}
	count.prepare(count.isFreq)
	sum := 0.0
	weight := 0.0
	for i := 1; i < count.size; i++ {
		impo := float64(count.size - i)
		elm := &count.line[i]
		sum += elm.value * impo
		weight += impo
	}
	if weight > 0 {
		return sum / weight
	} else {
		return sum
	}
}

func (count *Count) prepare(isFreq bool) {
	count.isFreq = isFreq
	now := time.Now().Unix()
	if count.size == 0 {
		count.line[0] = countElm{}
		count.at = now
		count.size = 1
	}
	delta := now - count.at
	if delta < 0 {
		return
	}
	if delta > 0 {
		for to := countSize - 1; to >= 0; to-- {
			from := to - int(delta)
			if from >= 0 {
				count.line[to] = count.line[from]
			} else if isFreq {
				count.line[to] = countElm{}
			} else {
				count.line[to] = countElm{value: count.lastValue}
			}
		}
		count.at = now
		count.size += int(delta)
		if count.size > countSize {
			count.size = countSize
		}
	}
}
