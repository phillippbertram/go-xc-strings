package internal

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// FindDuplicates searches for duplicate keys in a given .strings file or directory.
// FindDuplicates searches for duplicate keys and their values in a given .strings file or directory.
func FindDuplicates(path string) (DuplicatesMap, error) {
	result := make(DuplicatesMap)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	processFunc := func(path string) error {
		duplicates, err := processFileForDuplicates(path)
		if err != nil {
			return err
		}
		if len(duplicates) > 0 {
			result[path] = duplicates
		}
		return nil
	}

	if fileInfo.IsDir() {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".strings") {
				return processFunc(path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		err := processFunc(path)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func processFileForDuplicates(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	keys := make(map[string]string)
	duplicates := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		key, value := extractKeyValue(line)
		if key != "" {
			if _, found := keys[key]; found {
				duplicates[key] = value // Store duplicate with its value
			} else {
				keys[key] = value
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return duplicates, nil
}

// extractKeyValue extracts the key and value from a single line of a .strings file.
func extractKeyValue(line string) (string, string) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) < 2 {
		return "", ""
	}
	key := strings.Trim(parts[0], " \"")
	value := strings.Trim(parts[1], " \";")
	return key, value
}

// RemoveDuplicates removes all but the first occurrence of each duplicate key.
func RemoveDuplicates(basePath string, duplicates DuplicatesMap) error {
	for filePath, dupKeys := range duplicates {
		if err := removeExtraOccurrences(filePath, dupKeys); err != nil {
			return err
		}
	}
	return nil
}

func removeExtraOccurrences(filePath string, dupKeys map[string]string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	keysEncountered := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		key, _ := extractKeyValue(line)
		if _, found := dupKeys[key]; found && keysEncountered[key] {
			continue // skip this line because it's a duplicate
		}
		lines = append(lines, line)
		keysEncountered[key] = true
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
}
