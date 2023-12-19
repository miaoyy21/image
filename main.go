package main

import (
	"flag"
	"log"
	"x/download"
	"x/split"
	"x/srt"
)

var sMode string
var sFile string
var sDir string

func init() {
	flag.StringVar(&sMode, "m", "srt", "使用模式：[download]自动下载图片；[split]自动对文本按句分隔；[srt]自动提取字幕文件")
	flag.StringVar(&sFile, "f", "/Users/miaojingyi/Documents/media/short/Product/最强太子妃/最强太子妃 1.srt", "需要自动下载的文件或字幕文件")
	flag.StringVar(&sDir, "d", "", "需要自动换替换的文件目录")

	flag.Parse()
}

func main() {

	switch sMode {
	case "download":
		if len(sFile) < 1 {
			log.Println("需要指定自动下载的文件")
			flag.Usage()
			return
		}

		if err := download.Download(sFile); err != nil {
			log.Fatalf("Download Failure :: %s \n", err.Error())
		}
	case "split":
		if len(sDir) < 1 {
			log.Println("需要指定替换的的文件目录")
			flag.Usage()
			return
		}

		if err := split.Split(sDir); err != nil {
			log.Fatalf("Split Failure :: %s \n", err.Error())
		}
	case "srt":
		if len(sFile) < 1 {
			log.Println("需要指定提取的字幕文件")
			flag.Usage()
			return
		}

		if err := srt.Srt(sFile); err != nil {
			log.Fatalf("Srt Failure :: %s \n", err.Error())
		}
	default:
		flag.Usage()
	}
}
