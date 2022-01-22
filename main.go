package psloc

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	results := make(map[string]uint)
	for _, path := range os.Args {
		linesOfCode, format, err := sloc(path)
		if err != nil {
			fmt.Printf("invalid file %s", path)
			continue
		}
		results[format] += linesOfCode
	}
	fmt.Printf("%v", results)
}

type slocConfig struct {
	singleLineCommentMarker   []string
	multiLineBeginCommentMark string
	multiLineEndCommentMark   string
}

func sloc(path string) (uint, string, error) {
	fileType, err := guessFileType(path)
	if err != nil {
		return 0, "", err
	}

	file, err := os.Open(path)
	if err != nil {
	}
	defer file.Close()

	lines, err := countLinesOfCode(slocConfig{}, file)
	if err != nil {
		return 0, fileType, err
	}

	return lines, fileType, nil
}

func guessFileType(path string) (string, error) {
	return "go", nil
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

func countLinesOfCode(config slocConfig, file *os.File) (uint, error) {
	reader := bufio.NewReader(file)

	var counter uint
	var mlcc multiLineCommentContext
	var line []byte
	var isPrefix bool
	var err error
	for {
		var loc string
		for {
			line, isPrefix, err = reader.ReadLine()
			if err != nil {
				return 0, nil
			}
			loc += string(line)

			if !isPrefix {
				break
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
}

func isSingleLineComment(config slocConfig, line string) bool {
	for _, mark := range config.singleLineCommentMarker {
		if line[0:len(line)] == mark {
			return true
		}
	}
	return false
}

func startsWithMultilineBeginCommentMark(config slocConfig, line string) bool {
	return line[0:len(line)] == config.multiLineBeginCommentMark
}

func findMultilineEndCommentMarkInThisLine(config slocConfig, line string) bool {
	return strings.Contains(line, config.multiLineEndCommentMark)
}