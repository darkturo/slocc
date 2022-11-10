package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/darkturo/slocc/internal/pkg/filetype"
	"os"
	"path/filepath"
)

const SettingsFile = "~/.slocc.json"

type Settings struct {
	// Ignore is a list of file patterns to ignore
	Ignore IgnoredPaths `json:"ignore"`

	// Languages contains the configuration for each language regarding what are comments within that lenguage
	Languages map[filetype.FileType]Lang `json:"languages"`
}

func (s Settings) IsIgnored(path string) bool {
	for _, pattern := range s.Ignore {
		match, err := filepath.Match(pattern, path)
		if err != nil {
			fmt.Printf("Warning: pattern '%s'", pattern)
			continue
		}
		if match {
			return true
		}
	}
	return false
}

func (s Settings) GetLang(fileType filetype.FileType) Lang {
	return s.Languages[fileType]
}

// LoadSettings loads slocc settings from the SettingsFile
// If there is no settings available it uses slocc's default settings to create one.
// If the file exists but is unreadable, it warns and returns the default settings
func LoadSettings() Settings {
	settings := getDefaultSettings()

	// Check if SETTINGS_FILE exists
	_, err := os.Stat(SettingsFile)
	if os.IsNotExist(err) {
		// create settings file with the defaults
		configFile, err := os.Create(SettingsFile)
		if err != nil {
			fmt.Printf("Warning: Error creating '%s': %v", SettingsFile, err)
			return settings
		}
		defer configFile.Close()

		b, err := json.Marshal(getDefaultSettings())
		if err != nil {
			fmt.Printf("Warning: Error marshalling default settings: %v", err)
			return settings
		}

		_, err = configFile.Write(b)
		if err != nil {
			fmt.Printf("Warning: Error writing to '%s': %v", SettingsFile, err)
			return settings
		}
	} else {
		configFile, err := os.Open(SettingsFile)
		if err != nil {
			fmt.Printf("Warning: Error opening '%s': %v ... ignoring local settings", SettingsFile, err)
			return settings
		}
		defer configFile.Close()

		err = json.NewDecoder(configFile).Decode(&settings)
		if err != nil {
			fmt.Printf("Warning: Error decoding '%s': %v ... ignoring local settings", SettingsFile, err)
			return getDefaultSettings()
		}
	}

	// Add .gitignore patterns to the ignore list
	settings.Ignore = append(settings.Ignore, loadGitIgnore()...)

	return settings
}

// getDefaultSettings loads the default settings hardcoded in the application
func getDefaultSettings() Settings {
	return Settings{
		Ignore:    defaultIgnoredPaths,
		Languages: defaultLanguages,
	}
}

// loadIgnoredPaths reads the .gitignore file and returns a list of patterns to ignore
func loadGitIgnore() IgnoredPaths {
	var gitIgnoreList IgnoredPaths

	file, err := os.Open(".gitignore")
	if err != nil {
		return gitIgnoreList
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		gitIgnoreList = append(gitIgnoreList, scanner.Text())
	}
	return gitIgnoreList
}
