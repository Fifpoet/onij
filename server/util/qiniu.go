package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
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

const (
	hlsSegFop        = "avthumb/m3u8/noDomain/1/segtime/10/vcodec/copy/acodec/copy"
	faststartFop     = "avthumb/mp4/vcodec/copy/acodec/copy"
	headProbeBytes   = 262144
	qiniuHTTPTimeout = 30 * time.Second
)

func qiniuPipeline() string {
	return strings.TrimSpace(os.Getenv("QINIU_PIPELINE"))
}

func newOperationManager() *storage.OperationManager {
	cfg := storage.Config{
		Zone:          zone,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	return storage.NewOperationManager(getQiniuMac(), &cfg)
}

// HLSKey 约定：原对象 key 后追加 .m3u8
func HLSKey(storeKey string) string {
	return storeKey + ".m3u8"
}

func privateFopURL(key, domain, fop string) string {
	if len(key) == 0 {
		return ""
	}
	deadline := time.Now().Add(time.Hour).Unix()
	publicURL := storage.MakePublicURL(domain, key)
	urlToSign := fmt.Sprintf("%s?%s&e=%d", publicURL, fop, deadline)
	token := getQiniuMac().Sign([]byte(urlToSign))
	return urlToSign + "&token=" + token
}

func qiniuHTTPGet(rawURL string, extra http.Header) ([]byte, int, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, 0, err
	}
	for k, vs := range extra {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	client := &http.Client{Timeout: qiniuHTTPTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return data, resp.StatusCode, nil
}

type AvinfoFormat struct {
	Duration   string `json:"duration"`
	Size       string `json:"size"`
	FormatName string `json:"format_name"`
}

type AvinfoResult struct {
	Format AvinfoFormat `json:"format"`
}

func (a *AvinfoResult) DurationMs() int64 {
	if a == nil {
		return 0
	}
	sec, err := strconv.ParseFloat(strings.TrimSpace(a.Format.Duration), 64)
	if err != nil || sec <= 0 {
		return 0
	}
	return int64(sec * 1000)
}

// Avinfo 拉七牛音视频元信息（服务端走 FetchDomain）
func Avinfo(key string) (*AvinfoResult, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("empty key")
	}
	raw := privateFopURL(key, FetchDomain(), "avinfo")
	data, code, err := qiniuHTTPGet(raw, nil)
	if err != nil {
		return nil, err
	}
	if code >= 400 {
		return nil, fmt.Errorf("avinfo http %d: %s", code, strings.TrimSpace(string(data)))
	}
	var ret AvinfoResult
	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

// ObjectExists 对象是否在桶中（612 视为不存在）
func ObjectExists(key string) bool {
	if len(key) == 0 {
		return false
	}
	_, err := getManager().Stat(bk, key)
	return err == nil
}

// ObjectHeadHasMoovFirst 读对象头部，判断 moov 是否在 mdat 之前（faststart）
func ObjectHeadHasMoovFirst(key string) (bool, error) {
	if len(key) == 0 {
		return false, fmt.Errorf("empty key")
	}
	raw := privateFileURL(key, FetchDomain(), "")
	h := http.Header{}
	h.Set("Range", fmt.Sprintf("bytes=0-%d", headProbeBytes-1))
	data, code, err := qiniuHTTPGet(raw, h)
	if err != nil {
		return false, err
	}
	if code >= 400 && code != http.StatusPartialContent {
		return false, fmt.Errorf("head probe http %d", code)
	}
	moov := bytes.Index(data, []byte("moov"))
	mdat := bytes.Index(data, []byte("mdat"))
	if moov >= 0 && (mdat < 0 || moov < mdat) {
		return true, nil
	}
	return false, nil
}

func Pfop(key, fops string) (string, error) {
	om := newOperationManager()
	pid, err := om.Pfop(bk, key, fops, qiniuPipeline(), "", false)
	if err != nil {
		log.Printf("Pfop failed: key=%s err=%v\n", key, err)
		return "", err
	}
	return pid, nil
}

func Prefop(persistentID string) (storage.PrefopRet, error) {
	om := newOperationManager()
	return om.Prefop(persistentID)
}

// PfopFaststart 把 moov 挪到文件头，覆盖原 key
func PfopFaststart(key string) (string, error) {
	saveas := storage.EncodedEntry(bk, key)
	return Pfop(key, faststartFop+"|saveas/"+saveas)
}

// PfopHLS 单码率 m3u8，输出 key 为 HLSKey(src)
func PfopHLS(srcKey string) (string, error) {
	m3u8Key := HLSKey(srcKey)
	saveas := storage.EncodedEntry(bk, m3u8Key)
	return Pfop(srcKey, hlsSegFop+"|saveas/"+saveas)
}

func PrefopBusy(code int) bool {
	return code == 1 || code == 2
}

func PrefopOK(code int) bool {
	return code == 0
}

// SignKeyURL 用桶绑定域签发私有链（浏览器经 /cdn 反代）
func SignKeyURL(key string) string {
	return BrowserFileURL(key, "")
}

// RewriteM3U8 把清单里的 ts 行改成签名 URL（文本，不是视频 body）
func RewriteM3U8(playlist, m3u8Key string) string {
	dir := path.Dir(m3u8Key)
	lines := strings.Split(strings.ReplaceAll(playlist, "\r\n", "\n"), "\n")
	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasPrefix(trim, "#") {
			continue
		}
		segKey := trim
		if strings.HasPrefix(trim, "http://") || strings.HasPrefix(trim, "https://") {
			u, err := url.Parse(trim)
			if err == nil {
				segKey = strings.TrimPrefix(u.Path, "/")
			}
		} else if dir != "." && dir != "/" && !strings.Contains(trim, "/") {
			segKey = dir + "/" + trim
		}
		lines[i] = SignKeyURL(segKey)
	}
	return strings.Join(lines, "\n")
}

// FetchM3U8Text 拉 m3u8 清单文本
func FetchM3U8Text(m3u8Key string) (string, error) {
	data, err := FetchObject(m3u8Key)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
