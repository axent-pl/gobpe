package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func MustWriteToFile(filename string, data []byte) {
	err := os.WriteFile(filename, data, 0600)
	if err != nil {
		panic(err)
	}
}

func MustReadFile(filepath string) []byte {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return bytes
}

func LoadLinksFromFile(filepath string) ([]string, error) {
	var links []string = make([]string, 0)
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		links = append(links, scanner.Text())
	}
	return links, nil
}

func Download(links []string, path string) ([]string, error) {
	var paths []string = make([]string, 0)
	for _, link := range links {
		log.Printf("Downloading link %s", link)
		resp, err := http.Get(link)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("http status: %s", resp.Status)
		}

		segments := strings.Split(link, "/")
		filename := segments[len(segments)-1]
		filePath := filepath.Join(path, filename)

		if _, err := os.Stat(filePath); err == nil {
			log.Printf("Downloading link %s, already exists", link)
			paths = append(paths, filePath)
			continue
		}

		file, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return nil, err
		}
		paths = append(paths, filePath)
	}
	return paths, nil
}
