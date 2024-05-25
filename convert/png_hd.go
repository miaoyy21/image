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

func PNGToHD() error {
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
		dstImage := image.NewNRGBA(
			image.Rectangle{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{X: bounds.Dx() / 2, Y: bounds.Dy() / 2},
			},
		)

		for x := 0; x < bounds.Dx(); x = x + 2 {
			for y := 0; y < bounds.Dy(); y = y + 2 {
				r1, g1, b1, a1 := srcImage.(*image.NRGBA).At(x, y).RGBA()
				r2, g2, b2, a2 := srcImage.(*image.NRGBA).At(x+1, y).RGBA()
				r3, g3, b3, a3 := srcImage.(*image.NRGBA).At(x, y+1).RGBA()
				r4, g4, b4, a4 := srcImage.(*image.NRGBA).At(x+1, y+1).RGBA()

				rgba1 := r1>>8<<24 | g1>>8<<16 | b1>>8<<8 | a1>>8
				rgba2 := r2>>8<<24 | g2>>8<<16 | b2>>8<<8 | a2>>8
				rgba3 := r3>>8<<24 | g3>>8<<16 | b3>>8<<8 | a3>>8
				rgba4 := r4>>8<<24 | g4>>8<<16 | b4>>8<<8 | a4>>8

				colorsMap := make(map[uint32]int)
				colorsMap[rgba1]++
				colorsMap[rgba2]++
				colorsMap[rgba3]++
				colorsMap[rgba4]++

				var exists bool

				for c := range colorsMap {
					r, g, b := uint8(c>>24), uint8(c&0x00ff0000>>16), uint8(c&0x0000ff00>>8)
					if r>>4 == g>>4 && r>>4 == b>>4 && r>>4 == 0x0f {
						exists = true
					}
				}

				var r, g, b, a uint8

				if exists {
					r, g, b, a = 0xff, 0xff, 0xff, 0xff
				} else {
					switch len(colorsMap) {
					case 4:
						// 四种颜色
						r, g, b, a = uint8(r4>>8), uint8(g4>>8), uint8(b4>>8), uint8(a4>>8)
					case 3:
						// 三种颜色，必有一种相同
						for c, n := range colorsMap {
							if n != 2 {
								continue
							}

							r, g, b, a = uint8(c>>24), uint8(c&0x00ff0000>>16), uint8(c&0x0000ff00>>8), uint8(c&0x000000ff)
						}
					case 2:
						// 两种颜色，可能是A颜色1个和B颜色3个，或者是A颜色2个和B颜色2个
						colorsSlice := make([]int, 0, 2)
						for _, count := range colorsMap {
							colorsSlice = append(colorsSlice, count)
						}

						if colorsSlice[0] == colorsSlice[1] {
							// A颜色2个和B颜色2个
							r, g, b, a = uint8(r4>>8), uint8(g4>>8), uint8(b4>>8), uint8(a4>>8)
						} else {
							// A颜色1个和B颜色3个
							for c, n := range colorsMap {
								if n != 3 {
									continue
								}

								r, g, b, a = uint8((c&0xff000000)>>24), uint8((c&0x00ff0000)>>16), uint8((c&0x0000ff00)>>8), uint8(c&0x000000ff)
							}
						}
					case 1:
						// 一种颜色
						r, g, b, a = uint8(r4>>8), uint8(g4>>8), uint8(b4>>8), uint8(a4>>8)
					default:
						r, g, b, a = 0, 0, 0, 0
					}
				}

				newRGBA := color.NRGBA{R: r, G: g, B: b, A: a}
				dstImage.SetNRGBA(x/2, y/2, newRGBA)
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
