package lik

import (
	"fmt"
	"strings"
)

const (
	Dbg  = false
	ZERO = '\x00'
)

type JsonParse struct {
	Source []rune
	Pos    int
	Len    int
}

type JsonParser interface {
	scanItValue() Itemer
	scanItMap() Seter
	scanItList() Lister
	scanItString() string
	scanItImmediate() string
	getNextRune() rune
	stepNextRule()
}

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

func buildParse(data string) *JsonParse {
	parse := JsonParse{}
	parse.Source = []rune(data)
	parse.Len = len(parse.Source)
	return &parse
}

func (pars *JsonParse) scanItValue() Itemer {
	ch := pars.getNextRune()
	if ch == '{' {
		return pars.scanItMap()
	} else if ch == '[' {
		return pars.scanItList()
	} else if ch == '"' || ch == '\'' {
		str := pars.scanItString()
		return &DItemString{str}
	} else {
		str := pars.scanItImmediate()
		/*if str == "null" || str == "NULL" {
			return nil
		} else*/if str == "true" || str == "TRUE" {
			return &DItemBool{true}
		} else if str == "false" || str == "FALSE" {
			return &DItemBool{false}
		} else if ival, ok := StrToInt64If(str); ok {
			return &DItemInt{ival}
		} else if fval, ok := StrToFloatIf(str); ok {
			return &DItemFloat{fval}
		} else {
			return &DItemString{str}
		}
	}
}

func (pars *JsonParse) scanItMap() Seter {
	if Dbg {
		fmt.Printf("scanItMap: %d\n", pars.Pos)
	}
	info := BuildSet()
	if pars.getNextRune() == '{' {
		pars.stepNextRune()
	}
	for {
		ch := pars.getNextRune()
		if ch == ZERO || ch == '}' {
			break
		}
		var key string
		if ch == '"' || ch == '\'' {
			key = pars.scanItString()
		} else {
			key = pars.scanItImmediate()
		}
		if key == "" {
			break
		}
		if key == "definition" {
			key += ""
		}
		ch = pars.getNextRune()
		if ch != '=' && ch != ':' {
			break
		}
		pars.stepNextRune()
		item := pars.scanItValue()
		if item != nil {
			info.SetValue(key, item)
		}
		ch = pars.getNextRune()
		if ch == ',' {
			pars.stepNextRune()
		}
	}
	if pars.getNextRune() == '}' {
		pars.stepNextRune()
	}
	return info
}

func (pars *JsonParse) scanItList() Lister {
	info := BuildList()
	if pars.getNextRune() == '[' {
		pars.stepNextRune()
	}
	for {
		ch := pars.getNextRune()
		if ch == ZERO || ch == ']' {
			break
		}
		item := pars.scanItValue()
		if item == nil {
			break
		}
		info.AddItems(item)
		ch = pars.getNextRune()
		if ch == ',' {
			pars.stepNextRune()
		}
	}
	if pars.getNextRune() == ']' {
		pars.stepNextRune()
	}
	return info
}

func (pars *JsonParse) scanItString() string {
	info := ""
	chg := pars.getNextRune()
	pars.stepNextRune()
	for ; pars.Pos < pars.Len; pars.Pos++ {
		ch := pars.Source[pars.Pos]
		if ch == chg {
			pars.Pos++
			break
		} else if ch == '\\' {
			pars.Pos++
			if pars.Pos >= pars.Len {
				break
			}
			chn := pars.Source[pars.Pos]
			if chn == 'r' {
				ch = '\r'
			} else if chn == 'n' {
				ch = '\n'
			} else if chn == 't' {
				ch = '\t'
			} else if chn == 'b' {
				ch = '\b'
			} else {
				ch = chn
			}
		}
		info += string(ch)
	}
	return info
}

func (pars *JsonParse) scanItImmediate() string {
	info := ""
	for ; pars.Pos < pars.Len; pars.Pos++ {
		ch := pars.Source[pars.Pos]
		if ch == ZERO || ch == ' ' || ch < 0x20 || ch == ',' || ch == ':' || ch == '=' || ch == '}' || ch == ']' {
			break
		}
		info += string(ch)
	}
	return info
}

func (pars *JsonParse) getNextRune() rune {
	if Dbg {
		fmt.Printf("getNextRune: %d\n", pars.Pos)
	}
	for ; pars.Pos < pars.Len; pars.Pos++ {
		ch := pars.Source[pars.Pos]
		if ch != ZERO && ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r' {
			return ch
		}
	}
	return ZERO
}

func (pars *JsonParse) stepNextRune() {
	if pars.Pos < pars.Len {
		pars.Pos++
	}
}
