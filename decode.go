package lik

import "strings"

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
