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

var excluded = []string{
	".git/",
	".idea/",
}

func main() {
	files := make([]string, 0, len(os.Args))

	results := make(map[filetype.FileType]uint)
	for _, f := range os.Args[1:] {
		err := filepath.Walk(f, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() && !isExcluded(path) {
				files = append(files, path)
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

		if lang == filetype.Binary {
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

func isExcluded(path string) bool {
	for _, pattern := range excluded {
		if strings.HasPrefix(path, pattern) {
			return true
		}
	}

	return false
}

func sloc(path string) (uint, filetype.FileType, error) {
	fileType := filetype.Guess(path)

	file, err := os.Open(path)
	if err != nil {
		return 0, fileType, err
	}
	defer file.Close()

	lines, err := slocc.CountLinesOfCode(fileType, bufio.NewReader(file))
	if err != nil {
		return 0, fileType, err
	}

	return lines, fileType, nil
}
