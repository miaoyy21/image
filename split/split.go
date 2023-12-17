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

	lines, title, parse := make([]string, 0), "", strings.Repeat("=", 80)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		line = strings.ReplaceAll(line, "《", "")
		line = strings.ReplaceAll(line, "》", "")
		line = strings.ReplaceAll(line, "女性", "女人")
		line = strings.ReplaceAll(line, "男性", "男人")

		if len(line) < 1 {
			continue
		}

		if strings.EqualFold(line, parse) {
			break
		}

		if len(title) < 1 {
			title = line
			continue
		}

		if strings.EqualFold(line, strings.Repeat(">", len(parse))) {
			continue
		}

		s1 := strings.Split(line, "，")
		for _, s := range s1 {
			s2 := strings.Split(s, "。")

			for _, s0 := range s2 {
				if len(s0) < 1 {
					continue
				}

				lines = append(lines, s0)
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

	buf := new(bytes.Buffer)

	// 写入标题
	if _, err := fmt.Fprintf(buf, "%s \n", title); err != nil {
		return err
	}

	// 写入分隔符
	if _, err := fmt.Fprintf(buf, "%s \n", strings.Repeat(">", len(parse))); err != nil {
		return err
	}

	// 写入每行
	for _, line := range lines {
		if _, err := fmt.Fprintf(buf, "%s， \n", line); err != nil {
			return err
		}
	}

	// 分隔符
	if _, err := fmt.Fprintf(buf, "%s \n", parse); err != nil {
		return err
	}

	// 全部行
	if _, err := fmt.Fprintf(buf, "%s \n", strings.Join(lines, "，")); err != nil {
		return err
	}

	if _, err := io.Copy(newFile, buf); err != nil {
		return err
	}

	log.Printf("Processing %q Finished ... \n", path)

	return nil
}
