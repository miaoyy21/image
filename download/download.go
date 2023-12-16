package download

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func Download(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// 自动创建对应的文件目录
	dirName := fmt.Sprintf("%s_Download", fileName)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			return err
		}
	}

	dUrls := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "curl") {
			continue
		}

		if !strings.Contains(line, "https://images.pexels.com/photos/") {
			continue
		}

		dUrl := line[6 : len(line)-3]
		dUrls = append(dUrls, dUrl)
	}

	// 检查是否有错误发生
	if err := scanner.Err(); err != nil {
		return err
	}

	var wg sync.WaitGroup

	wg.Add(len(dUrls))
	for i, dUrl := range dUrls {
		out := filepath.Join(dirName, fmt.Sprintf("%03d", i))
		if len(dUrls) > 1000 {
			out = filepath.Join(dirName, fmt.Sprintf("%04d", i))
		}
		
		percent := fmt.Sprintf("% 4d (% 4d)", i, len(dUrls))
		go download(&wg, dUrl, out, percent)
	}

	wg.Wait()

	return nil
}

func download(wg *sync.WaitGroup, dUrl, out, percent string) {
	defer wg.Done()

	if strings.Contains(dUrl, ".png?") {
		out = fmt.Sprintf("%s.png", out)
	} else if strings.Contains(dUrl, ".jpeg?") {
		out = fmt.Sprintf("%s.jpeg", out)
	} else if strings.Contains(dUrl, ".jpg?") {
		out = fmt.Sprintf("%s.jpg", out)
	} else {
		log.Fatalf("unknown file type %s \n", dUrl)
	}

	resp, err := http.Get(dUrl)
	if err != nil {
		log.Fatalf("http.Get Failure :: %s \n", err.Error())
	}

	file, err := os.Create(out)
	if err != nil {
		log.Fatalf("os.Create Failure :: %s \n", err.Error())
	}

	if _, err := io.Copy(file, resp.Body); err != nil {
		log.Fatalf("io.Copy Failure :: %s \n", err.Error())
	}

	log.Printf("Downloading %s => %s ... \n", percent, dUrl)
}
