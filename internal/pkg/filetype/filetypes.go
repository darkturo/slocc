package filetype

type FileType string

const (
	Go     FileType = "go"
	Python FileType = "python"
	Perl   FileType = "perl"
	Ruby   FileType = "ruby"
	Bash   FileType = "bash"
	Cpp    FileType = "cpp"
	Cuda   FileType = "cuda"

	Binary FileType = "binary"
	Text   FileType = "text"
	// Other is a file type that is not supported by slocc
	Other FileType = "other"
)
