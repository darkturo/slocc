package config

type IgnoredPaths []string

var defaultIgnoredPaths = IgnoredPaths{
	".git/",
	".idea/",
	"vendor/",
	"node_modules/",
	"venv/",
	"__pycache__/",
	"*.egg-info/",
	"*.egg/",
	"*.pyc",
	"*.gz",
	"*.zip",
	"*.tar",
	"*.tar.gz",
}
