package config

import (
	"bufio"
	"os"
	"path/filepath"
)

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

func LoadIgnoredPaths() IgnoredPaths {

	// read .gitignore and add to excludedPaths
	file, err := os.Open(".gitignore")
	if err != nil {
		return defaultIgnoredPaths
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		defaultIgnoredPaths = append(defaultIgnoredPaths, scanner.Text())
	}
	return defaultIgnoredPaths
}

func IsExcluded(excluded IgnoredPaths, path string) bool {
	for _, pattern := range excluded {
		match, err := filepath.Match(pattern, path)
		if err != nil {
			return false
		}
		if match {
			return true
		}
	}

	return false
}
