package convert

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func RemoveBlack() error {
	return filepath.Walk("images", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dir, name := filepath.Split(path)

		if !strings.EqualFold(strings.ToLower(filepath.Ext(path)), ".png") {
			return nil
		}
		//
		//if !strings.Contains(strings.ToLower(name), "_1.png") {
		//	return nil
		//}

		orgFile, err := os.Open(path)
		if err != nil {
			return err
		}

		// 原始图像
		orgImage, _, err := image.Decode(orgFile)
		if err != nil {
			return err
		}

		if err := orgFile.Close(); err != nil {
			return err
		}

		bounds := orgImage.Bounds()

		// 目标图像
		dstImage := image.NewNRGBA(
			image.Rectangle{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{X: bounds.Dx(), Y: bounds.Dy()},
			},
		)

		for x := 0; x < bounds.Dx(); x++ {
			for y := 0; y < bounds.Dy(); y++ {
				r, g, b, a := orgImage.(*image.NRGBA).At(x, y).RGBA()
				if r>>12 <= 0 && g>>12 <= 0 && b>>12 <= 0 {
					dstImage.SetNRGBA(x, y, color.NRGBA{R: 0, G: 0, B: 0, A: 0})
				} else {
					newRGBA := color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
					dstImage.SetNRGBA(x, y, newRGBA)
				}
			}
		}

		dstPath := fmt.Sprintf("%sdst_%s", dir, name)
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}

		if err := png.Encode(dstFile, dstImage); err != nil {
			return err
		}

		if err := dstFile.Close(); err != nil {
			return err
		}

		log.Printf("%q Deal Finished ... \n", path)

		return nil
	})
}
