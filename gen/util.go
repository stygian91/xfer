package gen

import p "github.com/stygian91/xfer/parse"

func GetIdent(node p.Node) (string, bool) {
	if len(node.Children) == 0 {
		return "", false
	}

	first := node.Children[0]
	if first.Kind != p.IDENT {
		return "", false
	}

	value, valid := first.Value.(p.IdentValue)
	if !valid {
		return "", false
	}

	return value.Name, true
}
