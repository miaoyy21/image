package convert

import (
	"fmt"
	"github.com/xyproto/palgen"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func PNGToGIF() error {
	pngNames := make([]string, 0)

	if err := filepath.Walk("images", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.EqualFold(strings.ToLower(filepath.Ext(path)), ".png") {
			return nil
		}

		if !strings.HasPrefix(strings.ToLower(info.Name()), "dst_") {
			return nil
		}

		pngNames = append(pngNames, path)
		return nil
	}); err != nil {
		return err
	}

	gifNames := make(map[string][]string)
	for _, pngName := range pngNames {
		names := strings.Split(pngName, "_")

		gifName := fmt.Sprintf("%s.gif", strings.Join(names[:len(names)-1], "_"))
		pngs, ok := gifNames[gifName]
		if !ok {
			pngs = make([]string, 0)
		}

		pngs = append(pngs, pngName)

		// PNG图片排序
		sort.Slice(pngs, func(i, j int) bool {
			iIndex, err := strconv.Atoi(pngs[i][strings.LastIndex(pngs[i], "_")+1 : len(pngs[i])-4])
			if err != nil {
				log.Panic(err)
			}

			jIndex, err := strconv.Atoi(pngs[j][strings.LastIndex(pngs[j], "_")+1 : len(pngs[j])-4])
			if err != nil {
				log.Panic(err)
			}

			return iIndex < jIndex
		})

		gifNames[gifName] = pngs
	}

	for gifName, pngNames := range gifNames {
		log.Printf("GIF[ %s ] :: PNG[ %s ] ...\n", gifName, strings.Join(pngNames, ","))

		out := gif.GIF{LoopCount: len(pngNames)}
		for index, pngName := range pngNames {
			pngFile, err := os.Open(pngName)
			if err != nil {
				return err
			}

			pngImage, err := png.Decode(pngFile)
			if err != nil {
				return err
			}

			palPalette, err := palgen.Generate(pngImage, 256)
			if err != nil {
				return err
			}

			paletteImg := image.NewPaletted(pngImage.Bounds(), palPalette)

			draw.Draw(paletteImg, pngImage.Bounds(), pngImage, image.ZP, draw.Over)

			out.Image = append(out.Image, paletteImg)
			out.Delay = append(out.Delay, 15)

			if err := pngFile.Close(); err != nil {
				return err
			}

			log.Printf("GIF[ %s ] Dealing:: [ %d/%d ] ...\n", gifName, index, len(pngNames))
		}

		gifFile, err := os.Create(gifName)
		if err != nil {
			return err
		}

		if err := gif.EncodeAll(gifFile, &out); err != nil {
			return err
		}

		if err := gifFile.Close(); err != nil {
			return err
		}
	}

	return nil
}
