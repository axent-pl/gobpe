package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

func mustReadFile(filepath string) []byte {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return bytes
}

func mustWriteToFile(filename string, data []byte) {
	err := os.WriteFile(filename, data, 0600)
	if err != nil {
		panic(err)
	}
}

func listFiles(path string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
