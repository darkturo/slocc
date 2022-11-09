package filetype

var extensions = map[string]FileType{
	".go":   Go,
	".py":   Python,
	".pl":   Perl,
	".rb":   Ruby,
	".sh":   Bash,
	".bash": Bash,
}
