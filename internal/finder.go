package internal

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// locates keys in the .strings file that are not referenced in Swift files.
func FindUnusedKeys(referenceStringsPath, swiftDirectory string, ignorePatterns []string) ([]string, error) {
	entries, err := parseStringsFile(referenceStringsPath)
	if err != nil {
		return nil, err
	}

	keys := make(map[string]struct{})
	for _, entry := range entries {
		keys[entry.key] = struct{}{}
	}

	usedKeys := searchKeysInSwiftFiles(swiftDirectory, keys, ignorePatterns)
	unusedKeys := []string{}
	for key := range keys {
		if _, found := usedKeys[key]; !found {
			unusedKeys = append(unusedKeys, key)
		}
	}
	sort.Strings(unusedKeys)
	return unusedKeys, nil
}

func searchKeysInSwiftFiles(directory string, keys map[string]struct{}, ignorePatterns []string) map[string]struct{} {
	usedKeys := make(map[string]struct{})
	filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and files that match the ignore patterns
		for _, pattern := range ignorePatterns {
			if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Only process .swift files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".swift") {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Check if any of the keys are used in the file
			content := string(fileContent)
			for key := range keys {
				if strings.Contains(content, key) {
					usedKeys[key] = struct{}{}
				}
			}
		}
		return nil
	})
	return usedKeys
}
