package slocc

import (
	"github.com/darkturo/slocc/internal/pkg/filetype"
)

type languageConfig struct {
	SingleLineCommentMarker   []string
	MultiLineBeginCommentMark string
	MultiLineEndCommentMark   string
}

var languageConfigurations = map[filetype.FileType]languageConfig{
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
}
