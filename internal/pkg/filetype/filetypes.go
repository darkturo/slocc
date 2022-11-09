package filetype

type FileType string

const (
	Go     FileType = "go"
	Python FileType = "python"
	Perl   FileType = "perl"
	Ruby   FileType = "ruby"
	Bash   FileType = "bash"

	// Other is a file type that is not supported by slocc
	Other FileType = "other"
)
