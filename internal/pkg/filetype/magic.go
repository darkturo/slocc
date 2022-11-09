package filetype

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Guess returns the file type of the given file path
func Guess(path string) FileType {
	if fileType, ok := extensions[filepath.Ext(path)]; ok {
		return fileType
	}

	openFile, err := os.Open(path)
	if err != nil {
		return Shebang(bufio.NewReader(openFile))
	}
	return Other
}

// Shebang detects the FileType of a given input by reading the first line
func Shebang(input *bufio.Reader) FileType {
	line, _, err := input.ReadLine()
	if err == nil {
		if len(line) > 2 && line[0] == '#' && line[1] == '!' {
			switch {
			case strings.Contains(string(line), "python"):
				return Python
			case strings.Contains(string(line), "perl"):
				return Perl
			case strings.Contains(string(line), "ruby"):
				return Ruby
			case strings.Contains(string(line), "bash"):
				return Bash
			}
		}
	}

	return Other
}
