package slocc

import (
	"bufio"
	"github.com/darkturo/slocc/internal/pkg/filetype"
	"io"
)

// CountLinesOfCode counts the lines of code in a file  (excluding comments and empty lines)
func CountLinesOfCode(fileType filetype.FileType, file *bufio.Reader) (uint, error) {
	var counter uint
	var commentContext multiLineCommentContext
	var line []byte
	var isPrefix bool
	var err error

	config := languageConfigurations[fileType]
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

		if !commentContext.isInContext() {
			if isSingleLineComment(config, loc) {
				continue
			}

			if isMultilineComment(config, loc) {
				commentContext.enterContext()
				if findMultilineEnding(config, loc) {
					commentContext.exitContext()
				}
				continue
			}

			counter++
		} else {
			if findMultilineEnding(config, loc) {
				commentContext.exitContext()
			}
		}
	}
	return counter, nil
}
