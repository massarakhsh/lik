package likbase

import (
	"fmt"
	"github.com/massarakhsh/lik"
	"sort"
	"sync"
)

type ItTable struct {
	Part       string
	Title      string
	Sync       sync.Mutex
	JB         *JsonBase
	Elms       map[lik.IDB]*ItElm
	FaseModify int
}

func (it *ItTable) LoadElms() {
	it.JB.ControlTable(it.Part, FieldElms)
	it.Sync.Lock()
	oldelms := it.Elms
	it.Elms = make(map[lik.IDB]*ItElm)
	it.Sync.Unlock()
	if rows, isit := it.JB.QueryRow(fmt.Sprintf("SELECT id FROM `%s` ORDER BY id", it.Part)); isit {
		for rows.Next() {
			var id lik.IDB
			if err := rows.Scan(&id); err != nil {
				fmt.Println(err)
			} else {
				var elm *ItElm
				if oldelms != nil {
					elm, _ = oldelms[id]
				}
				if elm == nil {
					elm = &ItElm{Table: it}
				}
				elm.Id = id
				it.Sync.Lock()
				it.Elms[elm.Id] = elm
				it.Sync.Unlock()
				if bin := it.JB.GetBinary(it.Part, elm.Id, "info"); bin != nil && len(bin) > 0 {
					str := string(bin)
					elm.Info = lik.SetFromRequest(str)
				} else {
					elm.Info = lik.BuildSet()
				}
			}
		}
		if rows.Close() != nil {
		}
	}
}

func (it *ItTable) Purge() {
	it.JB.Execute(fmt.Sprintf("DELETE FROM `%s`", it.Part))
	it.Sync.Lock()
	it.Elms = make(map[lik.IDB]*ItElm)
	it.Sync.Unlock()
}

func (it *ItTable) Drop() {
	it.Purge()
	it.JB.Execute(fmt.Sprintf("DROP TABLE `%s`", it.Part))
	it.LoadElms()
}

func (it *ItTable) OnModify() {
	it.FaseModify++
}

func (it *ItTable) GetFaseModify() int {
	return it.FaseModify
}

func (it *ItTable) CreateElm() *ItElm {
	elm := &ItElm{Table: it}
	elm.Info = lik.BuildSet()
	it.SaveElm(elm)
	if elm.Id <= 0 {
		elm = nil
	}
	return elm
}

func (it *ItTable) RestoreElm(id lik.IDB, info lik.Seter) *ItElm {
	elm := &ItElm{Table: it}
	elm.Id = id
	elm.Info = info
	sets := elm.ToMap()
	elm.Id = it.JB.InsertElm(it.Part, sets)
	it.Sync.Lock()
	it.Elms[elm.Id] = elm
	it.Sync.Unlock()
	return elm
}

func (it *ItTable) SaveElm(elm *ItElm) {
	sets := elm.ToMap()
	if elm.Id > 0 {
		sets.DelItem("id")
		if it.JB.UpdateElm(it.Part, elm.Id, sets) {
			elm.IsModify = false
		}
	} else if elm.Id = it.JB.InsertElm(it.Part, sets); elm.Id > 0 {
		it.Sync.Lock()
		elm.IsModify = false
		it.Elms[elm.Id] = elm
		it.Sync.Unlock()
	}
}

func (it *ItTable) GetElm(id lik.IDB) *ItElm {
	elm, _ := it.Elms[id]
	return elm
}

func (it *ItTable) InsertElm(sets lik.Seter) *ItElm {
	elm := it.CreateElm()
	return it.UpdateElm(elm.Id, sets)
}

func (it *ItTable) UpdateElm(id lik.IDB, sets lik.Seter) *ItElm {
	elm := it.GetElm(id)
	if elm != nil {
		if !it.JB.UpdateElm(it.Part, id, sets) {
			elm = nil
		}
	}
	return elm
}

func (it *ItTable) DeleteElm(id lik.IDB) bool {
	result := false
	if id > 0 {
		it.JB.DeleteElm(it.Part, id)
		it.Sync.Lock()
		delete(it.Elms, id)
		it.Sync.Unlock()
		it.OnModify()
		result = true
	}
	return result
}

func (it *ItTable) GetListElm(dir bool) []*ItElm {
	ids := []int{}
	for id, _ := range it.Elms {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)
	list := []*ItElm{}
	mi := len(ids)
	for ni := 0; ni < mi; ni++ {
		pi := ni
		if !dir {
			pi = mi - 1 - ni
		}
		if elm := it.GetElm(lik.IDB(ids[pi])); elm != nil {
			list = append(list, elm)
		}
	}
	return list
}
