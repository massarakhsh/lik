package lik

import (
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

func SetFromRequest(data string) Seter {
	result := BuildSet()
	data = strings.Trim(data, " \n\r\t\b")
	if data == "" {
	} else if strings.HasPrefix(data, "{") {
		pars := buildParse(data)
		result = pars.scanItMap()
	} else {
		words := strings.Split(data, "&")
		for _, word := range words {
			if peq := strings.Index(word, "="); peq > 0 {
				key := word[0:peq]
				val := word[peq+1:]
				pars := buildParse(val)
				item := pars.scanItValue()
				if key != "" {
					result.SetValue(key, item)
				}
			}
		}
	}
	return result
}

func SetFromMap(data map[string]interface{}) Seter {
	result := BuildSet()
	for key, val := range data {
		result.SetValue(key, val)
	}
	return result
}

func SetFromStruct(data interface{}) Seter {
	return BuildItem(data).ToSet()
}

func ListFromRequest(data string) Lister {
	result := BuildList()
	data = strings.Trim(data, " \n\r\t")
	if data == "" {
	} else if strings.HasPrefix(data, "[") {
		pars := buildParse(data)
		result = pars.scanItList()
	}
	return result
}

func SetToStruct[T any](set Seter) T {
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
			nm := ToSnakeCase(name)
			item = set.GetItem(nm)
		}
		if item != nil {
			if val := vl.Field(f); val.IsValid() {
				if stp := tp.Field(f).Type.Kind().String(); stp == "string" {
					val.Set(reflect.ValueOf(item.ToString()))
				} else if stp == "int64" {
					val.Set(reflect.ValueOf(item.ToInt()))
				} else if stp == "int" {
					val.Set(reflect.ValueOf(int(item.ToInt())))
				} else if stp == "slice" {
					if list := item.ToList(); list != nil {
						var strs []string
						for n := 0; n < list.Count(); n++ {
							strs = append(strs, list.GetString(n))
						}
						val.Set(reflect.ValueOf(strs))
					}
				}
			}
		}
	}
	return res
}
