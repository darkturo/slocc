package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	files := make([]string, 0, len(os.Args))
	results := make(map[string]uint)
	for _, f := range os.Args[1:] {
		err := filepath.Walk(f, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("error inspecting %s\n", f)
			continue
		}
	}
	fmt.Printf("FILES: %v\n", files)

	for _, path := range files {
		linesOfCode, format, err := sloc(path)
		if err != nil {
			fmt.Printf("invalid file %s\n", path)
			continue
		}
		results[format] += linesOfCode
	}
	fmt.Printf("%v\n", results)
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
		return 0, fileType, err
	}
	defer file.Close()

	lines, err := countLinesOfCode(slocConfig{
		singleLineCommentMarker:   []string{"//"},
		multiLineBeginCommentMark: "/*",
		multiLineEndCommentMark:   "*/",
	}, file)
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
readLineLoop:
	for {
		var loc string

	collectStringLoop:
		for {
			line, isPrefix, err = reader.ReadLine()
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

func isSingleLineComment(config slocConfig, line string) bool {
	for _, mark := range config.singleLineCommentMarker {
		if len(line) > len(mark) && line[0:len(mark)] == mark {
			return true
		}
	}
	return false
}

func startsWithMultilineBeginCommentMark(config slocConfig, line string) bool {
	mark := config.multiLineBeginCommentMark
	return len(line) > len(mark) && line[0:len(mark)] == mark
}

func findMultilineEndCommentMarkInThisLine(config slocConfig, line string) bool {
	return strings.Contains(line, config.multiLineEndCommentMark)
}
