package download

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func DownloadText(fileName string) error {
	srcFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 自动创建对应的文件目录
	dirName := fmt.Sprintf("%s_Download", fileName)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			return err
		}
	}

	// 打开文件以追加模式写入
	target := filepath.Join(dirName, "TODO.txt")
	dstFile, err := os.OpenFile(target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("os.OpenFile Failure :: %s", err)
	}
	defer dstFile.Close()

	fmt.Println("文件追加写入成功")

	index := 0
	scanner := bufio.NewScanner(srcFile)
	writer := bufio.NewWriter(dstFile)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "curl") {
			continue
		}

		if !strings.Contains(line, "https://images.pexels.com/photos/") {
			continue
		}

		dUrl := line[6 : len(line)-3]
		ext, err := parseExt(dUrl)
		if err != nil {
			return err
		}

		index++
		text := fmt.Sprintf("curl %s --output %04d.%s", dUrl, index, ext)
		if _, err := fmt.Fprintln(writer, text); err != nil {
			return fmt.Errorf("fmt.Fprintln Failure :: %s", err)
		}
	}

	// 检查是否有错误发生
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner.Err Failure :: %s", err)
	}

	if err = writer.Flush(); err != nil {
		return fmt.Errorf("writer.Flush Failure :: %s", err)
	}

	log.Printf("Setting Download File Finished , Totals %d ... \n", index)

	return nil
}
