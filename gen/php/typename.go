package php

import (
	"fmt"

	p "github.com/stygian91/xfer/parse"
)

func TypenameToStr(node p.Node) (string, error) {
	if node.Kind != p.TYPENAME {
		return "", fmt.Errorf("php.TypenameToStr(): not a Typename")
	}

	if len(node.Children) < 1 {
		return "", fmt.Errorf("php.TypenameToStr(): no Typename child")
	}

	switch node.Children[0].Kind {
	case p.BOOLTYPE:
		return "bool", nil
	case p.INTTYPE:
		return "int", nil
	case p.FLOATTYPE:
		return "float", nil
	case p.STRINGTYPE:
		return "string", nil
	case p.CUSTOMTYPE:
		return node.Children[0].Value.(string), nil
	default:
		return "", fmt.Errorf("php.TypenameToStr(): invalid typename")
	}
}
