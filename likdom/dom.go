package likdom

import (
	"regexp"
	"strings"
)

type LikDom struct {
	Tag    string
	Unpair bool
	Attr   map[string]string
	Text   string
	Data   []Domer
}

type Domer interface {
	GetTag() string
	IsAttr(key string) bool
	GetAttr(key string) string
	SetAttr(attr ...string) Domer
	ToString() string
	ToPrefixString(prefix string) string
	GetDataCount() int
	GetDataItem(num int) (Domer, bool)
	GetDataTag(tag string) (Domer, bool)
	AppendItem(items ...Domer)
	BuildSpace() Domer
	BuildString(text ...string)
	BuildItem(tag string, attr ...string) Domer
	BuildItemClass(tag string, class string, attr ...string) Domer
	BuildUnpairItem(tag string, attr ...string) Domer
	BuildTable(attr ...string) Domer
	BuildTableClass(class string, attr ...string) Domer
	BuildTr(attr ...string) Domer
	BuildTrClass(class string, attr ...string) Domer
	BuildTd(attr ...string) Domer
	BuildTdClass(class string, attr ...string) Domer
	BuildTrTd(attr ...string) Domer
	BuildTrTdClass(class string, attr ...string) Domer
	BuildDiv(attr ...string) Domer
	BuildDivClass(class string, attr ...string) Domer
	BuildTextClassProc(text string, class string, proc string, attr ...string) Domer
}

func BuildSpace() Domer {
	dom := &LikDom{Attr: make(map[string]string)}
	return dom
}

func BuildString(text string) Domer {
	dom := &LikDom{Attr: make(map[string]string), Text: text}
	return dom
}

func BuildItem(tag string, attr ...string) Domer {
	dom := &LikDom{Tag: tag, Attr: make(map[string]string)}
	dom.SetAttr(attr...)
	return dom
}

func BuildItemClass(tag string, class string, attr ...string) Domer {
	item := &LikDom{Tag: tag, Attr: make(map[string]string)}
	if class != "" {
		item.SetAttr("class", class)
	}
	item.SetAttr(attr...)
	return item
}

func BuildUnpairItem(tag string, attr ...string) Domer {
	item := &LikDom{Tag: tag, Unpair: true, Attr: make(map[string]string)}
	item.SetAttr(attr...)
	return item
}

func BuildPageHtml() Domer {
	item := BuildItem("html")
	item.BuildItem("head")
	item.BuildItem("body")
	return item
}

func BuildDiv(attr ...string) Domer {
	return BuildItem("div", attr...)
}
func BuildDivClass(class string, attr ...string) Domer {
	return BuildItemClass("div", class, attr...)
}
func BuildDivClassId(class string, id string, attr ...string) Domer {
	item := BuildDivClass(class, attr...)
	if id != "" {
		item.SetAttr("id", id)
	}
	return item
}

func BuildTable(attr ...string) Domer {
	return BuildItem("table", attr...)
}
func BuildTableClass(class string, attr ...string) Domer {
	return BuildItemClass("table", class, attr...)
}
func BuildTableClassId(class string, id string, attr ...string) Domer {
	item := BuildTableClass(class, attr...)
	if id != "" {
		item.SetAttr("id", id)
	}
	return item
}

func BuildImgProc(class string, id string, proc string, img string, title string) Domer {
	item := BuildUnpairItem("img", "src", img)
	if class != "" {
		item.SetAttr("class", class)
	}
	if id != "" {
		item.SetAttr("id", id)
	}
	if proc != "" {
		item.SetAttr("onclick", proc)
	}
	if title != "" {
		item.SetAttr("title", title)
	}
	return item
}

func (dom *LikDom) GetTag() string {
	return dom.Tag
}

func (dom *LikDom) IsAttr(key string) bool {
	_, ok := dom.Attr[key]
	return ok
}

func (dom *LikDom) GetAttr(key string) string {
	val, _ := dom.Attr[key]
	return val
}

func (dom *LikDom) SetAttr(attr ...string) Domer {
	for na := 0; na < len(attr); na++ {
		key := attr[na]
		val := ""
		if match := regexp.MustCompile("^(.+?)=(.*)").FindStringSubmatch(key); match != nil {
			key = match[1]
			val = match[2]
		} else if na+1 < len(attr) {
			na++
			val = attr[na]
		}
		if len(key) > 0 {
			val = strings.Trim(val, "'\"")
			dom.Attr[key] = val
		} else if len(val) > 0 {
			dom.Attr[val] = ""
		}
	}
	return dom
}

func (dom *LikDom) AppendItem(items ...Domer) {
	for _, item := range items {
		if item != nil {
			dom.Data = append(dom.Data, item)
		}
	}
}

