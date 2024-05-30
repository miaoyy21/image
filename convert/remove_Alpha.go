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

func RemoveAlpha() error {
	return filepath.Walk("images", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dir, name := filepath.Split(path)

		if !strings.EqualFold(strings.ToLower(filepath.Ext(path)), ".png") {
			return nil
		}

		if strings.HasPrefix(strings.ToLower(name), "dst_") {
			return nil
		}

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
				//if a>>8 != 0xff {
				//	dstImage.SetNRGBA(x, y, color.NRGBA{R: 0, G: 0, B: 0, A: 0})
				//} else {

				var r0, g0, b0 uint8

				//r0, g0, b0 = 156, 39, 176
				//
				//r0 = uint8(float64(r0) * float64(r>>12) / float64(16))
				//g0 = uint8(float64(g0) * float64(g>>12) / float64(16))
				//b0 = uint8(float64(b0) * float64(b>>12) / float64(16))

				//if r>>8 <= 0x1f {
				//	r0, g0, b0 = 0x00, 0x00, 0x00
				//} else if r>>8 <= 0x3f {
				//	r0, g0, b0 = 0x20, 0x20, 0x20
				//} else if r>>8 <= 0x5f {
				//	r0, g0, b0 = 0x40, 0x40, 0x40
				//} else if r>>8 <= 0x7f {
				//	r0, g0, b0 = 0x60, 0x60, 0x60
				//} else if r>>8 <= 0x9f {
				//	r0, g0, b0 = 0x80, 0x80, 0x80
				//} else if r>>8 <= 0xbf {
				//	r0, g0, b0 = 0xa0, 0xa0, 0xa0
				//} else if r>>8 <= 0xdf {
				//	r0, g0, b0 = 0xc0, 0xc0, 0xc0
				//} else {
				//	r0, g0, b0 = 0xe0, 0xe0, 0xe0
				//}

				//if r>>8 <= 0x3f {
				//	r0, g0, b0 = 53, 14, 71
				//} else if r>>8 <= 0x5f {
				//	r0, g0, b0 = 83, 15, 91
				//} else if r>>8 <= 0x7f {
				//	r0, g0, b0 = 106, 38, 157
				//} else if r>>8 <= 0x9f {
				//	r0, g0, b0 = 111, 52, 147
				//} else {
				//	r0, g0, b0 = 152, 87, 205
				//}

				//if r0 <= 0x3f {
				//	r0 = 0x20
				//} else if r0 <= 0x7f {
				//	r0 = 0x60
				//} else if r0 <= 0xbf {
				//	r0 = 0xa0
				//} else {
				//	r0 = 0xe0
				//}

				_, _ = g, b
				r0, g0, b0 = uint8(r>>12<<4), uint8(r>>12<<4), uint8(r>>12<<4)
				if r>>8 <= 0x1f {
					r0, g0, b0 = 57, 19, 80
				} else if r>>8 <= 0x2f {
					r0, g0, b0 = 77, 19, 94
				} else if r>>8 <= 0x3f {
					r0, g0, b0 = 126, 46, 178
				} else if r>>8 <= 0x4f {
					r0, g0, b0 = 138, 78, 178
				} else if r>>8 <= 0x5f {
					r0, g0, b0 = 152, 87, 205
				} else if r>>8 <= 0x6f {
					r0, g0, b0 = 164, 92, 221
				} else if r>>8 <= 0x7f {
					r0, g0, b0 = 119, 63, 178
				} else if r>>8 <= 0x8f {
					r0, g0, b0 = 161, 97, 221
				} else if r>>8 <= 0x9f {
					r0, g0, b0 = 152, 87, 205
				} else {
					r0, g0, b0 = 126, 46, 178
				}

				newRGBA := color.NRGBA{R: r0, G: g0, B: b0, A: uint8(a >> 8)}
				dstImage.SetNRGBA(x, y, newRGBA)
				//}
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
