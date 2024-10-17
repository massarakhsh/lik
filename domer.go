package lik

import "strings"

type itDom struct {
	tag      string
	attrs    DItemSet
	isUnpair bool
	content  DItemList
}

type Domer interface {
	ToString() string
	GetTag() string
	IsAttr(key string) bool
	GetAttr(key string) string
	SetAttrs(attrs ...interface{})
	// ToPrefixString(prefix string) string
	// GetDataCount() int
	// GetDataItem(num int) (Domer, bool)
	// GetDataTag(tag string) (Domer, bool)
	// AppendItem(items ...Domer)
	// BuildSpace() Domer
	// BuildString(text ...string)
	// BuildItem(tag string, attr ...string) Domer
	// BuildItemClass(tag string, class string, attr ...string) Domer
	// BuildUnpairItem(tag string, attr ...string) Domer
	// BuildTable(attr ...string) Domer
	// BuildTableClass(class string, attr ...string) Domer
	// BuildTr(attr ...string) Domer
	// BuildTrClass(class string, attr ...string) Domer
	// BuildTd(attr ...string) Domer
	// BuildTdClass(class string, attr ...string) Domer
	// BuildTrTd(attr ...string) Domer
	// BuildTrTdClass(class string, attr ...string) Domer
	// BuildDiv(attr ...string) Domer
	// BuildDivClass(class string, attr ...string) Domer
	// BuildTextClassProc(text string, class string, proc string, attr ...string) Domer
}

func BuildDomer(tag string, attrs ...interface{}) Domer {
	it := &itDom{tag: tag}
	it.attrs.SetValues(attrs...)
	return it
}

func (it *itDom) GetTag() string {
	return it.tag
}

func (it *itDom) IsAttr(key string) bool {
	return it.attrs.IsItem(key)
}

func (it *itDom) GetAttr(key string) string {
	return it.attrs.GetString(key)
}

func (it *itDom) SetAttrs(attrs ...interface{}) {
	it.attrs.SetValues(attrs...)
}

func (it *itDom) ToString() string {
	return it.toString()
}

func (it *itDom) toString() string {
	code := ""
	if tag := it.tag; tag != "" {
		if it.tag == "html" {
			code += "<!DOCTYPE html>\n"
		}
		tags := tag
		for _, set := range it.attrs.Values() {
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
		// if !it.isUnpair {
		// 	for _, item := range it.content.Values() {
		// 		if item != nil {
		// 			code += item.ToString()
		// 		}
		// 	}
		// 	code += "</" + tag + ">"
		// }
		// } else if dom.Text != "" {
		// 	code += dom.Text
		// } else {
		// 	for _, item := range dom.Data {
		// 		if item != nil {
		// 			code += item.ToString()
		// 		}
		// 	}
	}
	return code
}

// func (dom *LikDom) ToPrefixString(prefix string) string {
// 	code := ""
// 	if dom.Tag != "" {
// 		if prefix == "" && dom.Tag == "html" {
// 			code += "<!DOCTYPE html>\n"
// 		}
// 		tag := dom.Tag
// 		tags := tag
// 		if dom.Text != "" {
// 			tags += " " + dom.Text
// 		}
// 		for key, val := range dom.Attr {
// 			tags += " " + key
// 			if val != "" {
// 				qval := strings.Replace(val, "\\", "\\\\", -1)
// 				qval = strings.Replace(qval, "\"", "\\\"", -1)
// 				qval = strings.Replace(qval, "\n", "\\n", -1)
// 				tags += "=\"" + qval + "\""
// 			}
// 		}
// 		code += prefix + "<" + tags + ">"
// 		if dom.Unpair {
// 			code += "\n"
// 		} else if strings.ToLower(tag) == "textarea" {
// 			for _, item := range dom.Data {
// 				if item != nil {
// 					code += item.ToString()
// 				}
// 			}
// 			code += "</" + tag + ">\n"
// 		} else {
// 			code += "\n"
// 			for _, item := range dom.Data {
// 				if item != nil {
// 					code += item.ToPrefixString(prefix + "    ")
// 				}
// 			}
// 			code += prefix + "</" + tag + ">\n"
// 		}
// 	} else if dom.Text != "" {
// 		code += prefix + dom.Text + "\n"
// 	} else {
// 		for _, item := range dom.Data {
// 			if item != nil {
// 				code += item.ToPrefixString(prefix)
// 			}
// 		}
// 	}
// 	return code
// }
