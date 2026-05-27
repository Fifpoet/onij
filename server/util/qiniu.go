package util

import (
	"bytes"
	"context"
	"fmt"
	_ "image/gif"  // 必须导入
	_ "image/jpeg" // 必须导入
	_ "image/png"  // 必须导入
	"log"
	"os"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

const (
	ak = "lKrcOB9iUaKvKhspGh1hgo4-Dy0yFlH0mgRSgPRY"
	bk = "onij"
	dm = "http://cloud.onij.fun"
)

var zone = &storage.ZoneHuadong

type UploadInfo struct {
	Name      string
	Url       string
	Bytes     []byte
	OssFolder []string
}

func getQiniuSK() string {
	for _, key := range []string{"QINIU_SK", "sk", "SK"} {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			return v
		}
	}
	panic(
		"七牛 SecretKey 未注入到 onij-server 进程：" +
			"请设置环境变量 QINIU_SK（或 sk）。" +
			"仅在 ssh 里 export 对 systemd/nohup 启动的进程无效，需在服务配置里写 Environment=",
	)
}

func getQiniuMac() *qbox.Mac {
	return qbox.NewMac(ak, getQiniuSK())
}

func GetFile(key string) ([]byte, error) {
	bucketManager := getManager()

	var res []byte
	r, err := bucketManager.Get(bk, key, nil)
	if err != nil {
		log.Printf("GetFiles, get file failed: err = %v \n", err)
		return nil, err
	}
	_, err = r.Body.Read(res)
	if err != nil {
		log.Printf("GetFiles, read file failed: err = %v \n", err)
		return nil, err
	}
	return res, nil
}

func UploadFile(ctx context.Context, info UploadInfo) (string, error) {
	putPolicy := storage.PutPolicy{Scope: bk}
	upToken := putPolicy.UploadToken(getQiniuMac())
	cfg := storage.Config{
		Zone:          zone,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	var err error
	ossPath := strings.Join(info.OssFolder, "/") + "/" + info.Name

	var bs []byte
	if len(info.Bytes) > 0 {
		bs = info.Bytes
	} else {
		bs, err = GET(info.Url, nil)
		if err != nil {
			return "", err
		}
	}
	err = formUploader.Put(ctx, &ret, upToken, ossPath, bytes.NewReader(bs), int64(len(bs)), &putExtra)
	if err != nil {
		return "", err
	}
	return ret.Key, nil
}

func DownloadFile(key string) string {
	if len(key) == 0 {
		return ""
	}
	deadline := time.Now().Add(time.Hour).Unix()

	privateAccessURL := storage.MakePrivateURL(getQiniuMac(), dm, key, deadline)
	return privateAccessURL
}

func DeleteFile(key string) error {
	bucketManager := getManager()

	err := bucketManager.Delete(bk, key)
	if err != nil {
		log.Printf("DeleteFile, delete file failed: err = %v \n", err)
		return fmt.Errorf("file deletion failed: %v", err)
	}
	fmt.Println("File deleted successfully.")
	return nil
}

func getManager() *storage.BucketManager {
	mac := getQiniuMac()
	cfg := storage.Config{
		Zone:          zone,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)

	return bucketManager
}

// 获取上传token
func GetUploadToken() string {
	putPolicy := storage.PutPolicy{Scope: bk}
	upToken := putPolicy.UploadToken(getQiniuMac())
	return upToken
}
