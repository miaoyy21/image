package main

import (
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// MP4 转为 GIF
	if err := MP4ToGIF(); err != nil {
		log.Printf("MP4ToGIF() Failure :: %s", err.Error())
		return
	}

	// 合并像素点
	if err := PNGToHD(); err != nil {
		log.Printf("PNGToHD() Failure :: %s", err.Error())
		return
	}
}

func MP4ToGIF() error {
	return filepath.Walk("videos", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.EqualFold(strings.ToLower(filepath.Ext(path)), ".mp4") {
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
	})
}

func PNGToHD() error {
	return filepath.Walk("images", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.EqualFold(strings.ToLower(filepath.Ext(path)), ".png") {
			return nil
		}

		if strings.HasSuffix(strings.ToLower(path), "_src.png") || strings.HasSuffix(strings.ToLower(path), "_dst.png") {
			return nil
		}

		// 拷贝1份原始文件
		srcPath := strings.ReplaceAll(path, ".png", "_src.png")
		if err := copyFile(path, srcPath); err != nil {
			return err
		}

		srcFile, err := os.Open(srcPath)
		if err != nil {
			return err
		}

		srcImage, err := png.Decode(srcFile)
		if err != nil {
			return err
		}

		if err := srcFile.Close(); err != nil {
			return err
		}

		bounds := srcImage.Bounds()
		log.Printf("%q : Width is %d,Height is %d \n", path, bounds.Dx(), bounds.Dy())
		for x := 0; x < bounds.Dx(); x++ {
			for y := 0; y < bounds.Dy(); y++ {

				//r1, g1, b1, a1 := originalImage.(*image.NRGBA).At(x-1, y-1).RGBA()
				//r2, g2, b2, a2 := originalImage.(*image.NRGBA).At(x, y-1).RGBA()
				//r3, g3, b3, a3 := originalImage.(*image.NRGBA).At(x+1, y-1).RGBA()
				//
				//r4, g4, b4, a4 := originalImage.(*image.NRGBA).At(x-1, y).RGBA()
				//r5, g5, b5, a5 := originalImage.(*image.NRGBA).At(x, y).RGBA()
				//r6, g6, b6, a6 := originalImage.(*image.NRGBA).At(x+1, y).RGBA()
				//
				//r7, g7, b7, a7 := originalImage.(*image.NRGBA).At(x-1, y+1).RGBA()
				//r8, g8, b8, a8 := originalImage.(*image.NRGBA).At(x, y+1).RGBA()
				//r9, g9, b9, a9 := originalImage.(*image.NRGBA).At(x+1, y+1).RGBA()
				//
				//if (a5>>8 == 86) || (a5>>8 == 171) {
				//	// 86 和 171 分别为残影的透明度
				//	continue
				//}
				//
				//a := (a1 + a2 + a3 + a4 + a5 + a6 + a7 + a8 + a9) / 9
				//if a == 0 {
				//	continue
				//}
				//
				//r := (r1 + r2 + r3 + r4 + r5 + r6 + r7 + r8 + r9) / 9
				//g := (g1 + g2 + g3 + g4 + g5 + g6 + g7 + g8 + g9) / 9
				//b := (b1 + b2 + b3 + b4 + b5 + b6 + b7 + b8 + b9) / 9
				//
				//if a >= 0xffff*3/4 {
				//	r, g, b, a = r5, g5, b5, a5
				//} else if a < 0xffff*1/4 {
				//	r, g, b, a = color.Transparent.RGBA()
				//}
				//
				newRGBA := color.RGBA{
					R: 255,
					G: 0,
					B: 0,
					A: 255,
				}

				srcImage.(*image.NRGBA).Set(x, y, newRGBA)
			}
		}

		dstPath := strings.ReplaceAll(path, ".png", "_dst.png")
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}

		if err := png.Encode(dstFile, srcImage); err != nil {
			return err
		}

		if err := dstFile.Close(); err != nil {
			return err
		}

		return nil
	})
}

func copyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	if err := dstFile.Close(); err != nil {
		return err
	}

	if err := srcFile.Close(); err != nil {
		return err
	}

	return nil
}
