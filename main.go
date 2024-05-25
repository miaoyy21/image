package main

import (
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	root := "assets"

	if err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.EqualFold(filepath.Ext(path), ".mp4") {
			return nil
		}

		// 1280*720 = 16:9 = 320*180
		newName := strings.ReplaceAll(path, ".mp4", ".gif")
		if _, err := os.ReadFile(newName); err != nil && os.IsNotExist(err) {
			return ffmpeg.Input(path).
				Output(newName, ffmpeg.KwArgs{"s": "1280x720", "r": "15"}).
				OverWriteOutput().ErrorToStdOut().Run()
		}

		return nil
	}); err != nil {
		log.Printf("filepath.Walk Failure :: %s", err.Error())
	}
}
