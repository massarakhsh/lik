package metric

import (
	"sync"
	"time"
)

const protoNo = 0
const protoValue = 1
const protoFreq = 2

type MetricValue struct {
	gate      sync.RWMutex
	proto     int
	lastValue float64

	countSeries int
	posSeries   []int
	lenSeries   []int
	listValues  []calcLine

	elms []MetricElm
}

type MetricElm struct {
	name string
	elm  *MetricValue
}

const maxElms int = 100
const duraElm int64 = 1000
const duraFactor = 5
const maxCalcule int64 = 1000 * 10

type calcLine struct {
	at       int64
	duration int64
	count    int64
	weight   float64
}

func (it *MetricValue) GetLast(path ...string) float64 {
	it.gate.RLock()
	defer it.gate.RUnlock()

	if len(path) == 0 {
		return it.lastValue
	} else if elm := it.seekMetric(false, path); elm != nil {
		return elm.GetLast(path[1:]...)
	} else {
		return 0
	}
}

func (it *MetricValue) SetLast(value float64, path ...string) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.lastValue = value
	if elm := it.seekMetric(true, path); elm != nil {
		elm.SetLast(value, path[1:]...)
	}
}

func (it *MetricValue) SetValueInt(value int64, path ...string) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.set(float64(value), protoValue)
	if elm := it.seekMetric(true, path); elm != nil {
		elm.SetValueInt(value, path[1:]...)
	}
}

func (it *MetricValue) SetValueFloat(value float64, path ...string) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.set(value, protoValue)
	if elm := it.seekMetric(true, path); elm != nil {
		elm.SetValueFloat(value, path[1:]...)
	}
}

func (it *MetricValue) Inc(path ...string) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.set(1, protoFreq)
	if elm := it.seekMetric(true, path); elm != nil {
		elm.Inc(path[1:]...)
	}
}

func (it *MetricValue) Add(value float64, path ...string) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.set(value, protoFreq)
	if elm := it.seekMetric(true, path); elm != nil {
		elm.Add(value, path[1:]...)
	}
}

func (it *MetricValue) GetValue(path ...string) float64 {
	it.gate.RLock()
	defer it.gate.RUnlock()

	if len(path) == 0 {
		return it.get()
	} else if elm := it.seekMetric(false, path); elm != nil {
		return elm.GetValue(path[1:]...)
	} else {
		return 0
	}
}

func (it *MetricValue) GetListValues(at time.Time, step time.Duration, need int, path ...string) []float64 {
	it.gate.RLock()
	defer it.gate.RUnlock()

	if len(path) == 0 {
		return it.getList(at, step, need)
	} else if elm := it.seekMetric(false, path); elm != nil {
		return elm.GetListValues(at, step, need, path[1:]...)
	} else {
		return nil
	}
}

func (it *MetricValue) GetCollect(path ...string) []string {
	it.gate.RLock()
	defer it.gate.RUnlock()

	if len(path) == 0 {
		var collect []string
		for _, elm := range it.elms {
			collect = append(collect, elm.name)
		}
		return collect
	} else if elm := it.seekMetric(false, path); elm != nil {
		return elm.GetCollect(path[1:]...)
	} else {
		return nil
	}
}

func (it *MetricValue) seekMetric(create bool, path []string) *MetricValue {
	if len(path) == 0 {
		return nil
	}
	name := path[0]

	var elm *MetricValue
	max := len(it.elms)
	pos := 0
	for pos = 0; pos < max; pos++ {
		if what := it.elms[pos].name; what == name {
			elm = it.elms[pos].elm
			break
		} else if what > name {
			break
		}
	}
	if elm == nil && create {
		elms := make([]MetricElm, max+1)
		for n := 0; n < pos; n++ {
			elms[n] = it.elms[n]
		}
		elm = &MetricValue{}
		elms[pos] = MetricElm{name: name, elm: elm}
		for n := pos; n < max; n++ {
			elms[n+1] = it.elms[n]
		}
		it.elms = elms
	}

	return elm
}
