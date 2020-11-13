package lik

import "strings"

func XML_ListToString(prefix string, list Lister) string {
	dump := ""
	if list != nil {
		for ne := 0; ne < list.Count(); ne++ {
			if str := list.GetString(ne); str != "" {
				dump += prefix + str + "\r\n"
			} else if elm := list.GetSet(ne); elm != nil {
				dump += XML_ElementToString(prefix, elm)
			}
		}
	}
	return dump
}

func XML_ElementToString(prefix string, elm Seter) string {
	dump := ""
	if elm != nil {
		if tag := elm.GetString("_tag"); tag != "" {
			value := elm.GetString("_value")
			content := elm.GetList("_content")
			if value != "" || content != nil || elm.Count() > 1 {
				dump += prefix + "<" + tag
				for _, atr := range elm.Values() {
					if !strings.HasPrefix(atr.Key, "_") {
						if str := atr.Val.ToString(); str != "" {
							dump += " " + atr.Key + "="
							dump += StrToQuotes(str)
						}
					}
				}
				if value == "" && content == nil {
					dump += " />\r\n"
				} else {
					dump += ">"
					if value != "" {
						dump += value
					}
					if content != nil {
						dump += "\r\n"
						dump += XML_ListToString(prefix+"    ", content)
						dump += prefix
					}
					dump += "</" + tag + ">\r\n"
				}
			}
		}
	}
	return dump
}
