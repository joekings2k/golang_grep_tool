package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
func walkDir(dir string, fun func(path string, info os.FileInfo, err error) error) error {
	return filepath.Walk(dir, fun)
}
func collectFiles(path []string) []string {
	var files []string
	for _, p := range path {
		fileinfo, err := os.Stat(p)
		if err != nil {
			log.Printf("Error accessing path %s: %v\n", p, err)
			continue
		}
		if fileinfo.IsDir() {

			//handle directory case
			err := walkDir(p, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				log.Printf("Error walking directory %s: %v\n", p, err)
				continue
			}
		} else {
			files = append(files, p)
		}

	}
	return files
}
func searchFile(fileName, searchTerm string) (bool, error) {
	match := false
	lineNumber := 0
	file, err := os.Open(fileName)
	if err != nil {
		return match, err
	}

	defer file.Close()
	searchTerm = strings.ToLower(searchTerm)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		originalLline := line
		lineNumber++
		line = strings.ToLower(line)

		if contains(line, searchTerm) {
			fmt.Printf("%s:%d: %s\n", fileName, lineNumber, originalLline)
			match = true
		}
	}
	if !match {
		fmt.Printf("%s:No matches found.", fileName)
	}
	if err := scanner.Err(); err != nil {
		return match, err
	}
	return match, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <filename> <search term>")
		return
	}
	searchTerm := os.Args[1]
	path := os.Args[2:]

	files := collectFiles(path)

	for _, file := range files {

		_, err := searchFile(file, searchTerm)
		if err != nil {
			log.Printf("Error searching in file %s: %v\n", file, err)
			continue
		}
	}

}
