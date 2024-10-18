package main

import (
	"fmt"

	"github.com/massarakhsh/lik"
)

func main() {
	test1()
	test2()
	test3()
	test4()
	test5()
}

type TP1 struct {
	Name   string
	Format string
	Number int
	LogT   bool
	LogF   bool
}

func test1() {
	data := TP1{Name: "Имя", Format: "Форма", Number: 13, LogT: true, LogF: false}
	set := lik.BuildSet(data)
	if data1 := lik.SetToType[TP1](set); data1 == data {
		fmt.Println("Test1: Ok")
	} else {
		fmt.Println("ERROR for Test1")
	}
}

type TP2 struct {
	Infa []string
}

func test2() {
	data := TP2{Infa: []string{"Фамилия", "Имя", "Отчество"}}
	set := lik.BuildSet(data)
	if data1 := lik.SetToType[TP2](set); fmt.Sprint(data1) == fmt.Sprint(data) {
		fmt.Println("Test2: Ok")
	} else {
		fmt.Println("ERROR for Test2")
	}
}

type TP3 struct {
	Infa []TP1
}

func test3() {
	data := TP3{}
	data.Infa = append(data.Infa, TP1{Name: "Имя"})
	data.Infa = append(data.Infa, TP1{Format: "Форм"})
	set := lik.BuildSet(data)
	if data1 := lik.SetToType[TP3](set); fmt.Sprint(data1) == fmt.Sprint(data) {
		fmt.Println("Test3: Ok")
	} else {
		fmt.Println("ERROR for Test3")
	}
}

type TP4 struct {
	Infa map[string]TP1
}

func test4() {
	data := TP4{}
	data.Infa = make(map[string]TP1)
	data.Infa["one"] = TP1{Name: "Имя"}
	data.Infa["two"] = TP1{Format: "Форм"}
	set := lik.BuildSet(data)
	if data1 := lik.SetToType[TP4](set); fmt.Sprint(data1) == fmt.Sprint(data) {
		fmt.Println("Test4: Ok")
	} else {
		fmt.Println("ERROR for Test4")
	}
}

type TP5 struct {
	Beta string
	TP5a
}

type TP5a struct {
	Alpha string
}

func test5() {
	data := TP5{}
	data.Alpha = "Альфа"
	data.Beta = "Бета"
	set := lik.BuildSet(data)
	if data1 := lik.SetToType[TP5](set); fmt.Sprint(data1) == fmt.Sprint(data) {
		fmt.Println("Test5: Ok")
	} else {
		fmt.Println("ERROR for Test5")
	}
}
