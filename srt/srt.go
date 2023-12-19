package srt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func Srt(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	buf, index := &bytes.Buffer{}, 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		index++

		if index%4 != 3 {
			continue
		}

		line := scanner.Text()
		buf.WriteString(line)
		buf.WriteByte('\n')
	}

	// 检查是否有错误发生
	if err := scanner.Err(); err != nil {
		return err
	}

	newFileName := fmt.Sprintf("%s.txt", fileName[:len(fileName)-4])
	newFile, err := os.Create(newFileName)
	if err != nil {
		return err
	}
	defer newFile.Close()

	if _, err := io.Copy(newFile, buf); err != nil {
		return err
	}

	log.Printf("Withdraw File %q Finished ... \n", newFileName)
	return nil
}
