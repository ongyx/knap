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
// If the format is invalid, ok is false.
func ParseEmbedSize(s string) (width, height int, ok bool) {
	m := reEmbedSize.FindStringSubmatchIndex(s)
	if m == nil {
		return 0, 0, false
	}

	wstart := m[2]
	wend := m[3]
	hstart := m[4]
	hend := m[5]

	w := s[wstart:wend]
	// SAFETY: w must only consist of digits as per reEmbedSize.
	width = util.Must(strconv.Atoi(w))

	if hstart > 0 && hend > 0 {
		h := s[hstart:hend]
		// SAFETY: h must only consist of digits as per reEmbedSize.
		height = util.Must(strconv.Atoi(h))
	}

	return width, height, true
}

// Returns the child text nodes of the given node without any Markdown formatting.
func NodeChildrenToText(node ast.Node, source []byte) string {
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
		panic("error returned from walking AST when there shouldn't be")
	}

	return buf.String()
}
