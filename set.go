package lik

import (
	"strings"
)

//	Интерфейс динамических структур
type Seter interface {
	Itemer
	Count() int
	Seek(key string) int
	IsItem(path string) bool
	GetItem(path string) Itemer
	GetBool(path string) bool
	GetInt(path string) int
	GetFloat(path string) float64
	GetString(path string) string
	GetList(path string) Lister
	GetSet(path string) Seter
	GetIDB(path string) IDB
	DelItem(path string) bool
	SetItem(val interface{}, path string) bool
	DelPos(pos int) bool
	ToJson() string
	Values() []SetElm
	Keys() []string
	Self() *DItemSet
}

func BuildSet(vals ...interface{}) Seter {
	itemap := &DItemSet{}
	for nv := 0; nv < len(vals); nv++ {
		vk := vals[nv]
		switch key := vk.(type) {
		case string:
			if match := RegExParse(key, "^(.+?)=(.*)"); match != nil {
				key = match[1]
				val := match[2]
				val = strings.Trim(val, "'\"")
				if ival, ok := StrToIntIf(val); ok {
					itemap.SetItem(ival, key)
				} else if fval, ok := StrToFloatIf(val); ok {
					itemap.SetItem(fval, key)
				} else if val == "true" {
					itemap.SetItem(true, key)
				} else if val == "false" {
					itemap.SetItem(false, key)
				} else {
					itemap.SetItem(val, key)
				}
			} else if len(key) > 0 && nv+1 < len(vals) {
				nv++
				itemap.SetItem(vals[nv], key)
			}
		}
	}
	return itemap
}

func (it *DItemSet) clone() Itemer {
	cpy := BuildSet()
	for _, set := range it.Val {
		if val := set.Val; val != nil {
			cpy.SetItem(val.Clone(), set.Key)
		}
	}
	return cpy
}

func (it *DItemSet) serialize() string {
	var text = "{"
	for n, set := range it.Values() {
		if set.Key != "" && set.Val != nil {
			if n > 0 {
				text += ","
			}
			text += StrToQuotes(set.Key) + ":"
			text += set.Val.Serialize()
		}
	}
	text += "}"
	return text
}

func (it *DItemSet) format(prefix string) string {
	var text = "{"
	if it.Count() > 0 {
		for n, set := range it.Values() {
			if set.Key != "" && set.Val != nil {
				if n > 0 {
					text += ","
				}
				text += "\n" + prefix + "  "
				text += StrToQuotes(set.Key) + ":"
				text += set.Val.Format(prefix + "  ")
			}
		}
		text += "\n" + prefix
	}
	text += "}"
	return text
}

func (it *DItemSet) ToJson() string {
	return it.Format("")
}

func (it *DItemSet) Count() int {
	return len(it.Val)
}

func (it *DItemSet) Seek(key string) int {
	for ns := 0; ns < len(it.Val); ns++ {
		if it.Val[ns].Key == key {
			return ns
		}
	}
	return -1
}

func (it *DItemSet) DelPos(pos int) bool {
	if pos < 0 || pos >= len(it.Val) {
		return false
	}
	vals := []SetElm{}
	for n := 0; n < len(it.Val); n++ {
		if n != pos {
			vals = append(vals, it.Val[n])
		}
	}
	it.Val = vals
	return true
}

func (it *DItemSet) IsItem(path string) bool {
	return it.GetItem(path) != nil
}

func (it *DItemSet) GetItem(path string) Itemer {
	var val Itemer
	name, ext := GetFirstExt(path)
	if name == "" && ext == "" {
		val = it
	} else if name == "" {
		val = it.GetItem(ext)
	} else if ns := it.Seek(name); ns >= 0 {
		val = getInfoItem(it.Val[ns].Val, ext)
	}
	return val
}

func (it *DItemSet) GetBool(path string) bool {
	if item := it.GetItem(path); item != nil {
		return item.ToBool()
	}
	return false
}

func (it *DItemSet) GetInt(path string) int {
	if item := it.GetItem(path); item != nil {
		return item.ToInt()
	}
	return 0
}

func (it *DItemSet) GetIDB(path string) IDB {
	return IDB(it.GetInt(path))
}

func (it *DItemSet) GetFloat(path string) float64 {
	if item := it.GetItem(path); item != nil {
		return item.ToFloat()
	}
	return 0
}

func (it *DItemSet) GetString(path string) string {
	if item := it.GetItem(path); item != nil {
		return item.ToString()
	}
	return ""
}

func (it *DItemSet) GetList(path string) Lister {
	if item := it.GetItem(path); item != nil {
		return item.ToList()
	}
	return nil
}

func (it *DItemSet) GetSet(path string) Seter {
	if item := it.GetItem(path); item != nil {
		return item.ToSet()
	}
	return nil
}

func (it *DItemSet) DelItem(path string) bool {
	modify := false
	name, ext := GetFirstExt(path)
	if name == "" && ext == "" {
	} else if name == "" {
		if it.DelItem(ext) {
			modify = true
		}
	} else if ext == "" {
		if ns := it.Seek(name); ns >= 0 {
			if it.DelPos(ns) {
				modify = true
			}
		}
	} else if elm := it.GetSet(name); elm != nil {
		if elm.DelItem(ext) {
			modify = true
		}
	}
	return modify
}

func (it *DItemSet) SetItem(val interface{}, path string) bool {
	modify := false
	if val == nil {
		if it.DelItem(path) {
			modify = true
		}
	} else if name, ext := GetFirstExt(path); name == "" && ext != "" {
		if it.SetItem(val, ext) {
			modify = true
		}
	} else if name != "" {
		ns := it.Seek(name)
		if ns < 0 {
			ns = len(it.Val)
			it.Val = append(it.Val, SetElm{name, nil})
			modify = true
		}
		set := &it.Val[ns]
		if ext == "" {
			if valnew := BuildItem(val); valnew != nil {
				if valnew.IsSet() || valnew.IsList() {
					modify = true
				} else if valold := set.Val; valold != nil {
					if valnew.ToString() != valold.ToString() {
						modify = true
					}
				} else {
					modify = true
				}
				set.Val = valnew
			}
		} else if set.Val != nil {
			if setInfoItem(val, set.Val, ext) {
				modify = true
			}
		} else if RegExCompare(ext, "^(\\d+)") {
			modify = true
			set.Val = BuildList()
			setInfoItem(val, set.Val, ext)
		} else {
			modify = true
			set.Val = BuildSet()
			setInfoItem(val, set.Val, ext)
		}
	}
	return modify
}

func (it *DItemSet) Values() []SetElm {
	return it.Val
}

func (it *DItemSet) Self() *DItemSet {
	return it
}

func (it *DItemSet) Keys() []string {
	keys := []string{}
	for ns := 0; ns < len(it.Val); ns++ {
		keys = append(keys, it.Val[ns].Key)
	}
	return keys
}
