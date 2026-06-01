package testutil

import "github.com/yuin/goldmark/ast"

// Finds a node of a specific type in the Markdown AST tree.
func FindNode[T ast.Node](tree ast.Node) T {
	var found T
	ast.Walk(tree, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
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
