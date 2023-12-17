package split

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
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

	buf := new(bytes.Buffer)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.ReplaceAll(line, "《", "")
		line = strings.ReplaceAll(line, "》", "")
		if len(line) < 1 {
			continue
		}

		if buf.Len() == 0 {
			if _, err := fmt.Fprintf(buf, "%s\n\n", line); err != nil {
				return err
			}

			continue
		}

		s1 := strings.Split(line, "，")
		for _, s := range s1 {
			s2 := strings.Split(s, "。")

			for _, s0 := range s2 {
				if len(s0) < 1 {
					continue
				}

				if _, err := fmt.Fprintf(buf, "%s\n", s0); err != nil {
					return err
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	newFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer newFile.Close()

	if _, err := io.Copy(newFile, buf); err != nil {
		return err
	}

	log.Printf("Processing %q Finished ... \n", path)

	return nil
}
