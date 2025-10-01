package metric

func (value *MetricValue) get() float64 {
	var duration MS
	count := int64(0)
	weight := 0.0
	now := NowMS()
	if len(value.lineLevels) > 0 {
		toLevel := &value.lineLevels[0]
		for n := 0; n < toLevel.size && duration < maxCalcule; n++ {
			pos := (toLevel.pos - n + duraSize) % duraSize
			elm := &value.listValues[pos]
			if n > 0 {
				count += elm.count
				weight += elm.weight
				duration += elm.duration
			} else if value.proto == protoFreq && now-elm.start > duraStart {
				duration += now - elm.start - duraStart
			}
		}
	}
	if value.proto == protoValue && count > 0 {
		return weight / float64(count)
	} else if value.proto == protoFreq && duration > 0 {
		return weight / float64(duration) * 1000
	}
	return 0
}
