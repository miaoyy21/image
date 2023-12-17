package main

import (
	"flag"
	"log"
	"x/download"
	"x/swap"
)

var sMode string
var sFile string
var sDir string

func init() {
	flag.StringVar(&sMode, "m", "download", "使用模式：[download]自动下载图片；[swap]自动换脸")
	flag.StringVar(&sFile, "f", "/Users/miaojingyi/Documents/media/women/002", "需要自动下载的文件")
	flag.StringVar(&sDir, "d", "", "需要自动换脸的文件目录")

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
	case "swap":
		if len(sDir) < 1 {
			log.Println("需要指定自动换脸的文件目录")
			flag.Usage()
			return
		}

		if err := swap.Swap(sFile); err != nil {
			log.Fatalf("Swap Failure :: %s \n", err.Error())
		}
	default:
		flag.Usage()
	}
}
