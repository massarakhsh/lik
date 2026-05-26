package metric

import (
	"sync"
	"time"
)

const countSize = 10

type Count struct {
	gate      sync.Mutex
	frequency bool
	at        int64
	size      int
	line      [countSize]countElm
	lastValue float64
}

type countElm struct {
	value  float64
	factor float64
}

func (count *Count) Inc() {
	count.Add(1)
}

func (count *Count) Add(val float64) {
	count.SetMod(val, true)
}

func (count *Count) Set(val float64) {
	count.SetMod(val, false)
	count.lastValue = val
}

func (count *Count) SetMod(val float64, frequency bool) {
	count.gate.Lock()
	defer count.gate.Unlock()

	count.prepare(frequency)
	count.line[0].value += val
	count.line[0].factor += 1
}

func (count *Count) Get() float64 {
	count.gate.Lock()
	defer count.gate.Unlock()

	if count.size <= 1 {
		elm := count.line[0]
		if elm.factor > 0 && !count.frequency {
			return float64(elm.value) / float64(elm.factor)
		} else {
			return float64(elm.value)
		}
	}
	count.prepare(count.frequency)
	sum := float64(0)
	factor := 0.0
	for i := 1; i < count.size; i++ {
		elm := count.line[i]
		sum += elm.value
		if !count.frequency {
			factor += float64(count.line[i].factor)
		} else {
			factor += 1
		}
	}
	if factor > 0 {
		return sum / factor
	} else {
		return sum
	}
}

func (count *Count) prepare(frequency bool) {
	count.frequency = frequency
	now := time.Now().Unix()
	if count.size == 0 {
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
			} else if frequency {
				count.line[to] = countElm{}
			} else {
				count.line[to] = countElm{value: count.lastValue, factor: 1}
			}
		}
		count.at = now
		count.size += int(delta)
		if count.size > countSize {
			count.size = countSize
		}
	}
}
