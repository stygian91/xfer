package parse

type NodeKind int

const (
	NILKIND = iota
	STRUCT
	FIELD
	IDENT
	TYPENAME
	FUNC_CALL
)

type Node struct {
	Kind NodeKind
	Value interface{}

	Children []Node
}
