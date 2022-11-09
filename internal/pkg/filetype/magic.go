package filetype

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Guess returns the file type of the given file path
// If the given file is a binary file, it returns Binary.
// Otherwise, it will try to guess the file type by its extension or by its shebang (scripts).
func Guess(path string) FileType {
	openFile, err := os.Open(path)
	if err != nil {
		return Other
	}
	defer openFile.Close()

	if looksLikeBinary(bufio.NewReader(openFile)) {
		return Other
	}

	if fileType, ok := extensions[filepath.Ext(path)]; ok {
		return fileType
	}

	return shebang(bufio.NewReader(openFile))
}

// looksLikeBinary returns true if the given file has a 0x0 byte
func looksLikeBinary(input *bufio.Reader) bool {
	bytes, err := input.Peek(1024)
	if err != nil && err != io.EOF {
		return false
	}

	for _, b := range bytes {
		if b == 0 {
			return true
		}
	}
	return false
}

// shebang detects the FileType of a given input by reading the first line
func shebang(input *bufio.Reader) FileType {
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
