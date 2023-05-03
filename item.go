package lik

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

type IDB int

// Базовый интерфейс динамических элементов
type Itemer interface {
	IsBool() bool
	IsInt() bool
	IsFloat() bool
	IsString() bool
	IsList() bool
	IsSet() bool
	ToBool() bool
	ToInt() int64
	ToFloat() float64
	ToString() string
	ToList() Lister
	ToSet() Seter
	Serialize() string
	Format(prefix string) string
	Clone() Itemer
}

type DItemBool struct {
	Val bool
}

type DItemInt struct {
	Val int64
}

type DItemFloat struct {
	Val float64
}

type DItemString struct {
	Val string
}

type DItemList struct {
	Val []Itemer
}

type DItemSet struct {
	Val []SetElm
}

type SetElm struct {
	Key string
	Val Itemer
}

// Interface Itemer

// Создание динамического элемента из параметра
// Тип элемента определяется типом параметра
func BuildItem(data interface{}) Itemer {
	var item Itemer
	switch val := data.(type) {
	case DItemBool:
		item = &DItemBool{val.Val}
	case *DItemBool:
		item = &DItemBool{val.Val}
	case DItemInt:
		item = &DItemInt{val.Val}
	case *DItemInt:
		item = &DItemInt{val.Val}
	case DItemFloat:
		item = &DItemFloat{val.Val}
	case *DItemFloat:
		item = &DItemFloat{val.Val}
	case DItemString:
		item = &DItemString{val.Val}
	case *DItemString:
		item = &DItemString{val.Val}
	case DItemList:
		item = &val
	case *DItemList:
		item = val
	case Lister:
		item = val
	case DItemSet:
		item = &val
	case *DItemSet:
		item = val
	case Seter:
		item = val
	default:
		tp := reflect.TypeOf(data)
		vl := reflect.ValueOf(data)
		if vl.Kind() == reflect.Ptr {
			tp = tp.Elem()
			vl = vl.Elem()
		}
		tps := tp.Kind().String()
		if tps == "bool" {
			item = &DItemBool{vl.Bool()}
		} else if tps == "int" {
			item = &DItemInt{vl.Int()}
		} else if tps == "int32" {
			item = &DItemInt{vl.Int()}
		} else if tps == "int64" {
			item = &DItemInt{vl.Int()}
		} else if tps == "uint" {
			item = &DItemInt{vl.Int()}
		} else if tps == "uint32" {
			item = &DItemInt{vl.Int()}
		} else if tps == "uint64" {
			item = &DItemInt{vl.Int()}
		} else if tps == "float32" {
			item = &DItemFloat{vl.Float()}
		} else if tps == "float64" {
			item = &DItemFloat{vl.Float()}
		} else if tps == "string" {
			item = &DItemString{vl.String()}
		} else if tps == "struct" {
			item = SetFromReflectStructure(tp, vl)
		} else if tps == "slice" {
			item = SetFromReflectSlice(tp, vl)
		} else {
			//fmt.Println("BuildItem ERROR: ", data)
		}
	}
	return item
}

func BuildItemReflect(val reflect.Value) Itemer {
	var item Itemer
	if !val.IsValid() {
	} else if val.CanInt() {
		item = &DItemInt{val.Int()}
	} else if val.CanFloat() {
		item = &DItemFloat{val.Float()}
	} else if val.CanConvert(reflect.TypeOf("")) {
		item = &DItemString{val.String()}
	} else if val.CanInterface() {
		item = BuildItem(val.Interface())
	}
	return item
}

func SetFromReflectStructure(tp reflect.Type, vl reflect.Value) Seter {
	set := BuildSet()
	cnt := vl.NumField()
	for f := 0; f < cnt; f++ {
		nam := tp.Field(f).Name
		if val := vl.Field(f); val.IsValid() {
			if !val.IsZero() {
				if item := BuildItemReflect(val); item != nil {
					set.SetValue(nam, item)
				}
			}
		}
	}
	return set
}

func SetFromReflectSlice(tp reflect.Type, vl reflect.Value) Lister {
	list := BuildList()
	cnt := vl.Len()
	for n := 0; n < cnt; n++ {
		if elm := vl.Index(n); elm.IsValid() {
			if item := BuildItemReflect(elm); item != nil {
				list.AddItems(item)
			}
		}
	}
	return list
}

func (it *DItemBool) IsBool() bool {
	return true
}
func (it *DItemInt) IsBool() bool {
	return false
}
func (it *DItemFloat) IsBool() bool {
	return false
}
func (it *DItemString) IsBool() bool {
	return false
}
func (it *DItemList) IsBool() bool {
	return false
}
func (it *DItemSet) IsBool() bool {
	return false
}

