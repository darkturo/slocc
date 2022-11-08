package slocc

import "strings"

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
func findMultilineEnding(config Config, line string) bool {
	return strings.Contains(line, config.MultiLineEndCommentMark)
}

// isMultilineComment returns true if the line starts a multi-line comment
func isMultilineComment(config Config, line string) bool {
	return strings.Contains(line, config.MultiLineBeginCommentMark)
}

// isSingleLineComment returns true if the line is a single-line comment
func isSingleLineComment(config Config, line string) bool {
	for _, marker := range config.SingleLineCommentMarker {
		if strings.HasPrefix(line, marker) {
			return true
		}
	}
	return false
}
