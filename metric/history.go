package metric

import (
	"time"
)

type MetricAt struct {
	At    time.Time
	Value float64
}

func (value *MetricValue) historyGet(atTime time.Time, duration time.Duration) []MetricAt {
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
		if toElm.start+toElm.duration < from {
			break
		}
		if toElm.start > to {
			continue
		}
		var val MetricAt
		val.At = time.UnixMilli(int64(toElm.start + toElm.duration/2))
		val.Value = value.calculeValue(toElm)
		values = append(values, val)
	}

	return values
}
