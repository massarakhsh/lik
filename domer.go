package lik

type itDomer struct {
	DItemSet
}

type Domer interface {
	Seter
	String() string
}

func BuildDomer() Domer {

}

func (it *itDomer) String() string {
	return "<html></html>"
}
