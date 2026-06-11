package testutil

import (
	"reflect"

	"github.com/yuin/goldmark/ast"
)

// Finds a child node of a specific type in the Markdown AST node.
func FindChildNode[T ast.Node](node ast.Node) T {
	var found T
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if t, ok := n.(T); ok {
			found = t
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})
	return found
}

// Finds a node of a specific reflected type in the Markdown AST node.
// Prefer [FindChildNode] for most cases unless you need to test dynamic types.
func FindChildNodeReflect(node ast.Node, ty reflect.Type) ast.Node {
	nodeType := reflect.TypeOf((*ast.Node)(nil)).Elem()
	if !ty.Implements(nodeType) {
		panic("ty " + ty.String() + " does not implement ast.Node")
	}

	var found ast.Node
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if reflect.TypeOf(n) == ty {
			found = n
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})
	return found
}
