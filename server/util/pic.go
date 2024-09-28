package util

import (
	"image"
	_ "image/gif"  // 引入 GIF 格式支持
	_ "image/jpeg" // 引入 JPEG 格式支持
	_ "image/png"  // 引入 PNG 格式支持
	"log"
	"mime/multipart"
	"os"
)

func GetLocalXY(file *os.File) (int, int) {
	// 解码图片
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Printf("getLocalXY, unable to decode image: %v \n", err)
		return 0, 0
	}

	// 返回图片的宽和高
	return img.Width, img.Height
}

func GetFormXY(fileHeader *multipart.FileHeader) (int, int) {
	file, err := fileHeader.Open()
	defer file.Close()
	if err != nil {
		log.Printf("getFormXY, open file failed: err = %v \n", err)
		return 0, 0
	}
	f, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Printf("getFormXY, unable to decode image: %v \n", err)
		return 0, 0
	}
	return f.Width, f.Height
}
