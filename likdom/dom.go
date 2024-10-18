package likdom

import (
	"strings"

	"github.com/massarakhsh/lik"
)

type likDom struct {
	tag      string
	isUnpair bool
	attrs    lik.DItemSet
	text     string
	content  []*likDom
}

type Domer interface {
	GetSelf() *likDom
	GetTag() string
	IsAttr(key string) bool
	GetAttr(key string) string
	SetAttr(attrs ...interface{})
	ToString() string
	ToPrefixString(prefix string) string
	GetDataCount() int
	GetDataItem(num int) (Domer, bool)
	GetDataTag(tag string) (Domer, bool)
	AppendItem(items ...Domer)
	BuildSpace() Domer
	BuildString(text ...string)
	BuildItem(tag string, attrs ...interface{}) Domer
	BuildItemClass(tag string, class string, attrs ...interface{}) Domer
	BuildUnpairItem(tag string, attrs ...interface{}) Domer
	BuildTable(attrs ...interface{}) Domer
	BuildTableClass(class string, attrs ...interface{}) Domer
	BuildTr(attrs ...interface{}) Domer
	BuildTrClass(class string, attrs ...interface{}) Domer
	BuildTd(attrs ...interface{}) Domer
	BuildTdClass(class string, attrs ...interface{}) Domer
	BuildTrTd(attrs ...interface{}) Domer
	BuildTrTdClass(class string, attrs ...interface{}) Domer
	BuildDiv(attrs ...interface{}) Domer
	BuildDivClass(class string, attrs ...interface{}) Domer
}

func BuildSpace() Domer {
	dom := &likDom{}
	return dom
}

func BuildString(text string) Domer {
	dom := &likDom{}
	dom.text = text
	return dom
}

func BuildItem(tag string, attrs ...interface{}) Domer {
	dom := &likDom{tag: tag}
	dom.SetAttr(attrs...)
	return dom
}

func BuildUnpairItem(tag string, attrs ...interface{}) Domer {
	dom := &likDom{tag: tag, isUnpair: true}
	dom.SetAttr(attrs...)
	return dom
}

func BuildItemClass(tag string, class string, attrs ...interface{}) Domer {
	dom := &likDom{tag: tag}
	dom.SetAttr(attrs...)
	if class != "" {
		dom.SetAttr("class", class)
	}
	return dom
}

func BuildPageHtml() Domer {
	item := BuildItem("html")
	item.BuildItem("head")
	item.BuildItem("body")
	return item
}

func BuildDiv(attrs ...interface{}) Domer {
	return BuildItem("div", attrs...)
}
func BuildDivClass(class string, attrs ...interface{}) Domer {
	return BuildItemClass("div", class, attrs...)
}
func BuildDivClassId(class string, id string, attrs ...interface{}) Domer {
	item := BuildDivClass(class, attrs...)
	if id != "" {
		item.SetAttr("id", id)
	}
	return item
}

func BuildTable(attrs ...interface{}) Domer {
	return BuildItem("table", attrs...)
}
func BuildTableClass(class string, attrs ...interface{}) Domer {
	return BuildItemClass("table", class, attrs...)
}
func BuildTableClassId(class string, id string, attrs ...interface{}) Domer {
	item := BuildTableClass(class, attrs...)
	if id != "" {
		item.SetAttr("id", id)
	}
	return item
}

func (dom *likDom) GetSelf() *likDom {
	return dom
}

func (dom *likDom) GetTag() string {
	return dom.tag
}

func (dom *likDom) IsAttr(key string) bool {
	return dom.attrs.IsItem(key)
}

func (dom *likDom) GetAttr(key string) string {
	return dom.attrs.GetString(key)
}

func (dom *likDom) SetAttr(attrs ...interface{}) {
	dom.attrs.SetValues(attrs...)
}

func (dom *likDom) AppendItem(items ...Domer) {
	for _, item := range items {
		if item != nil {
			dom.content = append(dom.content, item.GetSelf())
		}
	}
}

func (dom *likDom) BuildItem(tag string, attrs ...interface{}) Domer {
	item := BuildItem(tag, attrs...)
	dom.AppendItem(item)
	return item
}

func (dom *likDom) BuildItemClass(tag string, class string, attrs ...interface{}) Domer {
	item := BuildItemClass(tag, class, attrs...)
	dom.AppendItem(item)
	return item
}

func (dom *likDom) BuildUnpairItem(tag string, attrs ...interface{}) Domer {
	item := BuildUnpairItem(tag, attrs...)
	dom.AppendItem(item)
	return item
}

func (dom *likDom) BuildString(text ...string) {
	for _, txt := range text {
		item := BuildString(txt)
		dom.AppendItem(item)
	}
}

func (dom *likDom) BuildSpace() Domer {
	item := BuildSpace()
	dom.AppendItem(item)
	return item
}

