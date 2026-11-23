package metric

import (
	"strings"
	"sync"
	"time"

	"github.com/massarakhsh/lik"
)

type ProtoMetric int

const protoNo ProtoMetric = 0
const protoValue ProtoMetric = 1
const protoFreq ProtoMetric = 2

type MetricValue struct {
	gate      sync.RWMutex
	proto     ProtoMetric
	lastValue float64
	lastTotal float64

	lineLevels []lineLevel
	listValues []lineElm

	elms []metricElm
}

type metricElm struct {
	name string
	elm  *MetricValue
}

const maxCalcule MS = 16000
const duraStart MS = 4000
const duraSize int = 16
const duraFactor = 2

type lineLevel struct {
	size int
	pos  int
}

type lineElm struct {
	begin  MS
	count  int64
	weight float64
}

type MS int64

func NowMS() MS {
	return MS(time.Now().UnixMilli())
}

func (value *MetricValue) GetLast(name string) float64 {
	return value.GetPathLast(value.nameToPath(name)...)
}

func (value *MetricValue) GetPathLast(path ...string) float64 {
	value.gate.RLock()
	defer value.gate.RUnlock()

	if len(path) == 0 {
		return value.lastValue
	} else if elm, downPath := value.seekMetric(false, path); elm != nil {
		return elm.GetPathLast(downPath...)
	} else {
		return 0
	}
}

func (value *MetricValue) GetTotal(name string) float64 {
	return value.GetPathTotal(value.nameToPath(name)...)
}

func (value *MetricValue) GetPathTotal(path ...string) float64 {
	value.gate.RLock()
	defer value.gate.RUnlock()

	if len(path) == 0 {
		return value.lastTotal
	} else if elm, downPath := value.seekMetric(false, path); elm != nil {
		return elm.GetPathTotal(downPath...)
	} else {
		return 0
	}
}

func (value *MetricValue) SetValueInt(name string, ival int64) {
	value.SetPathInt(ival, value.nameToPath(name)...)
}

func (value *MetricValue) SetPathInt(ival int64, path ...string) {
	value.gate.Lock()
	defer value.gate.Unlock()

	value.setValue(protoValue, float64(ival))
	if elm, downPath := value.seekMetric(true, path); elm != nil {
		elm.SetPathInt(ival, downPath...)
	}
}

func (value *MetricValue) SetValueFloat(name string, fval float64) {
	value.SetPathFloat(fval, value.nameToPath(name)...)
}

func (value *MetricValue) SetPathFloat(fval float64, path ...string) {
	value.gate.Lock()
	defer value.gate.Unlock()

	value.setValue(protoValue, fval)
	if elm, downPath := value.seekMetric(true, path); elm != nil {
		elm.SetPathFloat(fval, downPath...)
	}
}

func (value *MetricValue) Inc(name string) {
	value.IncPath(value.nameToPath(name)...)
}

func (value *MetricValue) IncPath(path ...string) {
	value.gate.Lock()
	defer value.gate.Unlock()

	value.setValue(protoFreq, 1)
	if elm, downPath := value.seekMetric(true, path); elm != nil {
		elm.IncPath(downPath...)
	}
}

func (value *MetricValue) Add(name string, fval float64) {
	value.AddPath(fval, value.nameToPath(name)...)
}

func (value *MetricValue) AddPath(fval float64, path ...string) {
	value.gate.Lock()
	defer value.gate.Unlock()

	value.setValue(protoFreq, fval)
	if elm, downPath := value.seekMetric(true, path); elm != nil {
		elm.AddPath(fval, downPath...)
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
	} else if elm, downPath := value.seekMetric(false, path); elm != nil {
		return elm.GetPath(downPath...)
	} else {
		return 0
	}
}

func (value *MetricValue) GetHistoryPath(at time.Time, duration time.Duration, path ...string) []MetricAt {
	value.gate.RLock()
	defer value.gate.RUnlock()

	if len(path) == 0 {
		return value.getHistory(at, duration)
	} else if elm, downPath := value.seekMetric(false, path); elm != nil {
		return elm.GetHistoryPath(at, duration, downPath...)
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
	} else if elm, downPath := value.seekMetric(false, path); elm != nil {
		return elm.GetCollectPath(downPath...)
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
	} else if elm, downPath := value.seekMetric(false, path); elm != nil {
		return elm.SeekMetric(downPath...)
	} else {
		return nil
	}
}

func (value *MetricValue) seekMetric(create bool, path []string) (*MetricValue, []string) {
	if len(path) == 0 {
		return nil, nil
	}
	var downPath []string
	name := path[0]
	if strings.Contains(name, "/") {
		fullPath := lik.PathToNames(name)
		name = fullPath[0]
		downPath = fullPath[1:]
	}
	downPath = append(downPath, path[1:]...)

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

	return elm, downPath
}
