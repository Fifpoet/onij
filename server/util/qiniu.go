package util

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"image"
	_ "image/gif"  // 必须导入
	_ "image/jpeg" // 必须导入
	_ "image/png"  // 必须导入
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"
)

const (
	ak = "lKrcOB9iUaKvKhspGh1hgo4-Dy0yFlH0mgRSgPRY"
	bk = "onij"
	dm = "http://cloud.onij.fun"
)

func getQiniuMac() *qbox.Mac {
	sk := os.Getenv("sk")
	return qbox.NewMac(ak, sk)
}

func UploadFile(localFilePath string) (string, error) {
	putPolicy := storage.PutPolicy{Scope: bk}
	upToken := putPolicy.UploadToken(getQiniuMac())

	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	err := formUploader.PutFile(context.Background(), &ret, upToken, uuid.New().String(), localFilePath, &putExtra)
	if err != nil {
		log.Printf("UploadFile, upload file failed: err = %v \n", err)
		return "", fmt.Errorf("file upload failed: %v", err)
	}
	fmt.Printf("File uploaded successfully, key: %s; hash: %s \n", ret.Key, ret.Hash)
	return ret.Key, nil
}

func DownloadFile(key string) string {
	deadline := time.Now().Add(time.Hour).Unix()

	privateAccessURL := storage.MakePrivateURL(getQiniuMac(), dm, key, deadline)
	fmt.Printf("Download URL: %s\n", privateAccessURL)
	return privateAccessURL
}

func DeleteFile(key string) error {
	mac := getQiniuMac()
	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)

	err := bucketManager.Delete(bk, key)
	if err != nil {
		log.Printf("DeleteFile, delete file failed: err = %v \n", err)
		return fmt.Errorf("file deletion failed: %v", err)
	}
	fmt.Println("File deleted successfully.")
	return nil
}

func UploadFromReader(file multipart.File, size int64) (string, error) {

	putPolicy := storage.PutPolicy{Scope: bk}
	upToken := putPolicy.UploadToken(getQiniuMac())

	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	err := formUploader.Put(context.Background(), &ret, upToken, uuid.New().String(), file, size, &putExtra)
	if err != nil {
		log.Printf("UploadFromReader, upload file failed: err = %v \n", err)
		return "", fmt.Errorf("file upload failed: %v", err)
	}
	fmt.Printf("File uploaded successfully, key: %s\n", ret.Key)
	return ret.Key, nil
}

func GetLocalFileHash(localFilePath string) (string, int, int, error) {
	file, err := os.Open(localFilePath)
	if err != nil && !errors.Is(err, image.ErrFormat) {
		log.Printf("GetLocalFileHash, open file failed: err = %v \n", err)
		return "", 0, 0, err
	}
	defer file.Close()
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Printf("GetLocalFileHash, unable to decode image: %v \n", err)
		return "", 0, 0, err
	}

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Printf("GetLocalFileHash, copy file failed: err = %v \n", err)
		return "", 0, 0, err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), img.Width, img.Height, nil
}

func GetFileHash(file multipart.File) (string, int, int, error) {
	hash := md5.New()
	f, _, err := image.DecodeConfig(file)
	if err != nil && !errors.Is(err, image.ErrFormat) {
		log.Printf("getFormXY, unable to decode image: %v \n", err)
		return "", 0, 0, err
	}

	if _, err := io.Copy(hash, file); err != nil {
		log.Printf("GetFileHash, copy file failed: err = %v \n", err)
		return "", 0, 0, err

	}

	return fmt.Sprintf("%x", hash.Sum(nil)), f.Width, f.Height, nil
}
