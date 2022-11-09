package filetype

type FileType string

const (
	Go FileType = "go"

	// Other is a file type that is not supported by slocc
	Other FileType = "other"
)
