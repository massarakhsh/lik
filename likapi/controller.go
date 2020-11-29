package likapi

import (
	"fmt"
	"github.com/massarakhsh/lik/likdom"
	"math/rand"
)

type DataControl struct {
	Index  string
	Status string
	Mode   string
	IfShow		func(drive DataDriver, style string) likdom.Domer
	IfMarshal	func(drive DataDriver, status string)
	IfExecute	func(drive DataDriver, path []string)
}

type Controller interface {
	Init(status string)
	GetIndex() string
	GetStatus() string
	SetStatus(index string)
	GetMode() string
	Show(drive DataDriver, style string) likdom.Domer
	Marshal(drive DataDriver, status string)
	Execute(drive DataDriver, path []string)
}

func (it *DataControl) Init(status string) {
	if it.Index == "" {
		it.Index = fmt.Sprintf("%d", 100000000+rand.Intn(900000000))
	}
	if status != "" {
		it.SetStatus(status)
	}
}

func (it *DataControl) GetIndex() string {
	return it.Index
}

func (it *DataControl) GetStatus() string {
	return it.Status
}

func (it *DataControl) SetStatus(status string) {
	it.Status = status
}

func (it *DataControl) GetMode() string {
	return it.Mode
}

func (it *DataControl) PopCommand(path *[]string) string {
	if path == nil || len(*path) == 0 {
		return ""
	} else {
		cmd := (*path)[0]
		*path = (*path)[1:]
		return cmd
	}
}

func (it *DataControl) Show(drive DataDriver, style string) likdom.Domer {
	if it.IfShow != nil {
		return it.IfShow(drive, style)
	}
	return nil
}

func (it *DataControl) Marshal(drive DataDriver, status string) {
	if it.IfMarshal != nil {
		it.IfMarshal(drive, status)
	}
}

func (it *DataControl) Execute(drive DataDriver, path []string) {
	if it.IfExecute != nil {
		it.IfExecute(drive, path)
	}
}
