package lik

import (
	"reflect"
	"sort"
	"strings"
)

// Интерфейс динамических структур
type Seter interface {
	Itemer
	Count() int
	Seek(key string) int
	IsItem(path string) bool
	GetItem(path string) Itemer
	GetBool(path string) bool
	GetInt(path string) int64
	GetFloat(path string) float64
	GetString(path string) string
	GetList(path string) Lister
	GetSet(path string) Seter
	GetIDB(path string) IDB
	DelItem(path string) bool
	SetValue(path string, val interface{}) bool
	SetValues(vals ...interface{})
	AddSet(path string) Seter
	AddList(path string) Lister
	DelPos(pos int) bool
	Merge(set Seter)
	ToJson() string
	Values() []SetElm
	Keys() []string
	SortKeys() []string
	Self() *DItemSet
	SetString(key string, val string)
}

func BuildSet(vals ...interface{}) Seter {
	itemap := &DItemSet{}
	itemap.SetValues(vals...)
	return itemap
}

func BuildStringSet(vals ...string) Seter {
	var opts []interface{}
	for _, val := range vals {
		opts = append(opts, val)
	}
	return BuildSet(opts...)
}

func (it *DItemSet) SetValues(vals ...interface{}) {
	for nv := 0; nv < len(vals); nv++ {
		vk := vals[nv]
		switch key := vk.(type) {
		case string:
			if match := RegExParse(key, "^(.+?)=(.*)"); match != nil {
				key = match[1]
				val := match[2]
				val = strings.Trim(val, "'\"")
				if ival, ok := StrToIntIf(val); ok {
					it.SetValue(key, ival)
				} else if fval, ok := StrToFloatIf(val); ok {
					it.SetValue(key, fval)
				} else if val == "true" {
					it.SetValue(key, true)
				} else if val == "false" {
					it.SetValue(key, false)
				} else {
					it.SetValue(key, val)
				}
			} else if len(key) > 0 && nv+1 < len(vals) {
				nv++
				it.SetValue(key, vals[nv])
			} else if len(key) > 0 {
				it.SetValue(key, "")
			}
		default:
			if item := BuildItem(vk); item != nil && item.IsSet() {
				it.Merge(item.ToSet())
			}
		}
	}
}

func (it *DItemSet) clone() Itemer {
	cpy := BuildSet()
	for _, set := range it.Val {
		if val := set.Val; val != nil {
			cpy.SetValue(set.Key, val.Clone())
		}
	}
	return cpy
}

func (it *DItemSet) serialize(itf int) string {
	var text = "{"
	for _, set := range it.Values() {
		if set.Key != "" && set.Val != nil {
			if len(text) > 1 {
				text += ","
			}
			if itf == ITF_JSON {
				text += StrToQuotes(set.Key) + ":"
			} else {
				text += set.Key + ":"
			}
			text += set.Val.serializeAs(itf)
		}
	}
	text += "}"
	return text
}

func (it *DItemSet) sort_serialize() string {
	var text = "{"
	keys := it.SortKeys()
	for _, key := range keys {
		val := it.GetItem(key)
		if key != "" && val != nil {
			if len(text) > 1 {
				text += ","
			}
			text += StrToQuotes(key) + ":"
			text += val.SortSerialize()
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
	if path == "" || path == "/" {
		val = it
		// } else if name == "" {
		// 	val = it.GetItem(ext)
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

func (it *DItemSet) GetInt(path string) int64 {
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

func (it *DItemSet) SetValue(path string, val interface{}) bool {
	modify := false
	if val == nil {
		if it.DelItem(path) {
			modify = true
		}
	} else if name, ext := GetFirstExt(path); name == "" && ext != "" {
		if it.SetValue(ext, val) {
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

func (it *DItemSet) AddSet(path string) Seter {
	item := BuildSet()
	it.SetValue(path, item)
	return item
}

func (it *DItemSet) AddList(path string) Lister {
	item := BuildList()
	it.SetValue(path, item)
	return item
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

func (it *DItemSet) SortKeys() []string {
	keys := it.Keys()
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func (it *DItemSet) SetString(key string, val string) {
	if val == "true" {
		it.SetValue(key, true)
	} else if val == "false" {
		it.SetValue(key, false)
	} else if num, ok := StrToIntIf(val); ok {
		it.SetValue(key, num)
	} else if num, ok := StrToFloatIf(val); ok {
		it.SetValue(key, num)
	} else {
		it.SetValue(key, val)
	}
}

func (it *DItemSet) Merge(set Seter) {
	if set != nil {
		for _, pair := range set.Values() {
			if val := pair.Val; val.IsSet() {
				if to := it.GetSet(pair.Key); to != nil {
					to.Merge(val.ToSet())
				} else {
					it.SetValue(pair.Key, val)
				}
			} else {
				it.SetValue(pair.Key, val)
			}
		}
	}
}

func (it *DItemSet) setFromReflectStructure(tp reflect.Type, vl reflect.Value) {
	cnt := vl.NumField()
	for f := 0; f < cnt; f++ {
		tpf := tp.Field(f)
		if val := vl.Field(f); val.IsValid() {
			if !val.IsZero() {
				if item := BuildItemReflect(val); item != nil {
					if !tpf.Anonymous {
						nam := ToSnakeCase(tpf.Name)
						it.SetValue(nam, item)
					} else if set := item.ToSet(); set != nil {
						it.Merge(set)
					}
				}
			}
		}
	}
}

func (it *DItemSet) setFromReflectMap(tp reflect.Type, vl reflect.Value) {
	if keys := vl.MapKeys(); keys != nil {
		for _, key := range keys {
			nam := key.String()
			if val := vl.MapIndex(key); val.CanInterface() {
				it.SetValue(nam, val.Interface())
			}
		}
	}
}
