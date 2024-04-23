package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func CleanAndSortStringsFiles(stringsPath, stringsReferencePath, swiftDirectory string, ignorePatterns []string, sortFiles bool) error {
	// Detect unused keys
	unusedKeys, err := FindUnusedKeys(stringsReferencePath, swiftDirectory, ignorePatterns)
	if err != nil {
		return err
	}
	fmt.Printf("Found %d unused keys\n", len(unusedKeys))

	// using a map here as a more efficient data structure
	// to have a O(1) average-time complexity check instead of O(n)
	unusedKeysMap := SliceToMap(unusedKeys)

	if len(unusedKeysMap) == 0 {
		fmt.Print("No unused keys found, skipping deletion\n")
		return nil
	}

	if len(unusedKeysMap) == 0 && !sortFiles {
		fmt.Print("No unused keys found, skipping deletion and sorting\n")
		return nil
	}

	// CLEAN AND SORT
	err = filepath.Walk(stringsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".strings") {
			// Remove unused keys from the file
			if len(unusedKeysMap) > 0 {
				fmt.Printf("Cleaning %s\n", path)
				if err := cleanFile(path, unusedKeysMap); err != nil {
					return err
				}
			}

			// sort the file if requested
			if sortFiles {
				fmt.Printf("Sorting %s\n", path)
				if err := SortStringsFile(path); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

func cleanFile(path string, unusedKeys map[string]struct{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Implement logic to remove lines containing unused keys
	modifiedData, err := removeUnusedKeys(data, unusedKeys)
	if err != nil {
		return err
	}

	return os.WriteFile(path, modifiedData, 0644)
}

// Utility function to remove lines with unused keys
func removeUnusedKeys(data []byte, unusedKeys map[string]struct{}) ([]byte, error) {
	var modifiedData []byte
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		key := extractKeyFromLine(line)

		if _, exists := unusedKeys[key]; !exists {
			modifiedData = append(modifiedData, line...)
			modifiedData = append(modifiedData, '\n')
		}
	}
	return modifiedData, nil
}

// Utility function to extract key from a .strings file line
func extractKeyFromLine(line string) string {
	regex := regexp.MustCompile(`^"([^"]+)"`)
	matches := regex.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1] // Return the first group which is the key
	}
	return "" // Return an empty string if no match is found
}
