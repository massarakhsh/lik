package metric

func (value *MetricValue) set(proto ProtoMetric, fval float64) {
	value.proto = proto
	value.lastValue = fval

	level := 0
	duraLevel := duraStart
	now := NowMS() / duraStart * duraStart
	addElm := lineElm{start: now, duration: duraStart, count: 1, weight: fval}
	for {
		if level >= len(value.lineLevels) {
			value.lineLevels = append(value.lineLevels, lineLevel{})
		}
		toLevel := &value.lineLevels[level]
		pos := toLevel.pos
		loc := level*duraSize + pos
		if loc >= len(value.listValues) {
			toLevel.size++
			value.listValues = append(value.listValues, addElm)
			break
		}
		elm := &value.listValues[loc]
		if addElm.start/duraLevel == elm.start/duraLevel {
			elm.count += addElm.count
			elm.weight += addElm.weight
			if level > 0 {
				elm.duration += addElm.duration
			}
			break
		}
		pos++
		if pos >= duraSize {
			pos = 0
		}
		toLevel.pos = pos
		loc = level*duraSize + pos
		if toLevel.size >= duraSize {
			elm = &value.listValues[loc]
			pushElm := *elm
			*elm = addElm
			addElm = pushElm
			level++
			duraLevel *= duraFactor
		}
	}
}
