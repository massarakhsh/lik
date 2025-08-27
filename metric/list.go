package metric

import (
	"time"
)

type aTV struct {
	at     time.Time
	to     time.Time
	value  float64
	factor int
}

func (value *MetricValue) listGet(atTime time.Time, stepTime time.Duration, need int) []float64 {
	var values []float64
	level := 0
	index := 0

	var left aTV
	var right aTV
	used := false

	for need > 0 && level < len(value.lineLevels) {
		if right.factor == 0 || atTime.Before(left.at) {
			toLevel := &value.lineLevels[level]
			if index >= toLevel.size {
				level++
				index = 0
				continue
			}
			pos := (toLevel.pos - index + duraSize) % duraSize
			toElm := &value.listValues[level*duraSize+pos]
			index++
			var leftTo time.Time
			if left.factor == 0 {
				leftTo = atTime
			} else if right.factor == 0 {
				right = left
				leftTo = right.at
			} else if used {
				right = left
				leftTo = right.at
				used = false
			} else {
				summ := right.value*float64(right.factor) + left.value
				right.factor++
				right.value = summ / float64(right.factor)
				leftTo = right.at
			}
			left = aTV{at: toElm.at, to: leftTo, factor: 1}
			left.value = value.calculeValue(toElm, leftTo)
		} else {
			fVal := 0.0
			if atTime.Before(right.at) {
				fVal = (left.value*right.at.Sub(atTime).Seconds() + right.value*atTime.Sub(left.at).Seconds()) / right.at.Sub(left.at).Seconds()
				values = append(values, fVal)
				used = true
			}
			atTime = atTime.Add(-stepTime)
			need--
		}
	}
	return values
}

func (value *MetricValue) calculeValue(elm *lineValue, to time.Time) float64 {
	if value.proto == protoValue && elm.count > 0 {
		return elm.weight / float64(elm.count)
	} else if value.proto == protoFreq && to.After(elm.at) {
		return elm.weight / to.Sub(elm.at).Seconds()
	} else {
		return 0
	}
}
