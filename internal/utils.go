package internal

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Contains checks if a string is contained in a slice of strings.
func Contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// Convert a slice to a map for quicker lookups
func SliceToMap(slice []string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, item := range slice {
		result[item] = struct{}{} // Use an empty struct to minimize memory usage
	}
	return result
}

func MapToSlice(m map[string]struct{}) []string {
	result := make([]string, 0, len(m))
	for key := range m {
		result = append(result, key)
	}
	return result

}

// IsDirectory determines if a file represented
// by `path` is a directory or not
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func FindDefaultLanguageForXcodeProject(projectPath string) (string, error) {
	pbxprojPath, err := findPBXProjPath(projectPath)
	if err != nil {
		return "", fmt.Errorf("error finding .pbxproj file: %w", err)
	}
	fmt.Printf("Found .pbxproj file: %s\n", pbxprojPath)

	language, err := findDevelopmentRegionInPbxProj(pbxprojPath)
	if err != nil {
		return "", fmt.Errorf("error finding development region: %w", err)
	}

	return language, nil
}

// locates the .pbxproj file starting from the given path.
// If the path directly ends with .pbxproj, it returns the path.
// Otherwise, it searches within the directory and its subdirectories.
func findPBXProjPath(basePath string) (string, error) {
	// Normalize the base path
	basePath = filepath.Clean(basePath)

	// Check if the given path is directly a .pbxproj file
	if strings.HasSuffix(basePath, ".pbxproj") {
		if _, err := os.Stat(basePath); err != nil {
			return "", err
		}
		return basePath, nil
	}

	// Define an error to return in case no .pbxproj file is found
	var errNotFound = errors.New(".pbxproj file not found in the directory or its subdirectories")

	// Search for .pbxproj files in the directory and subdirectories
	foundPath := ""
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".pbxproj") {
			foundPath = path
			return filepath.SkipDir // Stop walking the directory tree
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", errNotFound
	}

	return foundPath, nil
}

// searches for the developmentRegion in an Xcode project.pbxproj file.
func findDevelopmentRegionInPbxProj(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Look for the developmentRegion key in the file
		if strings.Contains(line, "developmentRegion") {
			// Assuming the line format is developmentRegion = en;
			parts := strings.Split(line, "=")
			if len(parts) > 1 {
				// Clean up the parsed region value
				region := strings.TrimSpace(parts[1])
				region = strings.Trim(region, ";")
				region = strings.Trim(region, "\"")
				return region, nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("developmentRegion not found")
}

// constructs the path to the Localizable.strings file based on the development region.
func GetLocalizableStringsPath(baseDir, devRegion string) (string, error) {
	lprojPath := filepath.Join(baseDir, fmt.Sprintf("%s.lproj", devRegion), "Localizable.strings")

	// Check if the Localizable.strings file exists at the constructed path
	if _, err := os.Stat(lprojPath); os.IsNotExist(err) {
		return "", fmt.Errorf("Localizable.strings not found for region '%s' at path: %s", devRegion, lprojPath)
	} else if err != nil {
		return "", err
	}

	return lprojPath, nil
}

// extracts the key and value from a single line of a .strings file.
func extractKeyValue(line string) (string, string) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) < 2 {
		return "", ""
	}
	key := strings.Trim(parts[0], " \"")
	value := strings.Trim(parts[1], " \";")
	return key, value
}

// Helper function to extract a key from a line.
func extractKey(line string) string {
	key, _ := extractKeyValue(line)
	return key
}
