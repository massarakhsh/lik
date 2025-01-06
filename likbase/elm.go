package likbase

import (
	"time"

	"github.com/massarakhsh/lik"
)

type ItElm struct {
	Table    *ItTable
	Id       lik.IDB
	Info     lik.Seter
	IsModify bool
}

var (
	FieldElms = []DBField{
		DBField{"id", "LP"},
		DBField{"info", "T"},
	}
)

func (elm *ItElm) OnModify() {
	elm.IsModify = true
	elm.Table.OnModify()
	elm.Table.JB.StreemAdd(elm)
}

func (elm *ItElm) OnModifyWait() bool {
	elm.OnModify()
	return elm.Wait()
}

func (elm *ItElm) Wait() bool {
	for nw := 0; nw < 60; nw++ {
		if !elm.IsModify {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}

func (elm *ItElm) Delete() bool {
	return elm.Table.DeleteElm(elm.Id)
}

func (elm *ItElm) ToMap() lik.Seter {
	item := lik.BuildSet()
	if elm.Id > 0 {
		item.SetValue("id", elm.Id)
	}
	elm.ClearInfo()
	if elm.Info != nil {
		item.SetValue("info", elm.Info.Serialize(lik.ITF_JSON))
	} else {
		item.SetValue("info", "{}")
	}
	return item
}

func (elm *ItElm) ClearInfo() {
	elm.clearItemSet(elm.Info)
}

func (elm *ItElm) clearItemSet(item lik.Seter) bool {
	if item == nil {
		return false
	}
	it := item.Self()
	if it == nil {
		return false
	}
	for ns := len(it.Val) - 1; ns >= 0; ns-- {
		delit := false
		if val := it.Val[ns].Val; val == nil {
			delit = true
		} else if val.IsInt() {
			delit = val.ToInt() == 0
		} else if val.IsFloat() {
			delit = val.ToFloat() == 0
		} else if val.IsBool() {
			delit = !val.ToBool()
		} else if val.IsString() {
			delit = val.ToString() == ""
		} else if set := val.ToSet(); set != nil {
			delit = !elm.clearItemSet(set)
		} else if list := val.ToList(); list != nil {
			delit = !elm.clearItemList(list)
		} else {
			delit = true
		}
		if delit {
			item.DelPos(ns)
		}
	}
	return true
}

func (elm *ItElm) clearItemList(item lik.Lister) bool {
	if item == nil {
		return false
	}
	it := item.Self()
	if it == nil {
		return false
	}
	for ne := len(it.Val) - 1; ne >= 0; ne-- {
		val := it.Val[ne]
		delit := false
		if val == nil {
			delit = true
		} else if set := val.ToSet(); set != nil {
			delit = !elm.clearItemSet(set)
		} else if list := val.ToList(); list != nil {
			delit = !elm.clearItemList(list)
		}
		if delit {
			list := []lik.Itemer{}
			for n := 0; n < len(it.Val); n++ {
				if n != ne {
					list = append(list, val)
				}
			}
			it.Val = list
		}
	}
	return true
}

func (elm *ItElm) GetBool(path string) bool {
	value := false
	if elm != nil {
		if item := elm.GetItem(path); item != nil {
			value = item.ToBool()
		}
	}
	return value
}

func (elm *ItElm) GetString(path string) string {
	value := ""
	if item := elm.GetItem(path); item != nil {
		value = item.ToString()
	}
	return value
}

func (elm *ItElm) GetInt(path string) int {
	value := 0
	if elm != nil {
		if item := elm.GetItem(path); item != nil {
			value = int(item.ToInt())
		}
	}
	return value
}

func (elm *ItElm) GetIDB(path string) lik.IDB {
	return lik.IDB(elm.GetInt(path))
}

func (elm *ItElm) GetFloat(path string) float64 {
	value := 0.0
	if item := elm.GetItem(path); item != nil {
		value = item.ToFloat()
	}
	return value
}

func (elm *ItElm) GetList(path string) lik.Lister {
	var value lik.Lister
	if item := elm.GetItem(path); item != nil {
		value = item.ToList()
	}
	return value
}

func (elm *ItElm) GetSet(path string) lik.Seter {
	var value lik.Seter
	if item := elm.GetItem(path); item != nil {
		value = item.ToSet()
	}
	return value
}

func (elm *ItElm) GetItem(path string) lik.Itemer {
	var item lik.Itemer
	if elm.Info != nil {
		item = elm.Info.GetItem(path)
	}
	return item
}

func (elm *ItElm) SetValue(value interface{}, path string) bool {
	modify := false
	if elm.Info == nil {
		elm.Info = lik.BuildSet()
		elm.OnModify()
	}
	if elm.Info.SetValue(path, value) {
		elm.OnModify()
	}
	return modify
}
