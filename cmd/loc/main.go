package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func main() {
	// Map to store the line counts per extension
	lineCounts := make(map[string]int)

	// Walk through the directory
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if slices.Contains([]string{".go", ".tmpl", ".css", ".md", ".sql"}, ext) {
				lines, err := countLines(path)
				if err != nil {
					fmt.Printf("Error reading file %s: %v\n", path, err)
					return nil
				}
				lineCounts[ext] += lines
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path: %v\n", err)
		return
	}

	fmt.Println("Line counts by extension:")
	for ext, count := range lineCounts {
		fmt.Printf("%s: %d\n", ext, count)
	}
}

// countLines reads a file and returns the number of lines
func countLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var lines int
	buffer := make([]byte, 8192)
	for {
		n, err := file.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			return lines, err
		}
		lines += strings.Count(string(buffer[:n]), "\n")
		if n == 0 {
			break
		}
	}
	return lines, nil
}
