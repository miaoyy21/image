package main

import (
	"log"
	"x/convert"
)

func main() {
	// MP4 转为 GIF
	if err := convert.MP4ToGIF(); err != nil {
		log.Printf("MP4ToGIF() Failure :: %s", err.Error())
		return
	}

	// 合并像素点
	//if err := convert.PNGToHD(); err != nil {
	//	log.Printf("PNGToHD() Failure :: %s", err.Error())
	//	return
	//}

	// 将 PNG 转为GIF
	if err := convert.PNGToGIF(); err != nil {
		log.Printf("PNGToGIF() Failure :: %s", err.Error())
		return
	}
}
