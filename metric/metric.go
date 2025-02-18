package metric

import (
	"strings"
	"sync"
	"time"

	"github.com/massarakhsh/lik"
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

func (it *MetricValue) GetLast(name string) float64 {
	it.gate.RLock()
	defer it.gate.RUnlock()

	if it.isIt(name) {
		return it.lastValue
	} else if elm, next := it.seekMetric(name, false); elm != nil {
		return elm.GetLast(next)
	} else {
		return 0
	}
}

func (it *MetricValue) SetLast(name string, value float64) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.lastValue = value
	if elm, next := it.seekMetric(name, true); elm != nil {
		elm.SetLast(next, value)
	}
}

func (it *MetricValue) SetValueInt(name string, value int64) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.set(float64(value), protoValue)
	if elm, next := it.seekMetric(name, true); elm != nil {
		elm.SetValueInt(next, value)
	}
}

func (it *MetricValue) SetValueFloat(name string, value float64) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.set(value, protoValue)
	if elm, next := it.seekMetric(name, true); elm != nil {
		elm.SetValueFloat(next, value)
	}
}

func (it *MetricValue) Inc(name string) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.set(1, protoFreq)
	if elm, next := it.seekMetric(name, true); elm != nil {
		elm.Inc(next)
	}
}

func (it *MetricValue) Add(name string, value int64) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.set(float64(value), protoFreq)
	if elm, next := it.seekMetric(name, true); elm != nil {
		elm.Add(next, value)
	}
}

func (it *MetricValue) GetValue(name string) float64 {
	it.gate.RLock()
	defer it.gate.RUnlock()

	if it.isIt(name) {
		return it.get()
	} else if elm, next := it.seekMetric(name, false); elm != nil {
		return elm.GetValue(next)
	} else {
		return 0
	}
}

func (it *MetricValue) GetListValues(name string, at time.Time, step time.Duration, need int) []float64 {
	it.gate.RLock()
	defer it.gate.RUnlock()

	if it.isIt(name) {
		return it.getList(at, step, need)
	} else if elm, next := it.seekMetric(name, false); elm != nil {
		return elm.GetListValues(next, at, step, need)
	} else {
		return nil
	}
}

func (it *MetricValue) GetCollect(name string) []string {
	it.gate.RLock()
	defer it.gate.RUnlock()

	if it.isIt(name) {
		var collect []string
		for _, elm := range it.elms {
			collect = append(collect, elm.name)
		}
		return collect
	} else if elm, next := it.seekMetric(name, false); elm != nil {
		return elm.GetCollect(next)
	} else {
		return nil
	}
}

func (it *MetricValue) isIt(path string) bool {
	return lik.RegExCompare(path, "^/*$")
}

func (it *MetricValue) seekMetric(path string, create bool) (*MetricValue, string) {
	name, next := strings.TrimPrefix(path, "/"), ""
	if name == "" {
		return nil, ""
	}
	if match := lik.RegExParse(name, "([^/]*)/(.*)"); match != nil {
		name = match[1]
		next = match[2]
	}
	if name == "" {
		return it.seekMetric(next, create)
	}

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

	return elm, next
}
