package liktable

import (
	"fmt"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likdom"
	"math/rand"
)

type Table struct {
	options		lik.Seter
	columns    	lik.Lister
}

type Tabler interface {
}

func New(opts ...interface{}) *Table {
	it := &Table{}
	it.options = lik.BuildSet(opts...)
	return it
}

func (it *Table) Initialize(path string) likdom.Domer {
	div := likdom.BuildDivClass("grid")
	id := fmt.Sprintf("id_%d", 100000+rand.Int31n(900000))
	code := likdom.BuildTableClassId("grid", id, "path", path, "redraw=grid_redraw")
	div.AppendItem(code)
	return div
}

func (it *Table) Show() lik.Seter {
	grid := it.options.Clone().(lik.Seter)
	if !grid.IsItem("serverSide") {
		grid.SetItem(true, "serverSide")
	}
	if !grid.IsItem("processing") {
		grid.SetItem(grid.GetBool("serverSide"), "processing")
	}
	if !grid.IsItem("info") {
		grid.SetItem(false, "info")
	}
	if !grid.IsItem("paging") {
		grid.SetItem(false, "paging")
	}
	if !grid.IsItem("lengthChange") {
		grid.SetItem(false, "lengthChange")
	}
	if !grid.IsItem("pageLength") {
		grid.SetItem(15, "pageLength")
	}
	if !grid.IsItem("searching") {
		grid.SetItem(false, "searching")
	}
	if !grid.IsItem("select/style") {
		grid.SetItem("single", "select/style")
	}
	columns := lik.BuildList()
	if it.columns != nil {
		for nc := 0; nc < it.columns.Count(); nc++ {
			if col := it.columns.GetSet(nc); col != nil {
				columns.AddItems(col)
			}
		}
	}
	grid.SetItem(columns, "columns")
	grid.SetItem(it.showLanguage(), "language")
	return grid
}

func (it *Table) AddColumn(opts ...interface{}) {
	if it.columns == nil {
		it.columns = lik.BuildList()
	}
	it.columns.AddItemSet(opts...)
}

func (it *Table) showLanguage() lik.Seter {
	data := lik.BuildSet()
	data.SetItem("Поиск", "search")
	data.SetItem("Таблица пуста", "emptyTable")
	data.SetItem("Строки от _START_ до _END_, всего _TOTAL_", "info")
	data.SetItem("Загрузка ...", "loadingRecords")
	data.SetItem("Обработка ...", "processing")
	data.SetItem("Нет строк в таблице", "infoEmpty")
	data.SetItem("В начало", "paginate/first")
	data.SetItem("Назад", "paginate/previos")
	data.SetItem("Вперёд", "paginate/next")
	data.SetItem("В конец", "paginate/last")
	return data
}
