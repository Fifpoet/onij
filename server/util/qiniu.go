package util

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	_ "image/gif"  // 必须导入
	_ "image/jpeg" // 必须导入
	_ "image/png"  // 必须导入
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

const (
	ak = "lKrcOB9iUaKvKhspGh1hgo4-Dy0yFlH0mgRSgPRY"
	bk = "onij"
	dm = "http://cloud.onij.fun"
)

type UploadInfo struct {
	Name string

	// 文件信息, 三选一
	LocalPath string

	Bytes []byte

	File multipart.File
	Size int64

	// extra
	OssFolder []string
}

func getQiniuMac() *qbox.Mac {
	sk := os.Getenv("sk")
	if sk == "" {
		panic("sk is empty")
	}
	return qbox.NewMac(ak, sk)
}

func UploadFile(ctx context.Context, info UploadInfo) (string, error) {
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

	var err error
	ossPath := strings.Join(info.OssFolder, "/") + "/" + info.Name + "-" + uuid.New().String()[:8]
	if info.LocalPath != "" {
		err = formUploader.PutFile(ctx, &ret, upToken, ossPath, info.LocalPath, &putExtra)
	} else if len(info.Bytes) > 0 {
		err = formUploader.Put(ctx, &ret, upToken, ossPath, bytes.NewReader(info.Bytes), int64(len(info.Bytes)), &putExtra)
	} else if info.File != nil {
		err = formUploader.Put(ctx, &ret, upToken, ossPath, info.File, info.Size, &putExtra)
	} else {
		return "", errors.New("no file to upload")
	}
	if err != nil {
		log.Printf("UploadFile, upload file failed: err = %v \n", err)
		return "", fmt.Errorf("file upload failed: %v", err)
	}
	fmt.Printf("File uploaded successfully, key: %s; hash: %s \n", ret.Key, ret.Hash)
	return ret.Key, nil
}

func DownloadFile(key string) string {
	if len(key) == 0 {
		return ""
	}
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
