package config

import (
	"github.com/darkturo/slocc/internal/pkg/filetype"
)

type Lang struct {
	SingleLineCommentMarker   []string `json:"singleLineCommentMarker"`
	MultiLineBeginCommentMark string   `json:"multiLineBeginCommentMark"`
	MultiLineEndCommentMark   string   `json:"multiLineEndCommentMark"`
}

var defaultLanguages = map[filetype.FileType]Lang{
	filetype.Go: {
		SingleLineCommentMarker:   []string{"//"},
		MultiLineBeginCommentMark: "/*",
		MultiLineEndCommentMark:   "*/",
	},
	filetype.Python: {
		SingleLineCommentMarker:   []string{"#"},
		MultiLineBeginCommentMark: "\"\"\"",
		MultiLineEndCommentMark:   "\"\"\"",
	},
	filetype.Perl: {
		SingleLineCommentMarker:   []string{"#"},
		MultiLineBeginCommentMark: "=pod",
		MultiLineEndCommentMark:   "=cut",
	},
	filetype.Ruby: {
		SingleLineCommentMarker:   []string{"#"},
		MultiLineBeginCommentMark: "=begin",
		MultiLineEndCommentMark:   "=end",
	},
	filetype.Bash: {
		SingleLineCommentMarker:   []string{"#"},
		MultiLineBeginCommentMark: "",
		MultiLineEndCommentMark:   "",
	},
	filetype.Cpp: {
		SingleLineCommentMarker:   []string{"//"},
		MultiLineBeginCommentMark: "/*",
		MultiLineEndCommentMark:   "*/",
	},
	filetype.Cuda: {
		SingleLineCommentMarker:   []string{"//"},
		MultiLineBeginCommentMark: "/*",
		MultiLineEndCommentMark:   "*/",
	},
}
