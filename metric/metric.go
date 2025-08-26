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

const duraStart time.Duration = time.Second
const duraSize int = 64
const duraFactor = 2
const maxCalcule time.Duration = time.Second * 10

type lineLevel struct {
	size int
	pos  int
}

type lineValue struct {
	at     time.Time
	count  int64
	weight float64
}

func (value *MetricValue) GetLast(name string) float64 {
	return value.GetPathLast(value.nameToPath(name)...)
}

func (value *MetricValue) GetPathLast(path ...string) float64 {
	value.gate.RLock()
	defer value.gate.RUnlock()

	if len(path) == 0 {
		return value.lastValue
	} else if elm := value.seekMetric(false, path); elm != nil {
		return elm.GetPathLast(path[1:]...)
	} else {
		return 0
	}
}

func (value *MetricValue) SetLast(name string, fval float64) {
	value.SetPathLast(fval, value.nameToPath(name)...)
}

func (value *MetricValue) SetPathLast(fval float64, path ...string) {
	value.gate.Lock()
	defer value.gate.Unlock()

	value.lastValue = fval
	if elm := value.seekMetric(true, path); elm != nil {
		elm.SetPathLast(fval, path[1:]...)
	}
}

func (value *MetricValue) SetValueInt(name string, ival int64) {
	value.SetPathInt(ival, value.nameToPath(name)...)
}

func (value *MetricValue) SetPathInt(ival int64, path ...string) {
	value.gate.Lock()
	defer value.gate.Unlock()

	value.set(float64(ival), protoValue)
	if elm := value.seekMetric(true, path); elm != nil {
		elm.SetPathInt(ival, path[1:]...)
	}
}

func (value *MetricValue) SetValueFloat(name string, fval float64) {
	value.SetPathFloat(fval, value.nameToPath(name)...)
}

func (value *MetricValue) SetPathFloat(fval float64, path ...string) {
	value.gate.Lock()
	defer value.gate.Unlock()

	value.set(fval, protoValue)
	if elm := value.seekMetric(true, path); elm != nil {
		elm.SetPathFloat(fval, path[1:]...)
	}
}

func (value *MetricValue) Inc(name string) {
	value.IncPath(value.nameToPath(name)...)
}

func (value *MetricValue) IncPath(path ...string) {
	value.gate.Lock()
	defer value.gate.Unlock()

	value.set(1, protoFreq)
	if elm := value.seekMetric(true, path); elm != nil {
		elm.IncPath(path[1:]...)
	}
}

func (value *MetricValue) Add(name string, fval float64) {
	value.AddPath(fval, value.nameToPath(name)...)
}

func (value *MetricValue) AddPath(fval float64, path ...string) {
	value.gate.Lock()
	defer value.gate.Unlock()

	value.set(fval, protoFreq)
	if elm := value.seekMetric(true, path); elm != nil {
		elm.AddPath(fval, path[1:]...)
	}
}

func (value *MetricValue) GetValue(name string) float64 {
	return value.GetPath(value.nameToPath(name)...)
}

func (value *MetricValue) GetPath(path ...string) float64 {
	value.gate.RLock()
	defer value.gate.RUnlock()

	if len(path) == 0 {
		return value.get()
	} else if elm := value.seekMetric(false, path); elm != nil {
		return elm.GetPath(path[1:]...)
	} else {
		return 0
	}
}

func (value *MetricValue) GetListValues(name string, at time.Time, step time.Duration, need int) []float64 {
	return value.GetListPath(at, step, need, value.nameToPath(name)...)
}

func (value *MetricValue) GetListPath(at time.Time, step time.Duration, need int, path ...string) []float64 {
	value.gate.RLock()
	defer value.gate.RUnlock()

	if len(path) == 0 {
		return value.getList(at, step, need)
	} else if elm := value.seekMetric(false, path); elm != nil {
		return elm.GetListPath(at, step, need, path[1:]...)
	} else {
		return nil
	}
}

func (value *MetricValue) GetCollect(name string) []string {
	return value.GetCollectPath(value.nameToPath(name)...)
}

func (value *MetricValue) GetCollectPath(path ...string) []string {
	value.gate.RLock()
	defer value.gate.RUnlock()

	if len(path) == 0 {
		var collect []string
		for _, elm := range value.elms {
			collect = append(collect, elm.name)
		}
		return collect
	} else if elm := value.seekMetric(false, path); elm != nil {
		return elm.GetCollectPath(path[1:]...)
	} else {
		return nil
	}
}

func (value *MetricValue) nameToPath(name string) []string {
	if names := strings.Trim(name, "/"); names != "" {
		return strings.Split(names, "/")
	} else {
		return nil
	}
}

func (value *MetricValue) SeekMetric(path ...string) *MetricValue {
	value.gate.RLock()
	defer value.gate.RUnlock()

	if len(path) == 0 {
		return value
	} else if elm := value.seekMetric(false, path); elm != nil {
		return elm.SeekMetric(path[1:]...)
	} else {
		return nil
	}
}

func (value *MetricValue) seekMetric(create bool, path []string) *MetricValue {
	if len(path) == 0 {
		return nil
	}
	name := path[0]

	var elm *MetricValue
	max := len(value.elms)
	pos := 0
	for pos = 0; pos < max; pos++ {
		if what := value.elms[pos].name; what == name {
			elm = value.elms[pos].elm
			break
		} else if what > name {
			break
		}
	}
	if elm == nil && create {
		elms := make([]metricElm, max+1)
		for n := 0; n < pos; n++ {
			elms[n] = value.elms[n]
		}
		elm = &MetricValue{}
		elms[pos] = metricElm{name: name, elm: elm}
		for n := pos; n < max; n++ {
			elms[n+1] = value.elms[n]
		}
		value.elms = elms
	}

	return elm
}
