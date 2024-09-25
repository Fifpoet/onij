package util

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
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

func UploadFile(localFilePath, key string) (string, error) {
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
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFilePath, &putExtra)
	if err != nil {
		return "", fmt.Errorf("file upload failed: %v", err)
	}
	fmt.Printf("File uploaded successfully, key: %s\n", ret.Key)

	url := DownloadFile(key)
	return url, nil
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
		return fmt.Errorf("file deletion failed: %v", err)
	}
	fmt.Println("File deleted successfully.")
	return nil
}

func UploadFromReader(file multipart.File, size int64, key string) error {
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

	err := formUploader.Put(context.Background(), &ret, upToken, key, file, size, &putExtra)
	if err != nil {
		return fmt.Errorf("file upload failed: %v", err)
	}
	fmt.Printf("File uploaded successfully, key: %s\n", ret.Key)
	return nil
}
