package lik

import "fmt"

const (
	Dbg  = false
	ZERO = '\x00'
)

type JsonParse struct {
	Source []rune
	Pos    int
	Len    int
}

func buildParse(data string) *JsonParse {
	parse := JsonParse{}
	parse.Source = []rune(data)
	parse.Len = len(parse.Source)
	return &parse
}

func (pars *JsonParse) scanValue() Itemer {
	ch := pars.getNextRune()
	if ch == '{' {
		return pars.scanMap()
	} else if ch == '[' {
		return pars.scanList()
	} else if ch == '"' || ch == '\'' {
		str := pars.scanString()
		return &DItemString{str}
	} else {
		str := pars.scanImmediate()
		if str == "true" || str == "True" || str == "TRUE" {
			return &DItemBool{true}
		} else if str == "false" || str == "False" || str == "FALSE" {
			return &DItemBool{false}
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

func (pars *JsonParse) scanMap() Seter {
	if Dbg {
		fmt.Printf("scanItMap: %d\n", pars.Pos)
	}
	info := BuildSet()
	if pars.getNextRune() != '{' {
		pars.printError("need '{'")
		return nil
	}
	pars.stepNextRune()
	for {
		ch := pars.getNextRune()
		if ch == ZERO || ch == '}' {
			break
		}
		if ch != '"' || ch != '\'' {
			pars.printError("need key")
			return nil
		}
		key := pars.scanString()
		if key == "" {
			pars.printError("key can not be empty")
			return nil
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
	if pars.getNextRune() != '}' {
		pars.printError("need '}'")
		return nil
	}
	pars.stepNextRune()
	return info
}

func (pars *JsonParse) scanList() Lister {
	info := BuildList()
	if pars.getNextRune() != '[' {
		pars.printError("need '['")
		return nil
	}
	pars.stepNextRune()
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
		} else {
			break
		}
	}
	if pars.getNextRune() != ']' {
		pars.printError("need ']'")
		return nil
	}
	pars.stepNextRune()
	return info
}

func (pars *JsonParse) scanString() string {
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

func (pars *JsonParse) scanImmediate() string {
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

func (pars *JsonParse) printError(diag string) {
	fmt.Printf("Parsing error %s pos %d: ", diag, pars.Pos)
	text := " <<<"
	for pos := pars.Pos; pos >= 0; pos-- {
		if pos < pars.Len {
			text = string(pars.Source[pos]) + text
		}
	}
	fmt.Printf("%s\n", text)
}
