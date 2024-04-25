package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// RemoveKeyFromAllStringsFiles removes a specific key from all .strings files within a directory.
func RemoveKeyFromAllStringsFiles(key, directory string, excludeLanguages []string) ([]string, error) {
	var removedFromFiles []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".strings") {

			// Skip files that are in the excludeLanguages list
			// Check if parent directory is in the following format: "<language>.lproj"
			// if it is, ignore the file

			// Check if the parent directory is in the format "<language>.lproj"
			parentDir := filepath.Base(filepath.Dir(path))
			for _, lang := range excludeLanguages {
				if parentDir == fmt.Sprintf("%s.lproj", lang) {
					return nil
				}
			}

			wasRemoved, err := removeKeyFromFile(key, path)
			if err != nil {
				return err
			}
			if wasRemoved {
				removedFromFiles = append(removedFromFiles, path)
			}

		}
		return nil
	})

	return removedFromFiles, err
}

// removeKeyFromFile removes a specific key from a .strings file.
func removeKeyFromFile(key, filepath string) (bool, error) {
	// TODO: use parseStringsFile and writeStringsFile

	file, err := os.Open(filepath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Construct the regex to match the key with various whitespace patterns around `=`
	regexPattern := fmt.Sprintf(`^"\s*%s\s*"(\s*=\s*".*";)$`, regexp.QuoteMeta(key))
	keyRegex, err := regexp.Compile(regexPattern)
	if err != nil {
		return false, err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	wasRemoved := false

	// Read file and filter out the specified key
	for scanner.Scan() {
		line := scanner.Text()
		if !keyRegex.MatchString(line) {
			lines = append(lines, line)
		} else {
			wasRemoved = true
		}
	}
	if err := scanner.Err(); err != nil {
		return wasRemoved, err
	}

	// Write back if key was removed
	if wasRemoved {
		return wasRemoved, os.WriteFile(filepath, []byte(strings.Join(lines, "\n")), 0644)
	}

	return wasRemoved, nil
}
