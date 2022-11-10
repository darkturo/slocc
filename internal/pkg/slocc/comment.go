package slocc

import (
	"github.com/darkturo/slocc/internal/pkg/config"
	"strings"
)

// represents a multi-line comment, which can be nested
type multiLineCommentContext struct {
	level int
}

// enterContext enters a multi-line comment context
func (m *multiLineCommentContext) enterContext() {
	m.level++
}

// exitContext exits a multi-line comment context
func (m *multiLineCommentContext) exitContext() {
	if m.level > 0 {
		m.level--
	}
}

// isInContext returns true if the multi-line comment context is active
func (m multiLineCommentContext) isInContext() bool {
	return m.level > 0
}

// findMultilineEnding returns true if the line ends a multi-line comment
func findMultilineEnding(language config.Lang, line string) bool {
	return strings.Contains(line, language.MultiLineEndCommentMark)
}

// isMultilineComment returns true if the line starts a multi-line comment
func isMultilineComment(config config.Lang, line string) bool {
	return strings.Contains(line, config.MultiLineBeginCommentMark)
}

// isSingleLineComment returns true if the line is a single-line comment
func isSingleLineComment(config config.Lang, line string) bool {
	for _, marker := range config.SingleLineCommentMarker {
		if strings.HasPrefix(line, marker) {
			return true
		}
	}
	return false
}
