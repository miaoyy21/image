package split

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func Split(fileName string) error {
	err := filepath.Walk(fileName, split)
	if err != nil {
		return err
	}

	return nil
}

func split(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if f.IsDir() {
		return nil
	}

	if !strings.HasSuffix(path, "文案") {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		scanner.Text()

	}

	return nil
}
