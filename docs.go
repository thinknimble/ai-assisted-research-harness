package main

import (
	"os"
	"path/filepath"
	"strings"
)

var projectRoot string

func readFileContent(path string) (string, error) {
	resolved := path
	if !filepath.IsAbs(path) {
		resolved = filepath.Join(projectRoot, path)
	}
	cleaned := filepath.Clean(resolved)
	if !strings.HasPrefix(cleaned, projectRoot) {
		return "", os.ErrPermission
	}
	data, err := os.ReadFile(cleaned)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func listDir(dir string) ([]string, error) {
	root := filepath.Join(projectRoot, dir)
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			rel, _ := filepath.Rel(projectRoot, path)
			files = append(files, rel)
		}
		return nil
	})
	return files, err
}

func writeFile(dir, filename, content string) error {
	dest := filepath.Join(projectRoot, dir, filepath.Base(filename))
	cleaned := filepath.Clean(dest)
	if !strings.HasPrefix(cleaned, filepath.Join(projectRoot, dir)) {
		return os.ErrPermission
	}
	return os.WriteFile(cleaned, []byte(content), 0644)
}