func (it *DItemBool) IsInt() bool {
	return false
}
func (it *DItemInt) IsInt() bool {
	return true
}
func (it *DItemFloat) IsInt() bool {
	return false
}
func (it *DItemString) IsInt() bool {
	return false
}
func (it *DItemList) IsInt() bool {
	return false
}
func (it *DItemSet) IsInt() bool {
	return false
}

func (it *DItemBool) IsFloat() bool {
	return false
}
func (it *DItemInt) IsFloat() bool {
	return false
}
func (it *DItemFloat) IsFloat() bool {
	return true
}
func (it *DItemString) IsFloat() bool {
	return false
}
func (it *DItemList) IsFloat() bool {
	return false
}
func (it *DItemSet) IsFloat() bool {
	return false
}

func (it *DItemBool) IsString() bool {
	return false
}
func (it *DItemInt) IsString() bool {
	return false
}
func (it *DItemFloat) IsString() bool {
	return false
}
func (it *DItemString) IsString() bool {
	return true
}
func (it *DItemList) IsString() bool {
	return false
}
func (it *DItemSet) IsString() bool {
	return false
}

func (it *DItemBool) IsList() bool {
	return false
}
func (it *DItemInt) IsList() bool {
	return false
}
func (it *DItemFloat) IsList() bool {
	return false
}
func (it *DItemString) IsList() bool {
	return false
}
func (it *DItemList) IsList() bool {
	return true
}
func (it *DItemSet) IsList() bool {
	return false
}

func (it *DItemBool) IsSet() bool {
	return false
}
func (it *DItemInt) IsSet() bool {
	return false
}
func (it *DItemFloat) IsSet() bool {
	return false
}
func (it *DItemString) IsSet() bool {
	return false
}
func (it *DItemList) IsSet() bool {
	return false
}
func (it *DItemSet) IsSet() bool {
	return true
}

func (it *DItemBool) ToBool() bool {
	return it.Val
}
func (it *DItemInt) ToBool() bool {
	if it.Val > 0 {
		return true
	} else {
		return false
	}
}
func (it *DItemFloat) ToBool() bool {
	if it.Val > 0 {
		return true
	} else {
		return false
	}
}
func (it *DItemString) ToBool() bool {
	low := strings.ToLower(it.Val)
	if low == "true" {
		return true
	} else if low == "false" {
		return false
	}
	return false
}
func (it *DItemList) ToBool() bool {
	return false
}
func (it *DItemSet) ToBool() bool {
	return false
}

func (it *DItemBool) ToInt() int64 {
	if it.Val {
		return 1
	} else {
		return 0
	}
}
func (it *DItemInt) ToInt() int64 {
	return it.Val
}
func (it *DItemFloat) ToInt() int64 {
	return int64(math.Round(it.Val))
}
func (it *DItemString) ToInt() int64 {
	return StrToInt64(it.Val)
}
func (it *DItemList) ToInt() int64 {
	return 0
}
func (it *DItemSet) ToInt() int64 {
	return 0
}

func (it *DItemBool) ToFloat() float64 {
	if it.Val {
		return 1
	} else {
		return 0
	}
}
func (it *DItemInt) ToFloat() float64 {
	return float64(it.Val)
}
func (it *DItemFloat) ToFloat() float64 {
	return it.Val
}
func (it *DItemString) ToFloat() float64 {
	return StrToFloat(it.Val)
}
func (it *DItemList) ToFloat() float64 {
	return 0
}
func (it *DItemSet) ToFloat() float64 {
	return 0
}

func (it *DItemBool) ToString() string {
	if it.Val {
		return "true"
	} else {
		return "false"
	}
}
func (it *DItemInt) ToString() string {
	return fmt.Sprint(it.Val)
}
func (it *DItemFloat) ToString() string {
	return fmt.Sprint(it.Val)
}
func (it *DItemString) ToString() string {
	return it.Val
}
func (it *DItemList) ToString() string {
	return ""
}
func (it *DItemSet) ToString() string {
	return ""
}

func (it *DItemBool) ToList() Lister {
	return nil
}
func (it *DItemInt) ToList() Lister {
	return nil
}
func (it *DItemFloat) ToList() Lister {
	return nil
}
func (it *DItemString) ToList() Lister {
	return nil
}
func (it *DItemList) ToList() Lister {
	return it
}
func (it *DItemSet) ToList() Lister {
	return nil
}

