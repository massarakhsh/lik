package metric

import (
	"time"
)

type MetricAt struct {
	At    time.Time
	Value float64
}

func (value *MetricValue) getHistory(atTime time.Time, duration time.Duration) []MetricAt {
	var values []MetricAt
	level := 0
	index := 0
	to := MS(atTime.UnixMilli())
	from := MS(atTime.Add(-duration).UnixMilli())

	for level < len(value.lineLevels) {
		toLevel := &value.lineLevels[level]
		if index >= toLevel.size {
			level++
			index = 0
			continue
		}
		pos := (toLevel.pos - index + duraSize) % duraSize
		toElm := &value.listValues[level*duraSize+pos]
		index++
		if toElm.begin < from {
			break
		}
		if toElm.begin > to {
			continue
		}
		var val MetricAt
		val.At = time.UnixMilli(int64(toElm.begin))
		val.Value = value.calculeValue(toElm, to)
		values = append(values, val)
		to = toElm.begin
	}

	return values
}

func (value *MetricValue) calculeValue(elm *lineElm, to MS) float64 {
	if value.proto == protoValue && elm.count > 0 {
		return elm.weight / float64(elm.count)
	} else if value.proto == protoFreq && elm.begin < to {
		return elm.weight / float64(to-elm.begin) * 1000
	} else {
		return 0
	}
}
