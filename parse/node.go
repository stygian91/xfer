package parse

type NodeKind int

const (
	NILKIND = iota
	PROGRAM
	STRUCT
	FIELD
	IDENT
	TYPENAME
	FUNC_CALL
	BOOL
	INT
	FLOAT
	STRING
	VALIDATION
	BOOLTYPE
	STRINGTYPE
	INTTYPE
	FLOATTYPE
	CUSTOMTYPE
)

type Node struct {
	Kind NodeKind
	Value interface{}

	Children []Node
}