func (dom *likDom) BuildTable(attrs ...interface{}) Domer {
	return dom.BuildItem("table", attrs...)
}

func (dom *likDom) BuildTableClass(class string, attrs ...interface{}) Domer {
	return dom.BuildTagClass("table", class, attrs...)
}

func (dom *likDom) BuildTr(attrs ...interface{}) Domer {
	return dom.BuildItem("tr", attrs...)
}

func (dom *likDom) BuildTrClass(class string, attrs ...interface{}) Domer {
	return dom.BuildTagClass("tr", class, attrs...)
}

func (dom *likDom) BuildTd(attrs ...interface{}) Domer {
	return dom.BuildItem("td", attrs...)
}

func (dom *likDom) BuildTdClass(class string, attrs ...interface{}) Domer {
	return dom.BuildTagClass("td", class, attrs...)
}

func (dom *likDom) BuildTrTd(attrs ...interface{}) Domer {
	return dom.BuildItem("tr").BuildItem("td", attrs...)
}

func (dom *likDom) BuildTrTdClass(class string, attrs ...interface{}) Domer {
	td := dom.BuildTrTd(attrs...)
	td.SetAttr("class", class)
	return td
}

func (dom *likDom) BuildDiv(attrs ...interface{}) Domer {
	return dom.BuildItem("div", attrs...)
}

func (dom *likDom) BuildDivClass(class string, attrs ...interface{}) Domer {
	return dom.BuildTagClass("div", class, attrs...)
}

func (dom *likDom) BuildTagClass(tag string, class string, attrs ...interface{}) Domer {
	item := dom.BuildItem(tag, attrs...)
	if class != "" {
		item.SetAttr("class", class)
	}
	return item
}

func (dom *likDom) GetDataCount() int {
	return len(dom.content)
}

func (dom *likDom) GetDataItem(num int) (Domer, bool) {
	if num >= 0 && num < len(dom.content) {
		elm := dom.content[num]
		return elm, true
	}
	return nil, false
}

func (dom *likDom) GetDataTag(tag string) (Domer, bool) {
	for n := 0; n < len(dom.content); n++ {
		item := dom.content[n]
		if item != nil && item.GetTag() == tag {
			return item, true
		}
	}
	return nil, false
}

func (dom *likDom) ToString() string {
	code := ""
	if dom.tag != "" {
		if dom.tag == "html" {
			code += "<!DOCTYPE html>\n"
		}
		tag := dom.tag
		tags := tag
		if dom.text != "" {
			tags += " " + dom.text
		}
		for _, set := range dom.attrs.Values() {
			tags += " "
			tags += set.Key
			if val := set.Val.ToString(); val != "" {
				qval := strings.Replace(val, "\\", "\\\\", -1)
				qval = strings.Replace(qval, "\"", "\\\"", -1)
				qval = strings.Replace(qval, "\n", "\\n", -1)
				tags += "=\"" + qval + "\""
			}
		}
		code += "<" + tags + ">"
		if !dom.isUnpair {
			for _, item := range dom.content {
				if item != nil {
					code += item.ToString()
				}
			}
			code += "</" + tag + ">"
		}
	} else if dom.text != "" {
		code += dom.text
	} else {
		for _, item := range dom.content {
			if item != nil {
				code += item.ToString()
			}
		}
	}
	return code
}

func (dom *likDom) ToPrefixString(prefix string) string {
	code := ""
	if dom.tag != "" {
		if prefix == "" && dom.tag == "html" {
			code += "<!DOCTYPE html>\n"
		}
		tag := dom.tag
		tags := tag
		if dom.text != "" {
			tags += " " + dom.text
		}
		for _, set := range dom.attrs.Values() {
			tags += " " + set.Key
			if val := set.Val.ToString(); val != "" {
				qval := strings.Replace(val, "\\", "\\\\", -1)
				qval = strings.Replace(qval, "\"", "\\\"", -1)
				qval = strings.Replace(qval, "\n", "\\n", -1)
				tags += "=\"" + qval + "\""
			}
		}
		code += prefix + "<" + tags + ">"
		if dom.isUnpair {
			code += "\n"
		} else if strings.ToLower(tag) == "textarea" {
			for _, item := range dom.content {
				if item != nil {
					code += item.ToString()
				}
			}
			code += "</" + tag + ">\n"
		} else {
			code += "\n"
			for _, item := range dom.content {
				if item != nil {
					code += item.ToPrefixString(prefix + "    ")
				}
			}
			code += prefix + "</" + tag + ">\n"
		}
	} else if dom.text != "" {
		code += prefix + dom.text + "\n"
	} else {
		for _, item := range dom.content {
			if item != nil {
				code += item.ToPrefixString(prefix)
			}
		}
	}
	return code
}
