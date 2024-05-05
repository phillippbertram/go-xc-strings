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
	fmt.Printf("Processing %s\n", filePath)
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

	// find duplicates where the key has more than one value
	duplicates := make(map[string][]string)
	for key, values := range keys {
		if len(values) > 1 {
			duplicates[key] = values
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return duplicates, nil
}

// RemoveDuplicatesKeepLast removes all but the first occurrence of each duplicate key.
func RemoveDuplicatesKeepLast(basePath string, duplicates DuplicatesMap) error {
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

	lines := []string{}
	lastOccurrences := make(map[string]int) // Map to store the last occurrence of each key

	scanner := bufio.NewScanner(input)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		key, _ := extractKeyValue(line)
		if _, found := dupKeys[key]; found {
			lastOccurrences[key] = i // Store the index of the last occurrence of the key
		} else {
			lines = append(lines, line) // Append non-duplicate lines directly
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Reset and read the file again to collect the correct lines
	if _, err = input.Seek(0, 0); err != nil {
		return err
	}
	scanner = bufio.NewScanner(input)
	i = 0
	finalLines := make([]string, len(lines))
	copy(finalLines, lines) // Start with non-duplicate lines

	for scanner.Scan() {
		if index, found := lastOccurrences[extractKey(scanner.Text())]; found && index == i {
			finalLines = append(finalLines, scanner.Text()) // Append only the last occurrence
		}
		i++
	}

	return os.WriteFile(filePath, []byte(strings.Join(finalLines, "\n")), 0644)
}
