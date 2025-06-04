package metric

import "time"

func (it *MetricValue) set(value float64, proto ProtoMetric) {
	it.lastValue = value
	it.proto = proto

	level := 0
	duraLev := duraStart
	now := time.Now().UnixMilli() /*/ duraStart * duraStart*/
	addElm := lineValue{at: now, count: 1, weight: value}
	for {
		if level >= len(it.lineLevels) {
			it.lineLevels = append(it.lineLevels, lineLevel{})
		}
		toLevel := &it.lineLevels[level]
		pos := toLevel.pos
		loc := level*duraSize + pos
		if loc >= len(it.listValues) {
			toLevel.size++
			it.listValues = append(it.listValues, addElm)
			break
		}
		elm := &it.listValues[loc]
		if addElm.at/duraLev == elm.at/duraLev {
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
			elm = &it.listValues[loc]
			pushElm := *elm
			*elm = addElm
			addElm = pushElm
			level++
			duraLev *= duraFactor
		}
	}
}
