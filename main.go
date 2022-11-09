package main

import (
	"bufio"
	"fmt"
	"github.com/darkturo/slocc/internal/pkg/filetype"
	"github.com/darkturo/slocc/internal/pkg/slocc"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	files := make([]string, 0, len(os.Args))

	results := make(map[filetype.FileType]uint)
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
		SlocByLanguage map[filetype.FileType]uint
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

func sloc(path string) (uint, filetype.FileType, error) {
	fileType := filetype.Guess(path)

	file, err := os.Open(path)
	if err != nil {
		return 0, fileType, err
	}
	defer file.Close()

	lines, err := slocc.CountLinesOfCode(slocc.Config{
		SingleLineCommentMarker:   []string{"//"},
		MultiLineBeginCommentMark: "/*",
		MultiLineEndCommentMark:   "*/",
	}, bufio.NewReader(file))
	if err != nil {
		return 0, fileType, err
	}

	return lines, fileType, nil
}
