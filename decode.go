package lik

import (
	"os"
	"strings"
)

func SetFromFile(filename string) Seter {
	var set Seter
	if item := ItemFromFile(filename); item != nil {
		set = item.ToSet()
	}
	return set
}

func ItemFromFile(filename string) Itemer {
	var item Itemer
	if data, err := os.ReadFile(filename); err == nil {
		item = ItemFromString(string(data))
	}
	return item
}

func SetFromString(data string) Seter {
	var set Seter
	if item := ItemFromString(data); item != nil {
		set = item.ToSet()
	}
	return set
}

func ItemFromString(data string) Itemer {
	str := strings.Trim(data, " \n\r\t\b")
	if str == "" {
		return nil
	}
	pars := buildParse(str)
	item := pars.scanValue()
	if pars.stoped {
		return nil
	}
	return item
}

func SetFromBytes(code []byte) Seter {
	return SetFromString(string(code))
}

func SetFromRequest(data string) Seter {
	var set Seter
	data = strings.Trim(data, " \n\r\t\b")
	if strings.HasPrefix(data, "{") {
		set = SetFromString(data)
	} else {
		set = SetFromQuery(data)
	}
	return set
}

func SetFromQuery(data string) Seter {
	set := BuildSet()
	data = strings.Trim(data, " \n\r\t\b")
	words := strings.Split(data, "&")
	for _, word := range words {
		if peq := strings.Index(word, "="); peq > 0 {
			key := word[0:peq]
			if key != "" {
				val := word[peq+1:]
				if val != "" && (val[0] == '{' || val[0] == '[') {
					pars := buildParse(val)
					item := pars.scanValue()
					if item != nil {
						set.SetValue(key, item)
					}
				} else {
					set.SetString(key, val)
				}
			}
		}
	}
	return set
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
