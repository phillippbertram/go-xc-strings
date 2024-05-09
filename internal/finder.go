package internal

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func FindUnusedKeysInSwiftFiles(directory string, keys []string, ignorePatterns []string) []string {
	// fmt.Println("Searching for keys in Swift files...")
	// fmt.Println("Directory:", directory)
	// fmt.Println("Keys:", len(keys))
	// fmt.Println("Ignore patterns: ", ignorePatterns)

	keysMap := SliceToMap(keys) // more performant
	usedKeys := findKeysInSwiftFiles(directory, keys, ignorePatterns)

	// get unused keys
	unusedKeys := make(map[string]struct{})
	for key := range keysMap {
		if _, ok := usedKeys[key]; !ok {
			unusedKeys[key] = struct{}{}
		}
	}

	// Map to slice
	unusedKeysSlice := MapToSlice(unusedKeys)

	// sort the slice
	sort.Strings(unusedKeysSlice)

	return unusedKeysSlice
}

func findKeysInSwiftFiles(directory string, keys []string, ignorePatterns []string) map[string]struct{} {
	keysMap := SliceToMap(keys) // more performant
	usedKeys := make(map[string]struct{})

	_ = filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}
		// fmt.Printf("Processing %s\n", path)

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
			for key := range keysMap {
				if strings.Contains(content, key) {
					usedKeys[key] = struct{}{}
				}
			}
		}
		return nil
	})

	return usedKeys
}