func (it *DItemBool) ToSet() Seter {
	return nil
}
func (it *DItemInt) ToSet() Seter {
	return nil
}
func (it *DItemFloat) ToSet() Seter {
	return nil
}
func (it *DItemString) ToSet() Seter {
	return nil
}
func (it *DItemList) ToSet() Seter {
	return nil
}
func (it *DItemSet) ToSet() Seter {
	return it
}

func (it *DItemBool) Clone() Itemer {
	return &DItemBool{it.Val}
}
func (it *DItemInt) Clone() Itemer {
	return &DItemInt{it.Val}
}
func (it *DItemFloat) Clone() Itemer {
	return &DItemFloat{it.Val}
}
func (it *DItemString) Clone() Itemer {
	return &DItemString{it.Val}
}
func (it *DItemList) Clone() Itemer {
	return it.clone()
}
func (it *DItemSet) Clone() Itemer {
	return it.clone()
}

func (it *DItemBool) Serialize() string {
	return fmt.Sprint(it.Val)
}
func (it *DItemInt) Serialize() string {
	return fmt.Sprint(it.Val)
}
func (it *DItemFloat) Serialize() string {
	return fmt.Sprint(it.Val)
}
func (it *DItemString) Serialize() string {
	return StrToQuotes(it.Val)
}
func (it *DItemList) Serialize() string {
	return it.serialize()
}
func (it *DItemSet) Serialize() string {
	return it.serialize()
}

func (it *DItemBool) Format(prefix string) string {
	return it.Serialize()
}
func (it *DItemInt) Format(prefix string) string {
	return it.Serialize()
}
func (it *DItemFloat) Format(prefix string) string {
	return it.Serialize()
}
func (it *DItemString) Format(prefix string) string {
	return it.Serialize()
}
func (it *DItemList) Format(prefix string) string {
	return it.format(prefix)
}
func (it *DItemSet) Format(prefix string) string {
	return it.format(prefix)
}

/////////////////////////////////////////////

func GetSetString(item Seter, path string) string {
	val := ""
	if item != nil {
		val = item.GetString(path)
	}
	return val
}

func getInfoItem(info Itemer, path string) Itemer {
	var value Itemer
	if info == nil {
	} else if name, ext := GetFirstExt(path); name == "" && ext == "" {
		value = info
	} else if name == "" {
		value = getInfoItem(info, ext)
	} else if imap := info.ToSet(); imap != nil {
		value = getInfoItem(imap.GetItem(name), ext)
	} else if ilist := info.ToList(); ilist != nil {
		if idx, ok := StrToIntIf(name); ok {
			value = getInfoItem(ilist.GetItem(idx), ext)
		}
	}
	return value
}

func infoToString(info Itemer) string {
	value := ""
	if info == nil {
	} else if iset := info.ToSet(); iset != nil {
		for _, val := range iset.Values() {
			dt := infoToString(val.Val)
			if dt != "" {
				if value != "" {
					value += " "
				}
				value += dt
			}
		}
	} else if ilist := info.ToList(); ilist != nil {
		for _, sub := range ilist.Values() {
			dt := infoToString(sub)
			if dt != "" {
				if value != "" {
					value += ", "
				}
				value += dt
			}
		}
	} else {
		value = info.ToString()
	}
	return value
}

func setInfoItem(val interface{}, info Itemer, path string) bool {
	modify := false
	if info == nil {
	} else if name, ext := GetFirstExt(path); name == "" && ext == "" {
	} else if name == "" {
		modify = setInfoItem(val, info, ext)
	} else if imap := info.ToSet(); imap != nil {
		if ext == "" {
			if imap.SetValue(name, val) {
				modify = true
			}
		} else if item := imap.GetItem(name); item != nil {
			if setInfoItem(val, item, ext) {
				modify = true
			}
		} else if RegExCompare(ext, "^(\\d+)") {
			modify = true
			item := BuildList()
			imap.SetValue(name, item)
			setInfoItem(val, item, ext)
		} else {
			modify = true
			item := BuildSet()
			imap.SetValue(name, item)
			setInfoItem(val, item, ext)
		}
	} else if ilist := info.ToList(); ilist != nil {
		if idx, ok := StrToIntIf(name); ok {
			if ext == "" {
				if ilist.SetValue(idx, BuildItem(val)) {
					modify = true
				}
			} else if item := ilist.GetItem(idx); item != nil {
				if setInfoItem(val, item, ext) {
					modify = true
				}
			} else if RegExCompare(ext, "^(\\d+)") {
				modify = true
				item := BuildList()
				ilist.SetValue(idx, item)
				setInfoItem(val, item, ext)
			} else {
				modify = true
				item := BuildSet()
				ilist.SetValue(idx, item)
				setInfoItem(val, item, ext)
			}
		}
	}
	return modify
}
