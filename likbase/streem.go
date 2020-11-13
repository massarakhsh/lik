package likbase

import (
	"time"
)

type streemPot struct {
	At         time.Time
	Pred, Next *streemPot
	Elm        *ItElm
}

func StreemGo(jb *JsonBase) {
	jb.streemStop = false
	jb.streemFirst = nil
	jb.streemLast = nil
	jb.streemMap = make(map[string]*streemPot)
	go jb.streemGo()
}

func (it *JsonBase) StreemStop() {
	it.streemStop = true
}

func (it *JsonBase) StreemAdd(elm *ItElm) {
	it.streemSync.Lock()
	if pot := it.streemMap[SignId(elm.Table.Part, elm.Id)]; pot == nil {
		pot = &streemPot{Elm: elm}
		it.streemInsert(pot)
	}
	it.streemSync.Unlock()
}

func (it *JsonBase) streemGo() {
	for {
		if pot := it.streemFirst; pot != nil {
			if pot.Elm.IsModify {
				pot.Elm.Table.SaveElm(pot.Elm)
			}
			it.streemSync.Lock()
			it.streemExtract(pot)
			if pot.Elm.IsModify {
				it.streemInsert(pot)
			}
			it.streemSync.Unlock()
		} else if !it.streemStop {
			time.Sleep(time.Millisecond * 100)
		} else {
			break
		}
	}
}

func (it *JsonBase) streemExtract(pot *streemPot) {
	pred := pot.Pred
	next := pot.Next
	if pred != nil {
		pred.Next = next
		pot.Pred = nil
	} else {
		it.streemFirst = next
	}
	if next != nil {
		next.Pred = pred
		pot.Next = nil
	} else {
		it.streemLast = pred
	}
	delete(it.streemMap, SignId(pot.Elm.Table.Part, pot.Elm.Id))
}

func (it *JsonBase) streemInsert(pot *streemPot) {
	last := it.streemLast
	pot.Pred = last
	if last != nil {
		last.Next = pot
	} else {
		it.streemFirst = pot
	}
	pot.Next = nil
	it.streemLast = pot
	it.streemMap[SignId(pot.Elm.Table.Part, pot.Elm.Id)] = pot
	pot.At = time.Now()
}
