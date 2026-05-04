package rules

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Result struct {
	Title   string `json:"title"`
	Path    string `json:"path"`
	Excerpt string `json:"excerpt"`
}

func Search(dirs []string, query string) ([]Result, error) {
	needle := strings.ToLower(strings.TrimSpace(query))
	if needle == "" {
		return nil, nil
	}
	var results []Result
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil || entry.IsDir() {
				return walkErr
			}
			if !isTextFile(path) {
				return nil
			}
			result, ok, err := searchFile(path, needle)
			if err != nil || !ok {
				return err
			}
			results = append(results, result)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func searchFile(path string, needle string) (Result, bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return Result{}, false, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), needle) {
			return Result{Title: filepath.Base(path), Path: path, Excerpt: strings.TrimSpace(line)}, true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return Result{}, false, err
	}
	return Result{}, false, nil
}

func isTextFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".txt", ".json":
		return true
	default:
		return false
	}
}
