package parse

type NodeKind int

const (
	NILKIND = iota
	STRUCT
	FIELD
	IDENT
	TYPENAME
)

type Node struct {
	Kind NodeKind
	Value interface{}

	Children []Node
}
