package localizable

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type DuplicateKeys struct {
	Duplicates map[string][]Line
	FilePath   string
}

type StringsFileManager struct {
	Paths []string // This can include file paths or glob patterns
	Files []*StringsFile
}

func NewStringsFileManager(paths []string) (*StringsFileManager, error) {
	man := &StringsFileManager{
		Paths: paths,
		Files: make([]*StringsFile, 0),
	}

	err := man.parseFiles()
	if err != nil {
		return nil, err
	}

	return man, nil
}

func (m *StringsFileManager) FindDuplicates() map[string]*DuplicateKeys {
	duplicatesPerFile := make(map[string]*DuplicateKeys)

	for _, file := range m.Files {
		fmt.Printf("Finding duplicates in file: %s\n", file.Path)
		duplicates := file.FindDuplicateKeys()
		if len(duplicates) > 0 {
			duplicatesPerFile[file.Path] = &DuplicateKeys{
				FilePath:   file.Path,
				Duplicates: duplicates,
			}
		}
	}

	return duplicatesPerFile
}

func (m *StringsFileManager) Sanitize() {
	for _, file := range m.Files {
		fmt.Printf("Sanitizing file: %s\n", file.Path)
		file.Sanitize()
	}
}

func (m *StringsFileManager) Sort() {
	for _, file := range m.Files {
		fmt.Printf("Sorting file: %s\n", file.Path)
		file.Sort()
	}
}

func (m *StringsFileManager) Save() {
	for _, file := range m.Files {
		fmt.Printf("Saving file: %s\n", file.Path)
		file.Save()
	}
}

func (m *StringsFileManager) parseFiles() error {
	var err error
	for _, path := range m.Paths {
		fmt.Printf("Processing path: %s\n", path)

		// Check if the path is a directory
		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			fmt.Printf("Path is a directory: %s\n", path)

			// If it's a directory, walk the directory
			filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() && strings.HasSuffix(d.Name(), ".strings") {
					m.parseFile(p)
				}
				return nil
			})
		} else {
			// Handle it as a glob pattern
			matches, err := filepath.Glob(path)
			fmt.Printf("Glob Matches: %v\n", matches)
			if err != nil {
				log.Printf("Error interpreting glob or path '%s': %s", path, err)
				continue
			}
			for _, match := range matches {
				m.parseFile(match)
			}
		}
	}
	return err
}

func (manager *StringsFileManager) parseFile(path string) error {
	sf, err := NewStringsFile(path)
	if err != nil {
		return err
	}
	manager.Files = append(manager.Files, sf)
	return nil
}
