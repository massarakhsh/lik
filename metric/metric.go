package metric

import (
	"strings"
	"sync"
	"time"
)

type ProtoMetric int

const protoNo ProtoMetric = 0
const protoValue ProtoMetric = 1
const protoFreq ProtoMetric = 2

type MetricValue struct {
	gate      sync.RWMutex
	proto     ProtoMetric
	lastValue float64

	lineLevels []lineLevel
	listValues []lineValue

	elms []metricElm
}

type metricElm struct {
	name string
	elm  *MetricValue
}

const duraStart int64 = 1024
const duraSize int = 64
const duraFactor = 2
const maxCalcule int64 = 1000 * 10

type lineLevel struct {
	size int
	pos  int
}

type lineValue struct {
	at     int64
	count  int64
	weight float64
}

func (it *MetricValue) GetLast(name string) float64 {
	return it.GetPathLast(it.nameToPath(name)...)
}

func (it *MetricValue) GetPathLast(path ...string) float64 {
	it.gate.RLock()
	defer it.gate.RUnlock()

	if len(path) == 0 {
		return it.lastValue
	} else if elm := it.seekMetric(false, path); elm != nil {
		return elm.GetPathLast(path[1:]...)
	} else {
		return 0
	}
}

func (it *MetricValue) SetLast(name string, value float64) {
	it.SetPathLast(value, it.nameToPath(name)...)
}

func (it *MetricValue) SetPathLast(value float64, path ...string) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.lastValue = value
	if elm := it.seekMetric(true, path); elm != nil {
		elm.SetPathLast(value, path[1:]...)
	}
}

func (it *MetricValue) SetValueInt(name string, value int64) {
	it.SetPathInt(value, it.nameToPath(name)...)
}

func (it *MetricValue) SetPathInt(value int64, path ...string) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.set(float64(value), protoValue)
	if elm := it.seekMetric(true, path); elm != nil {
		elm.SetPathInt(value, path[1:]...)
	}
}

func (it *MetricValue) SetValueFloat(name string, value float64) {
	it.SetPathFloat(value, it.nameToPath(name)...)
}

func (it *MetricValue) SetPathFloat(value float64, path ...string) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.set(value, protoValue)
	if elm := it.seekMetric(true, path); elm != nil {
		elm.SetPathFloat(value, path[1:]...)
	}
}

func (it *MetricValue) Inc(name string) {
	it.IncPath(it.nameToPath(name)...)
}

func (it *MetricValue) IncPath(path ...string) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.set(1, protoFreq)
	if elm := it.seekMetric(true, path); elm != nil {
		elm.IncPath(path[1:]...)
	}
}

func (it *MetricValue) Add(name string, value float64) {
	it.AddPath(value, it.nameToPath(name)...)
}

func (it *MetricValue) AddPath(value float64, path ...string) {
	it.gate.Lock()
	defer it.gate.Unlock()

	it.set(value, protoFreq)
	if elm := it.seekMetric(true, path); elm != nil {
		elm.AddPath(value, path[1:]...)
	}
}

func (it *MetricValue) GetValue(name string) float64 {
	return it.GetPath(it.nameToPath(name)...)
}

func (it *MetricValue) GetPath(path ...string) float64 {
	it.gate.RLock()
	defer it.gate.RUnlock()

	if len(path) == 0 {
		return it.get()
	} else if elm := it.seekMetric(false, path); elm != nil {
		return elm.GetPath(path[1:]...)
	} else {
		return 0
	}
}

func (it *MetricValue) GetListValues(name string, at time.Time, step time.Duration, need int) []float64 {
	return it.GetListPath(at, step, need, it.nameToPath(name)...)
}

func (it *MetricValue) GetListPath(at time.Time, step time.Duration, need int, path ...string) []float64 {
	it.gate.RLock()
	defer it.gate.RUnlock()

	if len(path) == 0 {
		return it.getList(at, step, need)
	} else if elm := it.seekMetric(false, path); elm != nil {
		return elm.GetListPath(at, step, need, path[1:]...)
	} else {
		return nil
	}
}

func (it *MetricValue) GetCollect(name string) []string {
	return it.GetCollectPath(it.nameToPath(name)...)
}

func (it *MetricValue) GetCollectPath(path ...string) []string {
	it.gate.RLock()
	defer it.gate.RUnlock()

	if len(path) == 0 {
		var collect []string
		for _, elm := range it.elms {
			collect = append(collect, elm.name)
		}
		return collect
	} else if elm := it.seekMetric(false, path); elm != nil {
		return elm.GetCollectPath(path[1:]...)
	} else {
		return nil
	}
}

func (it *MetricValue) nameToPath(name string) []string {
	if names := strings.Trim(name, "/"); names != "" {
		return strings.Split(names, "/")
	} else {
		return nil
	}
}

func (it *MetricValue) SeekMetric(path ...string) *MetricValue {
	it.gate.RLock()
	defer it.gate.RUnlock()

	if len(path) == 0 {
		return it
	} else if elm := it.seekMetric(false, path); elm != nil {
		return elm.SeekMetric(path[1:]...)
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
		elms := make([]metricElm, max+1)
		for n := 0; n < pos; n++ {
			elms[n] = it.elms[n]
		}
		elm = &MetricValue{}
		elms[pos] = metricElm{name: name, elm: elm}
		for n := pos; n < max; n++ {
			elms[n+1] = it.elms[n]
		}
		it.elms = elms
	}

	return elm
}
