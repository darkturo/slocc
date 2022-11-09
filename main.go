package main

import (
	"bufio"
	"fmt"
	"github.com/darkturo/slocc/internal/pkg/filetype"
	"github.com/darkturo/slocc/internal/pkg/slocc"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	files := make([]string, 0, len(os.Args))
	excludedPaths := loadExcludedPaths()
	results := make(map[filetype.FileType]uint)
	for _, f := range os.Args[1:] {
		err := filepath.Walk(f, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() && !isExcluded(excludedPaths, path) {
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
		fileType := filetype.Guess(path)
		if fileType == filetype.Binary {
			continue
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("1 invalid file %s: %v\n", path, err)
			file.Close()
			continue
		}

		lines, err := slocc.CountLinesOfCode(fileType, bufio.NewReader(file))
		if err != nil {
			fmt.Printf("* invalid file %s: %v\n", path, err)
			file.Close()
			continue
		}

		results[fileType] += lines
		totalSLOC += lines
		file.Close()
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

func loadExcludedPaths() []string {
	excludedPaths := []string{
		".git/",
		".idea/",
		"vendor/",
		"node_modules/",
		"venv/",
		"__pycache__/",
		"*.egg-info/",
		"*.egg/",
		"*.pyc",
		"*.min.js",
		"*.min.css",
		"*.min.map",
		"*.map",
		"*.gz",
		"*.zip",
		"*.tar",
		"*.tar.gz",
	}

	// read .gitignore and add to excludedPaths
	file, err := os.Open(".gitignore")
	if err != nil {
		return excludedPaths
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		excludedPaths = append(excludedPaths, scanner.Text())
	}
	return excludedPaths
}

func isExcluded(excluded []string, path string) bool {
	for _, pattern := range excluded {
		match, err := filepath.Match(pattern, path)
		if err != nil {
			return false
		}
		if match {
			return true
		}
	}

	return false
}
