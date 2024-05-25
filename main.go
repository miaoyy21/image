package main

import (
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"image"
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

		originalFile, err := os.Open(path)
		if err != nil {
			return err
		}

		originalImage, _, err := image.Decode(originalFile)
		if err != nil {
			return fmt.Errorf("change Definition image.Decode [%s] : %s", path, err.Error())
		}

		if err := originalFile.Close(); err != nil {
			return err
		}

		bounds := originalImage.Bounds()
		log.Printf("%q : Width is %d,Height is %d \n", info.Name(), bounds.Dx(), bounds.Dy())
		for x := 1; x < bounds.Dx()-1; x++ {
			for y := 1; y < bounds.Dy()-1; y++ {

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
				//newRGBA := color.RGBA{
				//	R: uint8(r >> 8),
				//	G: uint8(g >> 8),
				//	B: uint8(b >> 8),
				//	A: uint8(a >> 8),
				//}
				//
				//originalImage.(*image.NRGBA).Set(x, y, newRGBA)
			}
		}

		//newFile, err := os.Create(path)
		//if err != nil {
		//	return fmt.Errorf("change Definition os.Create [%s] : %s", path, err.Error())
		//}
		//
		//if err := png.Encode(newFile, originalImage); err != nil {
		//	return fmt.Errorf("change Definition png.Encode [%s] : %s", path, err.Error())
		//}
		//
		//if err := newFile.Close(); err != nil {
		//	return err
		//}

		return nil
	})
}
