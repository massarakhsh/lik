package metric

import (
	"time"
)

func (it *MetricValue) getList(atTime time.Time, stepTime time.Duration, need int) []float64 {
	at := atTime.UnixMilli()
	step := stepTime.Milliseconds()
	var values []float64
	sumAt := int64(0)
	sumCount := int64(0)
	sumWeight := 0.0
	sumValue := 0.0
	nextAt := at
	nextValue := 0.0
	level := 0
	index := 0
	good := 0
	for need > 0 && level < len(it.lineLevels) {
		toLevel := &it.lineLevels[level]
		if index >= toLevel.size {
			level++
			index = 0
			continue
		}
		pos := (toLevel.pos - index + duraSize) % duraSize
		toElm := &it.listValues[level*duraSize+pos]
		if sumAt == 0 || toElm.at < sumAt {
			oldValue := sumValue
			oldAt := sumAt
			sumCount = toElm.count
			sumWeight = toElm.weight
			sumAt = toElm.at
			sumValue = it.calculeValue(sumCount, sumWeight, nextAt-sumAt)
			if oldAt > 0 {
				nextAt = oldAt
				nextValue = oldValue
			} else {
				nextValue = sumValue
			}
		}
		if over := at - sumAt; over >= 0 {
			if good > 1 {
				value := sumValue
				if dura := nextAt - sumAt; dura > 0 {
					value = (sumValue*float64(dura-over) + nextValue*float64(over)) / float64(dura)
				}
				values = append(values, value)
				need--
			} else {
				good++
			}
			at -= step
		}
		if at < sumAt {
			index++
		}
	}
	return values
}

func (it *MetricValue) calculeValue(count int64, weight float64, duration int64) float64 {
	if it.proto == protoValue && count > 0 {
		return weight / float64(count)
	} else if it.proto == protoFreq && duration > 0 {
		return weight / (float64(duration) / 1000)
	} else {
		return 0
	}
}
