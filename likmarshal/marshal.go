package likmarshal

import (
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
)

var indexData = 0
var registerData = lik.BuildSet()

func UpdateBase(base likbase.DBaser, tables []string, id string) {
	register := lik.BuildSet()
	for _, table := range tables {
		tbl := lik.BuildSet()
		register.SetItem(tbl, table)
		list := base.GetListAll(table)
		for ne := 0; ne < list.Count(); ne++ {
			if elm := list.GetSet(ne); elm != nil {
				tbl.SetItem(elm, elm.GetString(id))
			}
		}
	}
	Update(register)
}

func Update(register lik.Seter) bool {
	news := false
	news = true
	indexData++
	if register != nil {
		registerData = register.Clone().ToSet()
	} else {
		registerData = lik.BuildSet()
	}
	return news
}

func Answer(index int) lik.Seter {
	indexData++
	answer := lik.BuildSet("index", indexData, "register", registerData)
	return answer
}
