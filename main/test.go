package main

import (
	"fmt"

	"github.com/massarakhsh/lik"
)

type TP struct {
	Name   string
	Format string
	Lst    []string
}

func main() {
	data := TP{Name: "Имя", Format: "Форма"}
	data.Lst = []string{"Раз", "Два"}
	set := lik.BuildSet(data)
	fmt.Println(set.Serialize())
	data1 := lik.SetToStruct[TP](set)
	_ = data1
	fmt.Println("Ok")
}
