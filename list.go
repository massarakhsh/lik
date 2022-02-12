package lik

//	Интерфейс динамических списков
type Lister interface {
	Itemer
	Count() int
	GetItem(idx int) Itemer
	GetBool(idx int) bool
	GetInt(idx int) int64
	GetFloat(idx int) float64
	GetString(idx int) string
	GetList(idx int) Lister
	GetSet(idx int) Seter
	GetIDB(idx int) IDB
	AddItems(vals ...interface{})
	AddItemers(vals []Itemer)
	InsertItem(val interface{}, idx int)
	AddItemSet(vals ...interface{}) Seter
	SetItem(val interface{}, idx int) bool
	DelItem(idx int) bool
	SwapItem(pos1 int, pos2 int)
	ToCsv(dlm string) string
	ToXml() string
	Values() []Itemer
	Self() *DItemList
}

func BuildList(vals ...interface{}) Lister {
	list := &DItemList{}
	list.AddItems(vals...)
	return list
}

func BuildFromArray(vals []Seter) Lister {
	list := &DItemList{}
	for _, val := range vals {
		list.AddItems(val)
	}
	return list
}

func (it *DItemList) ToItem() Itemer {
	return it
}

func (it *DItemList) clone() Itemer {
	cpy := BuildList()
	for _, elm := range it.Val {
		cpy.AddItems(elm.Clone())
	}
	return cpy
}

func (it *DItemList) serialize() string {
	var text = "["
	for n, val := range it.Val {
		if n > 0 {
			text += ","
		}
		text += val.Serialize()
	}
	text += "]"
	return text
}

func (it *DItemList) format(prefix string) string {
	var text = "["
	if it.Count() > 0 {
		for n, val := range it.Val {
			if n > 0 {
				text += ","
			}
			text += "\n" + prefix + "  "
			text += val.Format(prefix + "  ")
		}
		text += "\n" + prefix
	}
	text += "]"
	return text
}

func (it *DItemList) Count() int {
	return len(it.Val)
}

func (it *DItemList) GetItem(idx int) Itemer {
	if idx >= 0 && idx < len(it.Val) {
		return it.Val[idx]
	}
	return nil
}

func (it *DItemList) GetBool(idx int) bool {
	if item := it.GetItem(idx); item != nil {
		return item.ToBool()
	}
	return false
}

func (it *DItemList) GetInt(idx int) int64 {
	if item := it.GetItem(idx); item != nil {
		return item.ToInt()
	}
	return 0
}

func (it *DItemList) GetFloat(idx int) float64 {
	if item := it.GetItem(idx); item != nil {
		return item.ToFloat()
	}
	return 0
}

func (it *DItemList) GetString(idx int) string {
	if item := it.GetItem(idx); item != nil {
		return item.ToString()
	}
	return ""
}

func (it *DItemList) GetIDB(idx int) IDB {
	return IDB(it.GetInt(idx))
}

func (it *DItemList) GetList(idx int) Lister {
	if item := it.GetItem(idx); item != nil {
		return item.ToList()
	}
	return nil
}

func (it *DItemList) GetSet(idx int) Seter {
	if item := it.GetItem(idx); item != nil {
		return item.ToSet()
	}
	return nil
}

func (it *DItemList) SetItem(val interface{}, idx int) bool {
	modify := false
	if idx >= 0 && idx < len(it.Val) {
		if val != nil {
			valnew := BuildItem(val)
			if valold := it.Val[idx]; valold != nil {
				if valnew.IsSet() != valold.IsSet() ||
					valnew.IsList() != valold.IsList() ||
					valnew.ToString() != valold.ToString() {
					modify = true
				}
			} else {
				modify = true
			}
			it.Val[idx] = valnew
		} else if idx >= 0 && idx < len(it.Val) {
			list := []Itemer{}
			for ni := 0; ni < len(it.Val); ni++ {
				if ni != idx {
					list = append(list, it.Val[ni])
				}
			}
			it.Val = list
			modify = true
		}
	}
	return modify
}

func (it *DItemList) DelItem(idx int) bool {
	return it.SetItem(nil, idx)
}

func (it *DItemList) AddItems(vals ...interface{}) {
	for _, val := range vals {
		if val != nil {
			it.Val = append(it.Val, BuildItem(val))
		}
	}
}

func (it *DItemList) AddItemers(vals []Itemer) {
	for _, val := range vals {
		if val != nil {
			it.Val = append(it.Val, BuildItem(val))
		}
	}
}

func (it *DItemList) InsertItem(val interface{}, idx int) {
	if val != nil {
		item := BuildItem(val)
		list := []Itemer{}
		for ni := 0; ni < len(it.Val); ni++ {
			if ni == idx {
				list = append(list, item)
				item = nil
			}
			list = append(list, it.Val[ni])
		}
		if item != nil {
			list = append(list, item)
		}
		it.Val = list
	}
}

func (it *DItemList) AddItemSet(vals ...interface{}) Seter {
	item := BuildSet(vals...)
	it.AddItems(item)
	return item
}

func (it *DItemList) SwapItem(pos1 int, pos2 int) {
	if pos1 != pos2 && pos1 < it.Count() && pos2 < it.Count() {
		itm := it.Val[pos1]
		it.Val[pos1] = it.Val[pos2]
		it.Val[pos2] = itm
	}
}

func (it *DItemList) Values() []Itemer {
	return it.Val
}

func (it *DItemList) Self() *DItemList {
	return it
}

func (it *DItemList) ToCsv(dlm string) string {
	dump := ""
	ml := it.Count()
	for nl := 0; nl < ml; nl++ {
		if line := it.GetList(nl); line != nil {
			me := line.Count()
			for ne := 0; ne < me; ne++ {
				dump += line.GetString(ne)
				if ne+1 < me {
					dump += dlm
				}
			}
			dump += "\r\n"
		}
	}
	return dump
}

func (it *DItemList) ToXml() string {
	dump := XML_ListToString("", it)
	return dump
}
