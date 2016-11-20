package main

import (
	"flag"

	"./zm/decodeqstp"
	"./zm/split"
	"./zm/zmmd5"
)

func main() {
	useSplit := flag.Bool("s", false, "分解暴风包")
	useMd5 := flag.Bool("m", false, "计算某个文件夹内的md5")
	useDecode := flag.Bool("d", false, "解码qstp串")

	path := flag.String("p", "", "参数 文件或文件夹路径")
	flag.Parse()

	use := *useSplit || *useMd5 || *useDecode
	if path == nil || *path == "" || !use {
		flag.PrintDefaults()
		return
	}

	if *useSplit {
		split.Split(path)
		return
	}

	if *useMd5 {
		zmmd5.CalMd5(path)
		return
	}

	if *useDecode {
		decodeqstp.Decode(path)
		return
	}
}
