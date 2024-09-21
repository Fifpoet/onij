package util

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"os"
	"time"
)

const (
	ak = "lKrcOB9iUaKvKhspGh1hgo4-Dy0yFlH0mgRSgPRY"
	bk = "onij"
	dm = "http://onij.s3.cn-east-1.qiniucs.com"
)

func getQiniuMac() *qbox.Mac {
	sk := os.Getenv("sk")
	return qbox.NewMac(ak, sk)
}

func UploadFile(localFilePath, key string) error {
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
		return fmt.Errorf("file upload failed: %v", err)
	}
	fmt.Printf("File uploaded successfully, key: %s\n", ret.Key)
	return nil
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
		Zone:          &storage.ZoneHuanan,
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
