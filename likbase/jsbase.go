package likbase

import (
	"sync"
)

const (
	DB_Driver = "mysql"
)

type JsonBase struct {
	DBase
	streemSync              sync.Mutex
	streemFirst, streemLast *streemPot
	streemMap               map[string]*streemPot
	streemStop              bool
}

type JsonBaser interface {
	DBaser
	StopDB()
	BuildTable(part string, name string) *ItTable
}

func OpenJsonBase(serv string, name string, user string, pass string) JsonBaser {
	jsonbase := &JsonBase{}
	logon := ""
	if user != "" && pass != "" {
		logon = user + ":" + pass
	}
	addr := "tcp(" + serv + ":3306)"
	if jsonbase.openDBase(DB_Driver, logon, addr, name) {
		StreemGo(jsonbase)
	} else {
		jsonbase = nil
	}
	return jsonbase
}

func (it *JsonBase) StopDB() {
	it.StreemStop()
}

func SortById(list []*ItElm) []*ItElm {
	size := len(list)
	for doze := 1; doze < size; doze *= 2 {
		trans := []*ItElm{}
		for pos := 0; pos < size; pos += doze * 2 {
			posa := pos
			enda := pos + doze
			if enda > size {
				enda = size
			}
			posb := enda
			endb := posb + doze
			if endb > size {
				endb = size
			}
			for posa < enda || posb < endb {
				cmp := 0
				if posa >= enda {
					cmp = 1
				} else if posb >= endb {
					cmp = -1
				} else {
					ida := list[posa].Id
					idb := list[posb].Id
					if ida < idb {
						cmp = 1
					} else if ida > idb {
						cmp = -1
					}
				}
				if cmp <= 0 {
					trans = append(trans, list[posa])
					posa++
				} else {
					trans = append(trans, list[posb])
					posb++
				}
			}
		}
		list = trans
	}
	return list
}

func (it *JsonBase) BuildTable(part string, name string) *ItTable {
	table := &ItTable{JB: it, Part: part, Title: name}
	return table
}
