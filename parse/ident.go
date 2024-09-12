package parse

type IdentValue struct {
	Name string
}

func NewIdent(name string) Node {
	return Node{Kind: IDENT, Value: IdentValue{Name: name}}
}
