package converter

import (
	"bytes"
	"regexp"
	"strconv"

	"github.com/ongyx/knap/internal/util"
	"github.com/yuin/goldmark/ast"
)

// Matches an embed size in the format (w)x(h), where h is optional.
var reEmbedSize = regexp.MustCompile(`^(\d+)(?:x(\d+))?$`)

// Parses an embed size of the format '(w)' or '(w)x(h)', where w and h are integers.
//
// If the format is invalid, width and height will be 0.
func ParseEmbedSize(b []byte) (width, height int) {
	m := reEmbedSize.FindSubmatchIndex(b)

	if len(m) >= 4 {
		w := b[m[2]:m[3]]
		// SAFETY: w must only consist of digits as per reEmbedSize.
		width = util.Must(strconv.Atoi(string(w)))
	}

	if len(m) == 6 {
		h := b[m[4]:m[5]]
		// SAFETY: h must only consist of digits as per reEmbedSize.
		height = util.Must(strconv.Atoi(string(h)))
	}

	return width, height
}

// Returns the child text nodes of the given node without any Markdown formatting.
func NodeChildrenToText(node ast.Node, source []byte) []byte {
	var buf bytes.Buffer

	walker := func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if t, ok := node.(*ast.Text); ok {
			buf.Write(t.Value(source))
		}

		return ast.WalkContinue, nil
	}

	// The error here is guaranteed to be nil if the walker doesn't return any error.
	if err := ast.Walk(node, walker); err != nil {
		panic("erorr returned from walking AST when there shouldn't be")
	}

	return buf.Bytes()
}
