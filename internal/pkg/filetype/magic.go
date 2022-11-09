package filetype

import "path/filepath"

var extensions = map[string]FileType{
	".go": Go,
}

// Guess returns the file type of the given file path
func Guess(path string) FileType {
	if fileType, ok := extensions[filepath.Ext(path)]; ok {
		return fileType
	}
	return Other
}
