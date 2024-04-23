package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func SortStringsFiles(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".strings" {
				return SortStringsFile(path)
			}
			return nil
		})
	} else {
		return SortStringsFile(path)
	}
}

func SortStringsFile(filepath string) error {
	entries, err := parseStringsFile(filepath)
	if err != nil {
		return err
	}

	// Sort entries
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].key < entries[j].key
	})

	// Write back to the same file or a new file
	return writeSortedStringsFile(filepath, entries)
}

type Entry struct {
	key   string
	value string
}

func parseStringsFile(filepath string) ([]Entry, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []Entry
	scanner := bufio.NewScanner(file)
	regex := regexp.MustCompile(`"(.+?)"\s*=\s*"(.*?)";`)
	for scanner.Scan() {
		matches := regex.FindStringSubmatch(scanner.Text())
		if len(matches) > 1 {
			entries = append(entries, Entry{key: matches[1], value: matches[2]})
		}
	}
	return entries, scanner.Err()
}

func writeSortedStringsFile(filepath string, entries []Entry) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lastPrefix string
	for _, entry := range entries {
		currentPrefix := strings.Split(entry.key, "_")[0]
		if currentPrefix != lastPrefix && lastPrefix != "" {
			_, err := file.WriteString("\n") // Add a blank line between groups
			if err != nil {
				return err
			}
		}
		_, err := file.WriteString(fmt.Sprintf("\"%s\"=\"%s\";\n", entry.key, entry.value))
		if err != nil {
			return err
		}
		lastPrefix = currentPrefix
	}
	return nil
}
