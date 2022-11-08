package slocc

import (
	"bufio"
	"io"
	"strings"
)

type Config struct {
	SingleLineCommentMarker   []string
	MultiLineBeginCommentMark string
	MultiLineEndCommentMark   string
}

type multiLineCommentContext struct {
	level int
}

func (m *multiLineCommentContext) enterContext() {
	m.level++
}

func (m *multiLineCommentContext) exitContext() {
	if m.level > 0 {
		m.level--
	}
}

func (m multiLineCommentContext) isInContext() bool {
	return m.level > 0
}

func isSingleLineComment(config Config, line string) bool {
	for _, mark := range config.SingleLineCommentMarker {
		if len(line) > len(mark) && line[0:len(mark)] == mark {
			return true
		}
	}
	return false
}

func startsWithMultilineBeginCommentMark(config Config, line string) bool {
	mark := config.MultiLineBeginCommentMark
	return len(line) > len(mark) && line[0:len(mark)] == mark
}

func findMultilineEndCommentMarkInThisLine(config Config, line string) bool {
	return strings.Contains(line, config.MultiLineEndCommentMark)
}

func CountLinesOfCode(config Config, file *bufio.Reader) (uint, error) {
	var counter uint
	var mlcc multiLineCommentContext
	var line []byte
	var isPrefix bool
	var err error
readLineLoop:
	for {
		var loc string

	collectStringLoop:
		for {
			line, isPrefix, err = file.ReadLine()
			if err != nil {
				if err == io.EOF {
					break readLineLoop
				}
				return 0, nil
			}
			loc += string(line)

			if !isPrefix {
				break collectStringLoop
			}

		}

		if !mlcc.isInContext() {
			if isSingleLineComment(config, loc) {
				continue
			}

			if startsWithMultilineBeginCommentMark(config, loc) {
				mlcc.enterContext()
				if findMultilineEndCommentMarkInThisLine(config, loc) {
					mlcc.exitContext()
				}
				continue
			}

			counter++
		} else {
			if findMultilineEndCommentMarkInThisLine(config, loc) {
				mlcc.exitContext()
			}
		}
	}
	return counter, nil
}
