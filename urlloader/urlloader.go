package urlloader

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	domainsURL = "https://tranco-list.eu/download/Z249G/full"
	tmpFile    = "domain_list.txt"
)

type URLLoader struct {
	file    *os.File
	scanner *bufio.Scanner
}

func New() (*URLLoader, error) {
	tmpFilePath := filepath.Join(os.TempDir(), tmpFile)
	_, err := os.Stat(tmpFilePath)
	if os.IsNotExist(err) {
		err := downloadFile(domainsURL, tmpFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to download file: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	file, err := os.Open(tmpFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	scanner := bufio.NewScanner(file)
	return &URLLoader{file: file, scanner: scanner}, nil
}

func (dl *URLLoader) NextURL() (string, error) {
	if dl.scanner.Scan() {
		splitted := strings.Split(dl.scanner.Text(), ",")
		if len(splitted) < 2 {
			return "", fmt.Errorf("domain line is invalid: %v", dl.scanner.Text())
		}
		return "https://" + splitted[1], nil
	}
	if err := dl.scanner.Err(); err != nil {
		return "", fmt.Errorf("scanner error: %w", err)
	}
	return "", io.EOF
}

func (dl *URLLoader) Close() error {
	err := dl.file.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	return nil
}

func downloadFile(url string, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get file from url: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	return nil
}