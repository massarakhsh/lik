package metric

import "time"

func (it *MetricValue) set(value float64, proto int) {
	it.lastValue = value
	it.proto = proto

	duraMax := duraElm
	addElm := calcLine{at: time.Now().UnixMilli(), count: 1, weight: value}
	seria := 0
	for {
		if seria >= it.countSeries {
			it.countSeries = seria + 1
			it.posSeries = append(it.posSeries, 0)
			it.lenSeries = append(it.lenSeries, 0)
		}
		pos := it.posSeries[seria]
		if pos >= it.lenSeries[seria] {
			it.lenSeries[seria]++
			it.listValues = append(it.listValues, addElm)
			break
		}
		elm := &it.listValues[seria*maxElms+pos]
		if addElm.at/duraMax == elm.at/duraMax {
			elm.count += addElm.count
			elm.weight += addElm.weight
			elm.duration = addElm.at - elm.at + addElm.duration
			break
		}
		end := addElm.at + addElm.duration
		addElm.at = addElm.at / duraMax * duraMax
		elm.duration = addElm.at - elm.at
		addElm.duration = end - addElm.at
		pos++
		if pos >= maxElms {
			pos = 0
		}
		it.posSeries[seria] = pos
		if it.lenSeries[seria] >= maxElms {
			elm = &it.listValues[seria*maxElms+pos]
			pushElm := *elm
			*elm = addElm
			addElm = pushElm
			duraMax *= duraFactor
			seria++
		}
	}
}
