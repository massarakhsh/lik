package lik

import (
	"fmt"
	"reflect"
	"strings"
)

func StrToQuotes(str string) string {
	val := strings.Replace(str, "\\", "\\\\", -1)
	val = strings.Replace(val, "\"", "\\\"", -1)
	val = strings.Replace(val, "\n", "\\n", -1)
	val = strings.Replace(val, "\r", "\\r", -1)
	return "\"" + val + "\""
}

func SetToType[T any](set Seter) T {
	return ItemToType[T](set)
}

func ListToType[T any](list Lister) T {
	return ItemToType[T](list)
}

func ItemToType[T any](item Itemer) T {
	var res T
	result := &res
	tp := reflect.TypeOf(result)
	vl := reflect.ValueOf(result)
	item_to_reflect(item, tp, vl)
	return res
}

func item_to_reflect(item Itemer, tp reflect.Type, vl reflect.Value) bool {
	if item == nil || !vl.IsValid() {
		return false
	}
	result := true

	if stp := tp.Kind().String(); stp == "ptr" {
		return item_to_reflect(item, tp.Elem(), vl.Elem())
	} else if stp == "string" {
		if val := reflect.ValueOf(item.ToString()); val.CanConvert(tp) {
			if vl.CanSet() {
				vl.Set(val.Convert(tp))
			}
		}
	} else if stp == "bool" {
		if val := reflect.ValueOf(item.ToBool()); val.CanConvert(tp) {
			if vl.CanSet() {
				vl.Set(val.Convert(tp))
			}
		}
	} else if stp == "int" || stp == "uint" || stp == "int32" || stp == "uint32" || stp == "int64" || stp == "uint64" {
		if val := reflect.ValueOf(item.ToInt()); val.CanConvert(tp) {
			if vl.CanSet() {
				vl.Set(val.Convert(tp))
			}
		}
	} else if stp == "struct" {
		result = set_to_struct(item.ToSet(), tp, vl)
	} else if stp == "map" {
		result = set_to_map(item.ToSet(), tp, vl)
	} else if stp == "slice" {
		result = list_to_slice(item.ToList(), tp, vl)
	} else {
		if val := reflect.ValueOf(item); val.CanConvert(tp) {
			if vl.CanSet() {
				vl.Set(val.Convert(tp))
			}
		}
	}
	return result
}

func set_to_struct(set Seter, tp reflect.Type, vl reflect.Value) bool {
	if set == nil || !vl.IsValid() {
		return false
	}
	result := true
	cnt := vl.NumField()
	for f := 0; f < cnt; f++ {
		tpf := tp.Field(f)
		var item Itemer
		if tpf.Anonymous {
			item = set
		} else if item = set.GetItem(tpf.Name); item == nil {
			item = set.GetItem(ToSnakeCase(tpf.Name))
		}
		if item != nil {
			if !item_to_reflect(item, tpf.Type, vl.Field(f)) {
				result = false
			}
		}
	}
	return result
}

func set_to_map(set Seter, tp reflect.Type, vl reflect.Value) bool {
	if set == nil || !vl.IsValid() {
		return false
	}
	result := true
	if tpkey := tp.Key(); tpkey.Kind().String() == "string" {
		tpelm := tp.Elem()
		tpmap := reflect.MapOf(tpkey, tpelm)
		vmap := reflect.MakeMap(tpmap)
		for _, pair := range set.Values() {
			valelm := reflect.New(tpelm)
			item_to_reflect(pair.Val, tpelm, valelm.Elem())
			vmap.SetMapIndex(reflect.ValueOf(pair.Key), valelm.Elem())
		}
		vl.Set(vmap)
	} else {
		fmt.Println("LikERROR: export map only by string keys")
		result = false
	}
	return result
}

func list_to_slice(list Lister, tp reflect.Type, vl reflect.Value) bool {
	if list == nil || !vl.IsValid() {
		return false
	}
	result := true
	tpelm := tp.Elem()
	tpsli := reflect.SliceOf(tpelm)
	count := list.Count()
	slice := reflect.MakeSlice(tpsli, count, count)
	for n := 0; n < list.Count(); n++ {
		if item := list.GetItem(n); item != nil {
			if !item_to_reflect(item, tpelm, slice.Index(n)) {
				result = false
			}
		}
	}
	vl.Set(slice)
	return result
}
