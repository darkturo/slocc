package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	files := make([]string, 0, len(os.Args))
	results := make(map[string]uint)
	for _, f := range os.Args[1:] {
		err := filepath.Walk(f, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() && !isExcluded(path) {
				binary, err := looksLikeBinary(path)
				if err != nil {
					return err
				}
				if !binary {
					files = append(files, path)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("error inspecting %s\n", f)
			continue
		}
	}

	totalSLOC := uint(0)
	for _, path := range files {
		loc, lang, err := sloc(path)
		if err != nil {
			fmt.Printf("invalid file %s\n", path)
			continue
		}
		results[lang] += loc
		totalSLOC += loc
	}

	tmpl, err := template.New("slocc output").
		Parse(`
SLOC	SLOC-by-Language (Sorted)
{{.Sloc}}	{{range $key, $value := .SlocByLanguage}} {{$key}}={{$value}},{{end}}
`)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, struct {
		Sloc           uint
		SlocByLanguage map[string]uint
	}{
		Sloc:           totalSLOC,
		SlocByLanguage: results,
	})
	if err != nil {
		panic(err)
	}
}

var excluded = []string{
	".git/",
	".idea/",
}

func isExcluded(path string) bool {
	for _, pattern := range excluded {
		if strings.HasPrefix(path, pattern) {
			return true
		}
	}

	return false
}

func looksLikeBinary(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil {
		return false, err
	}

	for i := 0; i < n; i++ {
		if buffer[i] == 0 {
			return true, nil
		}
	}
	return false, nil
}

type slocConfig struct {
	singleLineCommentMarker   []string
	multiLineBeginCommentMark string
	multiLineEndCommentMark   string
}

func sloc(path string) (uint, string, error) {
	fileType := guessFileType(path)

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

var extensions = map[string]string{
	".go": "go",
}

func guessFileType(path string) string {
	if fileType, ok := extensions[filepath.Ext(path)]; ok {
		return fileType
	}
	return "other"
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
