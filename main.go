package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type result struct {
	fileName   string
	lineNumber int
	line       string
}

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
func searchFile(fileName, searchTerm string, results chan result) (bool, error) {
	match := false
	lineNumber := 0
	file, err := os.Open(fileName)
	if err != nil {
		return match, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		originalLine := line
		lineNumber++
		line = strings.ToLower(line)

		if contains(line, searchTerm) {
			results <- result{fileName: fileName, lineNumber: lineNumber, line: originalLine}
			match = true
		}
	}
	if !match {

	}
	if err := scanner.Err(); err != nil {
		return match, err
	}
	return match, nil
}
func searchWorker(jobs <-chan string, searchTerm string, results chan result) {
	for file := range jobs {
		_, err := searchFile(file, searchTerm, results)
		if err != nil {
			log.Printf("Error searching in file %s: %v\n", file, err)
			continue
		}
	}
}

func main() {
	var wg sync.WaitGroup
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <filename> <search term>")
		return
	}
	searchTerm := strings.ToLower(os.Args[1])
	path := os.Args[2:]

	files := collectFiles(path)
	jobs := make(chan string)
	results := make(chan result)

	numOfWorker := 4
	wg.Add(numOfWorker)

	for i := 0; i < numOfWorker; i++ {
		go func() {

			defer wg.Done()
			searchWorker(jobs, searchTerm, results)
		}()
	}

	for _, file := range files {
		jobs <- file
	}
	close(jobs)
	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Printf("%s:%d: %s\n", r.fileName, r.lineNumber, r.line)
	}

}
