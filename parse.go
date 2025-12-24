package lik

import (
	"fmt"
	"strings"
)

const ZERO = '\x00'

type ParseJson struct {
	Source []rune
	Pos    int
	Len    int
	stoped bool
}

var DebugPrint = false
var DebugDiagnos = ""

func buildParse(data string) *ParseJson {
	parse := &ParseJson{}
	parse.Source = []rune(data)
	parse.Len = len(parse.Source)
	return parse
}

func (pars *ParseJson) scanValue() Itemer {
	ch := pars.getNextRune()
	switch ch {
	case '{':
		return pars.scanMap()
	case '[':
		return pars.scanList()
	case '"', '\'':
		str := pars.scanString()
		return &DItemString{str}
	default:
		str := pars.scanImmediate()
		if lstr := strings.ToLower(str); lstr == "true" {
			return &DItemBool{true}
		} else if lstr == "false" {
			return &DItemBool{false}
		} else if lstr == "null" || lstr == "nil" {
			return nil
		} else if ival, ok := StrToInt64If(str); ok {
			return &DItemInt{ival}
		} else if fval, ok := StrToFloatIf(str); ok {
			return &DItemFloat{fval}
		} else {
			pars.printError("need value")
			return nil
		}
	}
}

func (pars *ParseJson) scanMap() Seter {
	info := BuildSet()
	if pars.getNextRune() != '{' {
		pars.printError("need '{'")
		return nil
	}
	pars.stepNextRune()
	for {
		if pars.stoped {
			return nil
		}
		ch := pars.getNextRune()
		if ch == ZERO || ch == '}' {
			break
		}
		if ch != '"' && ch != '\'' {
			pars.printError("need key")
			return nil
		}
		key := pars.scanString()
		if key == "" {
			pars.printError("key can not be empty")
			return nil
		}
		ch = pars.getNextRune()
		if ch != ':' {
			pars.printError("need `:`")
			break
		}
		pars.stepNextRune()
		item := pars.scanValue()
		if item != nil {
			info.SetValue(key, item)
		}
		ch = pars.getNextRune()
		if ch == ',' {
			pars.stepNextRune()
		} else {
			break
		}
	}
	if pars.getNextRune() != '}' {
		pars.printError("need '}'")
		return nil
	}
	pars.stepNextRune()
	if pars.stoped {
		return nil
	}
	return info
}

func (pars *ParseJson) scanList() Lister {
	info := BuildList()
	if pars.getNextRune() != '[' {
		pars.printError("need '['")
		return nil
	}
	pars.stepNextRune()
	for {
		if pars.stoped {
			return nil
		}
		ch := pars.getNextRune()
		if ch == ZERO || ch == ']' {
			break
		}
		item := pars.scanValue()
		if item != nil {
			info.AddItems(item)
		}
		ch = pars.getNextRune()
		if ch == ',' {
			pars.stepNextRune()
		} else {
			break
		}
	}
	if pars.getNextRune() != ']' {
		pars.printError("need ']'")
		return nil
	}
	pars.stepNextRune()
	if pars.stoped {
		return nil
	}
	return info
}

func (pars *ParseJson) scanString() string {
	info := ""
	chg := pars.getNextRune()
	pars.stepNextRune()
	for ; pars.Pos < pars.Len; pars.Pos++ {
		ch := pars.Source[pars.Pos]
		if ch == chg {
			pars.Pos++
			break
		}
		if ch == '\\' {
			pars.Pos++
			if pars.Pos >= pars.Len {
				pars.printError("need symbol")
				return ""
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

func (pars *ParseJson) scanImmediate() string {
	info := ""
	for ; pars.Pos < pars.Len && !pars.stoped; pars.Pos++ {
		ch := pars.Source[pars.Pos]
		if ch == ZERO || ch == ' ' || ch < 0x20 || ch == ',' || ch == ':' || ch == '=' || ch == '}' || ch == ']' {
			break
		}
		info += string(ch)
	}
	return info
}

func (pars *ParseJson) getNextRune() rune {
	for ; pars.Pos < pars.Len; pars.Pos++ {
		ch := pars.Source[pars.Pos]
		if ch != ZERO && ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r' {
			return ch
		}
	}
	return ZERO
}

func (pars *ParseJson) stepNextRune() {
	if pars.Pos < pars.Len {
		pars.Pos++
	}
}

func (pars *ParseJson) printError(diag string) {
	text := fmt.Sprintf("Parsing error %s pos %d: ", diag, pars.Pos)
	for dep := -20; dep < 20; dep++ {
		if dep == 0 {
			text += " <<>> "
		}
		if pos := pars.Pos + dep; pos >= 0 && pos < pars.Len {
			text += string(pars.Source[pos])
		}
	}
	DebugDiagnos = text
	if DebugPrint {
		fmt.Printf("%s\n", text)
	}
	pars.stoped = true
}
