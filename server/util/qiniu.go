package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
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

// PublicDomain 七牛桶绑定的 CDN 域名（须为 cloud.onij.fun 等对象域，勿填主站 onij.fun）
func PublicDomain() string {
	if v := strings.TrimSpace(os.Getenv("QINIU_PUBLIC_DOMAIN")); v != "" {
		return strings.TrimRight(v, "/")
	}
	return dm
}

// FetchDomain 服务端拉取对象用的域名（本地开发 cloud.onij.fun 可能无法解析）
func FetchDomain() string {
	if v := strings.TrimSpace(os.Getenv("QINIU_FETCH_DOMAIN")); v != "" {
		return strings.TrimRight(v, "/")
	}
	return "http://onij.bkt.clouddn.com"
}

var zone = &storage.ZoneHuadong

// GetUploadURL 七牛上传域名（与 ZoneHuadong 一致，不能用 cloud.onij.fun 下载域）
func GetUploadURL() string {
	if zone != nil && len(zone.SrcUpHosts) > 0 {
		return "https://" + zone.SrcUpHosts[0]
	}
	return "https://up-z0.qiniup.com"
}

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

	r, err := bucketManager.Get(bk, key, nil)
	if err != nil {
		log.Printf("GetFile, get file failed: err = %v \n", err)
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func MimeTypeByName(name string) string {
	switch strings.ToLower(path.Ext(name)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".svg":
		return "image/svg+xml"
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mp3":
		return "audio/mpeg"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
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
	var ossPath string
	if len(info.OssFolder) > 0 {
		ossPath = strings.Join(info.OssFolder, "/") + "/" + info.Name
	} else {
		ossPath = info.Name
	}

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
	return BrowserFileURL(key, "")
}

// BrowserFileURL 浏览器跟随的私有链（须用桶绑定域名，bkt.clouddn.com 会 404）
func BrowserFileURL(key string, attname string) string {
	return privateFileURL(key, PublicDomain(), attname)
}

// ResolvePrivateFileURL 服务端拉取对象用的私有链
func ResolvePrivateFileURL(key string, attname string) string {
	return privateFileURL(key, FetchDomain(), attname)
}

func privateFileURL(key string, domain string, attname string) string {
	if len(key) == 0 {
		return ""
	}
	deadline := time.Now().Add(time.Hour).Unix()
	url := storage.MakePrivateURL(getQiniuMac(), domain, key, deadline)
	if attname == "" {
		return url
	}
	sep := "?"
	if strings.Contains(url, "?") {
		sep = "&"
	}
	return url + sep + "attname=" + urlEncode(attname)
}

func urlEncode(s string) string {
	return strings.ReplaceAll(url.QueryEscape(s), "+", "%20")
}

// FetchObject 服务端经私有链拉取对象（不依赖 cloud.onij.fun 本地 DNS）
func FetchObject(key string) ([]byte, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("empty key")
	}
	deadline := time.Now().Add(time.Hour).Unix()
	url := storage.MakePrivateURL(getQiniuMac(), FetchDomain(), key, deadline)
	data, err := GET(url, nil)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("empty object")
	}
	return data, nil
}

func DeleteFile(key string) error {
	bucketManager := getManager()

	err := bucketManager.Delete(bk, key)
	if err != nil {
		if isQiniuNoSuchKey(err) {
			return nil
		}
		log.Printf("DeleteFile, delete file failed: err = %v \n", err)
		return fmt.Errorf("file deletion failed: %v", err)
	}
	return nil
}

func isQiniuNoSuchKey(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "no such file") ||
		strings.Contains(msg, "612") ||
		strings.Contains(msg, "not found")
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
