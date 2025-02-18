package metric

import (
	"time"
)

func (it *MetricValue) getList(atTime time.Time, stepTime time.Duration, need int) []float64 {
	at := atTime.UnixMilli()
	step := stepTime.Milliseconds()
	var values []float64
	count := int64(0)
	weight := 0.0
	var duration int64
	seria := 0
	index := 0
	started := false
	for need > 0 && seria < it.countSeries {
		if index >= it.lenSeries[seria] {
			seria++
			index = 0
			continue
		}
		pos := (it.posSeries[seria] - index + maxElms) % maxElms
		elm := &it.listValues[seria*maxElms+pos]
		if started || at >= elm.at {
			count += elm.count
			weight += elm.weight
			duration += elm.duration
		}
		if at < elm.at {
			index++
			continue
		}
		if started {
			value := 0.0
			if it.proto == protoValue && count > 0 {
				value = weight / float64(count)
			} else if it.proto == protoFreq && duration > 0 {
				value = weight / (float64(duration) / 1000)
			}
			values = append(values, value)
			need--
		} else {
			started = true
		}
		count = 0
		weight = 0.0
		duration = 0
		at -= step
		if at < elm.at {
			index++
		}
	}
	return values
}
