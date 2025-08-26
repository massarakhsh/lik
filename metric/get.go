package metric

import (
	"time"
)

func (value *MetricValue) get() float64 {
	var duration time.Duration = 0
	count := int64(0)
	weight := 0.0
	if len(value.lineLevels) > 0 {
		now := time.Now()
		toLevel := &value.lineLevels[0]
		for n := 0; n < toLevel.size && duration < maxCalcule; n++ {
			pos := (toLevel.pos - n + duraSize) % duraSize
			elm := &value.listValues[pos]
			count += elm.count
			weight += elm.weight
			duration = now.Sub(elm.at)
		}
	}
	if value.proto == protoValue && count > 0 {
		return weight / float64(count)
	} else if value.proto == protoFreq && duration > 0 {
		return weight / duration.Seconds()
	}
	return 0
}