func (dom *LikDom) BuildItem(tag string, attr ...string) Domer {
	item := BuildItem(tag, attr...)
	dom.AppendItem(item)
	return item
}

func (dom *LikDom) BuildItemClass(tag string, class string, attr ...string) Domer {
	item := BuildItemClass(tag, class, attr...)
	dom.AppendItem(item)
	return item
}

func (dom *LikDom) BuildUnpairItem(tag string, attr ...string) Domer {
	item := BuildUnpairItem(tag, attr...)
	dom.AppendItem(item)
	return item
}

func (dom *LikDom) BuildString(text ...string) {
	for _, txt := range text {
		item := BuildString(txt)
		dom.AppendItem(item)
	}
}

func (dom *LikDom) BuildSpace() Domer {
	item := BuildSpace()
	dom.AppendItem(item)
	return item
}

func (dom *LikDom) GetDataCount() int {
	return len(dom.Data)
}

func (dom *LikDom) GetDataItem(num int) (Domer, bool) {
	if num >= 0 && num < len(dom.Data) {
		elm := dom.Data[num]
		return elm, true
	}
	return nil, false
}

func (dom *LikDom) GetDataTag(tag string) (Domer, bool) {
	for n := 0; n < len(dom.Data); n++ {
		item := dom.Data[n]
		if item != nil && item.GetTag() == tag {
			return item, true
		}
	}
	return nil, false
}

func (dom *LikDom) ToString() string {
	return dom.ToPrefixString("")
}

func (dom *LikDom) ToPrefixString(prefix string) string {
	code := ""
	if dom.Tag != "" {
		tag := dom.Tag
		tags := tag
		if dom.Text != "" {
			tags += " " + dom.Text
		}
		for key, val := range dom.Attr {
			tags += " " + key
			if val != "" {
				qval := strings.Replace(val, "\\", "\\\\", -1)
				qval = strings.Replace(qval, "\"", "\\\"", -1)
				qval = strings.Replace(qval, "\n", "\\n", -1)
				tags += "=\"" + qval + "\""
			}
		}
		code += prefix + "<" + tags + ">"
		if dom.Unpair {
			code += "\n"
		} else if strings.ToLower(tag) == "textarea" {
			for _, item := range dom.Data {
				if item != nil {
					code += item.ToString()
				}
			}
			code += "</" + tag + ">\n"
		} else {
			code += "\n"
			for _, item := range dom.Data {
				if item != nil {
					code += item.ToPrefixString(prefix + "    ")
				}
			}
			code += prefix + "</" + tag + ">\n"
		}
	} else if dom.Text != "" {
		code += prefix + dom.Text + "\n"
	} else {
		for _, item := range dom.Data {
			if item != nil {
				code += item.ToPrefixString(prefix)
			}
		}
	}
	return code
}

func (dom *LikDom) BuildTable(attr ...string) Domer {
	return dom.BuildItem("table", attr...)
}

func (dom *LikDom) BuildTableClass(class string, attr ...string) Domer {
	return dom.BuildTagClass("table", class, attr...)
}

func (dom *LikDom) BuildTr(attr ...string) Domer {
	return dom.BuildItem("tr", attr...)
}

func (dom *LikDom) BuildTrClass(class string, attr ...string) Domer {
	return dom.BuildTagClass("tr", class, attr...)
}

func (dom *LikDom) BuildTd(attr ...string) Domer {
	return dom.BuildItem("td", attr...)
}

func (dom *LikDom) BuildTdClass(class string, attr ...string) Domer {
	return dom.BuildTagClass("td", class, attr...)
}

func (dom *LikDom) BuildTrTd(attr ...string) Domer {
	return dom.BuildItem("tr").BuildItem("td", attr...)
}

func (dom *LikDom) BuildTrTdClass(class string, attr ...string) Domer {
	td := dom.BuildTrTd(attr...)
	td.SetAttr("class", class)
	return td
}

func (dom *LikDom) BuildDiv(attr ...string) Domer {
	return dom.BuildItem("div", attr...)
}

func (dom *LikDom) BuildDivClass(class string, attr ...string) Domer {
	return dom.BuildTagClass("div", class, attr...)
}

func (dom *LikDom) BuildTagClass(tag string, class string, attr ...string) Domer {
	item := dom.BuildItem(tag, attr...)
	if class != "" {
		item.SetAttr("class", class)
	}
	return item
}

func (dom *LikDom) BuildTextClassProc(text string, class string, proc string, attr ...string) Domer {
	item := dom.BuildItem("a", attr...)
	if text != "" {
		item.BuildString(text)
	}
	if class != "" {
		item.SetAttr("class", class)
	}
	if proc != "" {
		item.SetAttr("href", "#")
		item.SetAttr("onclick", proc)
	}
	return item
}
