package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/darkturo/slocc/internal/pkg/config"
	"github.com/darkturo/slocc/internal/pkg/filetype"
	"github.com/darkturo/slocc/internal/pkg/slocc"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

func init() {
	flag.Usage = func() {
		program := os.Args[0]
		if program[0:2] == "./" {
			program = program[2:]
		}
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <SOURCE_CODE_DIRs>\n", program)
		flag.PrintDefaults()
	}
}

type options struct {
	showHelp    bool
	directories []string
}

func parseArgs() options {
	opts := options{}

	flag.BoolVar(&opts.showHelp, "h", false, "Show usage information")
	flag.Parse()

	if opts.showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no directories specified\n")
		os.Exit(0)
	}

	opts.directories = flag.Args()

	return opts
}

func main() {
	opts := parseArgs()
	if opts.showHelp {
		flag.Usage()
		os.Exit(0)
	}

	files := make([]string, 0, len(os.Args))

	settings := config.LoadSettings()
	counter := slocc.New(settings)

	results := make(map[filetype.FileType]uint)
	for _, f := range opts.directories {
		err := filepath.Walk(f, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() && !settings.IsIgnored(path) {
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

		lines, err := counter.Count(fileType, bufio.NewReader(file))
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
