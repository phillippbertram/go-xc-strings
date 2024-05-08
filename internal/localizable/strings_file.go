package localizable

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

type Line struct {
	Key        string // Empty if the line is not a key-value pair
	Value      string // Empty if the line is not a key-value pair
	Text       string // Raw text of the line
	LineNumber int    // Line number in the file
}

type StringsFile struct {
	Path  string
	Lines []Line
}

func NewStringsFile(path string) (*StringsFile, error) {
	str := &StringsFile{
		Path:  path,
		Lines: make([]Line, 0),
	}

	err := str.parse()
	return str, err
}

func (sf *StringsFile) parse() error {
	file, err := os.Open(sf.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	sf.Lines = nil // Reset lines
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		parsedLine := Line{
			Text:       line,
			LineNumber: lineNumber,
		}

		// Attempt to parse as key-value if possible
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key := strings.Trim(parts[0], " \"")
			value := strings.Trim(parts[1], " \";")
			parsedLine.Key = key
			parsedLine.Value = value
		}

		sf.Lines = append(sf.Lines, parsedLine)
		lineNumber++
	}

	return scanner.Err()
}

func (sf *StringsFile) GetLinesForKey(key string) []Line {
	var linesForKey []Line

	for _, line := range sf.Lines {
		if line.Key == key {
			linesForKey = append(linesForKey, line)

		}
	}

	return linesForKey
}

func (sf *StringsFile) FindDuplicateKeys() map[string][]Line {
	keyLines := make(map[string][]Line)
	duplicates := make(map[string][]Line)

	for _, line := range sf.Lines {
		if line.Key != "" {
			keyLines[line.Key] = append(keyLines[line.Key], line)
			if len(keyLines[line.Key]) > 1 { // This line is checked on every addition
				duplicates[line.Key] = keyLines[line.Key]
			}
		}
	}

	// Ensure only keys with duplicates are in the final map
	for key, lines := range duplicates {
		if len(lines) < 2 {
			delete(duplicates, key)
		}
	}

	return duplicates
}

func (sf *StringsFile) RemoveKey(key string) []Line {
	var removedLines []Line
	var newLines []Line
	for _, line := range sf.Lines {
		if line.Key != key {
			newLines = append(newLines, line)
		} else {
			removedLines = append(removedLines, line)
		}
	}
	sf.Lines = newLines
	return removedLines
}

func (sf *StringsFile) Sort() {
	// Filter out empty lines and non-key-value lines
	var keyLines []Line
	for _, line := range sf.Lines {
		if line.Key != "" && strings.TrimSpace(line.Text) != "" {
			keyLines = append(keyLines, line)
		}
	}

	// Sort lines by key
	sort.Slice(keyLines, func(i, j int) bool {
		return keyLines[i].Key < keyLines[j].Key
	})

	// Group by prefix and insert empty lines
	var sortedLines []Line
	currentPrefix := ""

	for i, line := range keyLines {
		if i == 0 || !strings.HasPrefix(line.Key, currentPrefix) {
			if i != 0 {
				// Add an empty line to separate groups
				sortedLines = append(sortedLines, Line{Text: ""})
			}
			currentPrefix = string(line.Key[0]) // Using the first character as prefix
		}
		sortedLines = append(sortedLines, line)
	}

	sf.Lines = sortedLines
}

func (sf *StringsFile) Save() error {
	file, err := os.Create(sf.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range sf.Lines {
		if _, err := writer.WriteString(line.Text + "\n"); err != nil {
			return err
		}
	}

	return writer.Flush()
}
