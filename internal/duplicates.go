package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindDuplicates(path string) (DuplicatesMap, error) {
	result := make(DuplicatesMap)
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".strings") {
			duplicates, err := processFileForDuplicates(filePath)
			if err != nil {
				return err
			}
			if len(duplicates) > 0 {
				result[filePath] = duplicates
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func processFileForDuplicates(filePath string) (map[string][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	keys := make(map[string][]string)
	for scanner.Scan() {
		line := scanner.Text()
		key, value := extractKeyValue(line)
		if key != "" {
			keys[key] = append(keys[key], value)
		}
	}

	duplicates := make(map[string][]string)
	for key, values := range keys {
		if len(values) > 1 {
			seen := make(map[string]bool)
			var uniqueValues []string
			for _, value := range values {
				if _, found := seen[value]; !found {
					seen[value] = true
					uniqueValues = append(uniqueValues, value)
				}
			}
			if len(uniqueValues) > 1 {
				duplicates[key] = uniqueValues
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return duplicates, nil
}

// RemoveDuplicates removes all but the first occurrence of each duplicate key.
func RemoveDuplicates(basePath string, duplicates DuplicatesMap) error {
	for filePath, keys := range duplicates {
		if err := removeExtraOccurrences(filePath, keys); err != nil {
			return fmt.Errorf("failed to remove duplicates from %s: %w", filePath, err)
		}
	}
	return nil
}

func removeExtraOccurrences(filePath string, dupKeys map[string][]string) error {
	input, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer input.Close()

	var lines []string
	encounteredKeys := make(map[string]bool)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		key, value := extractKeyValue(line)
		if dupValues, found := dupKeys[key]; found {
			// Check if this value is the first occurrence and not already added.
			if !encounteredKeys[key] && len(dupValues) > 0 && value == dupValues[0] {
				lines = append(lines, line)
				encounteredKeys[key] = true // Mark this key as added.
			}
		} else {
			// Key is not a duplicate; add it to the output.
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Rewrite the file with duplicates removed.
	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
}
