package metric

import (
	"time"
)

type aTV struct {
	at     MS
	value  float64
	factor int
}

func (value *MetricValue) listGet(atTime time.Time, stepTime time.Duration, need int) []float64 {
	var values []float64
	level := 0
	index := 0
	at := MS(atTime.UnixMilli())
	step := MS(stepTime.Milliseconds())

	var left aTV
	var right aTV
	used := false

	for need > 0 && level < len(value.lineLevels) {
		if right.factor == 0 || at < left.at {
			toLevel := &value.lineLevels[level]
			if index >= toLevel.size {
				level++
				index = 0
				continue
			}
			pos := (toLevel.pos - index + duraSize) % duraSize
			toElm := &value.listValues[level*duraSize+pos]
			index++
			if left.factor == 0 {
			} else if right.factor == 0 {
				right = left
			} else if used {
				right = left
				used = false
			} else {
				summ := right.value*float64(right.factor) + left.value
				right.factor++
				right.value = summ / float64(right.factor)
			}
			left = aTV{at: toElm.start, factor: 1}
			left.value = value.calculeValue(toElm)
		} else {
			fVal := 0.0
			if at < right.at {
				fVal = (left.value*float64(right.at-at) + right.value*float64(at-left.at)) / float64(right.at-left.at)
				values = append(values, fVal)
				used = true
			}
			at -= step
			need--
		}
	}
	return values
}

func (value *MetricValue) calculeValue(elm *lineElm) float64 {
	if value.proto == protoValue && elm.count > 0 {
		return elm.weight / float64(elm.count)
	} else if value.proto == protoFreq && elm.duration > 0 {
		return elm.weight / float64(elm.duration) * 1000
	} else {
		return 0
	}
}
