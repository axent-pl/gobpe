package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var downloadPath string

	flag.StringVar(&downloadPath, "dst", "", "Download path")
	flag.Parse()

	if downloadPath == "" {
		usage()
		os.Exit(1)
	}

	links, err := readStdin()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		usage()
		os.Exit(1)
	}

	download(links, downloadPath)
}

func usage() {
	flag.Usage()
}

func readStdin() ([]string, error) {
	var lines []string

	fi, err := os.Stdin.Stat()
	if err != nil {
		return lines, err
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return lines, fmt.Errorf("no data in stdin")
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return lines, err
	}

	if len(lines) == 0 {
		return lines, fmt.Errorf("no lines")
	}
	return lines, nil
}

func download(links []string, path string) ([]string, error) {
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
