package convert

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func PNGToGIF() error {
	pngNames := make([]string, 0)

	if err := filepath.Walk("images", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		pngName := info.Name()

		if !strings.EqualFold(strings.ToLower(filepath.Ext(path)), ".png") {
			return nil
		}

		if strings.HasPrefix(strings.ToLower(pngName), "src_") || strings.HasPrefix(strings.ToLower(pngName), "dst_") {
			return nil
		}

		pngNames = append(pngNames, pngName)
		return nil
	}); err != nil {
		return err
	}

	log.Println(pngNames)

	return nil
}
