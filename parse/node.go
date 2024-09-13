package parse

type NodeKind int

const (
	NILKIND = iota
	STRUCT
	FIELD
	IDENT
	TYPENAME
	FUNC_CALL
	BOOL
	INT
	FLOAT
	STRING
)

type Node struct {
	Kind NodeKind
	Value interface{}

	Children []Node
}
