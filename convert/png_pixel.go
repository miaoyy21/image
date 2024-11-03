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

func PNGToPixel() error {
	return filepath.Walk("images", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dir, name := filepath.Split(path)

		if !strings.EqualFold(strings.ToLower(filepath.Ext(path)), ".png") {
			return nil
		}

		if strings.HasPrefix(strings.ToLower(name), "src_") || strings.HasPrefix(strings.ToLower(name), "dst_") {
			return nil
		}

		// 拷贝1份原始文件
		srcPath := fmt.Sprintf("%ssrc_%s", dir, name)
		if err := copyFile(path, srcPath); err != nil {
			return err
		}

		srcFile, err := os.Open(srcPath)
		if err != nil {
			return err
		}

		// 原始图像
		srcImage, _, err := image.Decode(srcFile)
		if err != nil {
			return err
		}

		if err := srcFile.Close(); err != nil {
			return err
		}

		bounds := srcImage.Bounds()

		// 目标图像
		dstImage := image.NewRGBA(
			image.Rectangle{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{X: bounds.Dx(), Y: bounds.Dy()},
			},
		)

		for x := 0; x < bounds.Dx(); x = x + 1 {
			for y := 0; y < bounds.Dy(); y = y + 1 {
				r0, g0, b0, a0 := srcImage.(*image.RGBA).At(x, y).RGBA()

				r1, g1, b1, a1 := r0>>12, g0>>12, b0>>12, a0>>12

				newRGBA := color.RGBA{R: uint8(r1 << 4), G: uint8(g1 << 4), B: uint8(b1 << 4), A: uint8(a1 << 4)}
				dstImage.SetRGBA(x, y, newRGBA)
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
