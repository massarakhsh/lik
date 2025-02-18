package metric

import (
	"time"
)

func (it *MetricValue) get() float64 {
	var dura int64
	count := int64(0)
	weight := 0.0
	if it.countSeries > 0 {
		now := time.Now().UnixMilli()
		for n := 0; n < it.lenSeries[0]; n++ {
			pos := (it.posSeries[0] - n + maxElms) % maxElms
			elm := &it.listValues[pos]
			if now-elm.at >= maxCalcule {
				break
			}
			count += elm.count
			weight += elm.weight
			dura = now - elm.at
		}
	}
	if it.proto == protoValue && count > 0 {
		return weight / float64(count)
	} else if it.proto == protoFreq && dura > 0 {
		return weight / (float64(dura) / 1000)
	}
	return 0
}
