package php

import (
	"fmt"
	"strings"

	"github.com/stygian91/xfer/gen"
	p "github.com/stygian91/xfer/parse"
)

func Struct(node p.Node) (string, error) {
	builder := strings.Builder{}
	if node.Kind != p.STRUCT {
		return "", fmt.Errorf("php.Struct(): tried to generate struct code but non-struct node was passed.")
	}

	name, exists := gen.GetIdent(node)
	if !exists {
		return "", fmt.Errorf("php.Struct(): struct does not have a name")
	}

	builder.WriteString(fmt.Sprintf("class %s {", name))

	if len(node.Children) <= 1 {
		goto END
	}

	for _, fieldNode := range node.Children[1:] {
		fieldName, exists := gen.GetIdent(fieldNode)
		if !exists {
			return "", fmt.Errorf("php.Struct(): struct field does not have a name")
		}

		if len(fieldNode.Children) < 2 {
			return "", fmt.Errorf("php.Struct(): struct field does not have a type. Node: %+v", fieldNode)
		}

		fieldType, err := TypenameToStr(fieldNode.Children[1])
		if err != nil {
			return "", err
		}

		builder.WriteString(fmt.Sprintf("\npublic %s $%s;", fieldType, fieldName))
	}

END:
	builder.WriteString("\n}")

	return builder.String(), nil
}
