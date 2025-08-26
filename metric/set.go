package metric

import "time"

func (value *MetricValue) set(fval float64, proto ProtoMetric) {
	value.lastValue = fval
	value.proto = proto

	level := 0
	duraLev := duraStart
	now := time.Now()
	addElm := lineValue{at: now, count: 1, weight: fval}
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
		if addElm.at.Sub(elm.at) <= duraLev {
			elm.count += addElm.count
			elm.weight += addElm.weight
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
			duraLev *= duraFactor
		}
	}
}
