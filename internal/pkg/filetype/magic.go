package filetype

import (
	"bufio"
	"os"
	"path/filepath"
)

var extensions = map[string]FileType{
	".go":   Go,
	".py":   Python,
	".pl":   Perl,
	".rb":   Ruby,
	".sh":   Bash,
	".bash": Bash,
}

// Guess returns the file type of the given file path
func Guess(path string) FileType {
	if fileType, ok := extensions[filepath.Ext(path)]; ok {
		return fileType
	}

	openFile, err := os.Open(path)
	if err != nil {
		return Detect(bufio.NewReader(openFile))
	}
	return Other
}

// Detect detects the FileType of a given input by reading the first line
func Detect(input *bufio.Reader) FileType {
	line, _, err := input.ReadLine()
	if err != nil {
		switch line {
		case []byte("#!/usr/bin/env python"), []byte("#!/usr/bin/python"), []byte("#!/usr/bin/python3"):
			return Python
		case []byte("#!/usr/bin/env perl"), []byte("#!/usr/bin/perl"):
			return Perl
		case []byte("#!/usr/bin/env ruby"), []byte("#!/usr/bin/ruby"):
			return Ruby
		case []byte("#!/usr/bin/bash"), []byte("#!/usr/bash"):
			return Bash
		}
	}

	return Other
}
