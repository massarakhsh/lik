package main

import (
	"fmt"
	"reflect"

	"github.com/massarakhsh/lik"
)

type TP struct {
	Name   string
	Format string
}

func main() {
	data := TP{Name: "Имя", Format: "Форма"}
	set := lik.BuildSet(data)
	fmt.Println(set.Serialize())
	data1 := SetToStruct[TP](set)
	if data == data1 {
		fmt.Println("Ok")
	} else {
		fmt.Println("Fault")
	}
}

func SetToStruct[T any](set lik.Seter) T {
	var res T
	result := &res
	tp := reflect.TypeOf(result)
	vl := reflect.ValueOf(result)
	if vl.Kind() != reflect.Pointer || vl.IsNil() {
		return res
	}
	tp = tp.Elem()
	vl = vl.Elem()
	cnt := vl.NumField()
	for f := 0; f < cnt; f++ {
		name := tp.Field(f).Name
		item := set.GetItem(name)
		if item == nil {
			nm := lik.ToSnakeCase(name)
			item = set.GetItem(nm)
		}
		if item != nil {
			if val := vl.Field(f); val.IsValid() {
				if stp := tp.Field(f).Type.Kind().String(); stp == "string" {
					val.Set(reflect.ValueOf(item.ToString()))
				} else {
					_ = stp
				}
			}
		}
	}
	return res
}
