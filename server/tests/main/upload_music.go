package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/dhowden/tag"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/gorm"
	"io"
	"onij/biz/prm"
	"onij/infra"
	"onij/infra/mysql"
	"onij/inject"
	"onij/logic"
	"onij/model/api"
	"onij/util"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

const targetDir = "C:\\KwDownload\\song\\tmp"
const fileSuffix = ".mp3"
const lyricsSuffix = ".lrc"

var allInfra *infra.AllInfra
var allLogic *logic.AllLogic
var ctx context.Context

func uploadMusic() {
	allInfra = inject.InitDalForTest()
	allLogic = inject.InitLogicForTest()
	ctx = context.Background()

	mp3Files, err := findMP3Files(targetDir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Found MP3 files: %v \n", mp3Files)

	//var models []*mysql.Music
	var musicName string
	var artistNames []string
	var composerName string
	var writerName string
	for _, path := range mp3Files {
		mp3FileName := filepath.Base(path)
		splits := strings.Split(mp3FileName, "-")
		if len(splits) < 2 {
			fmt.Println("Error: Invalid file name:", path)
			continue
		}
		musicName = splits[1]
		singerSplits := strings.Split(splits[0], "&") // 分隔符待定
		artistNames = singerSplits

		searchMusic, err := allInfra.MusicDal.GetByFullName(mp3FileName)
		if err != nil {
			return
		}
		if len(searchMusic) > 0 {
			fmt.Printf("歌曲已存在: <<%s>>\n", mp3FileName)
			continue
		}

		writerName, composerName = parseArtistFromLyr(path + lyricsSuffix)
		fmt.Printf("<<%s>>艺术家: %s %s %s\n", musicName, artistNames, composerName, writerName)

		artistIds := make([]int64, 0)
		for _, artistName := range artistNames {
			artistId, err := processArtist(artistName, api.ArtistType_AT_Singer)
			if err != nil {
				return
			}
			artistIds = append(artistIds, artistId)
		}
		composerId, err := processArtist(composerName, api.ArtistType_AT_Composer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		writerId, err := processArtist(writerName, api.ArtistType_AT_Writer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("保存艺术家成功")

		mp3FileId, lycFileId, err := processFile(path)
		if err != nil {
			return
		}
		fmt.Println("保存文件成功")

		err = allInfra.MusicDal.Upsert(&mysql.Music{
			Id:          util.IdGen.Generate(),
			RootId:      0,
			Name:        musicName,
			FullName:    mp3FileName,
			ArtistIds:   util.Int64List2Str(artistIds),
			ComposerId:  composerId,
			WriterId:    writerId,
			IssueTime:   0,
			PerformType: 0,
			MvUrl:       "",
			Mp3FileId:   mp3FileId,
			LyricFileId: lycFileId,
		}, nil)
		if err != nil {
			fmt.Println("save music Error:", err)
			return
		}
		fmt.Println("保存音乐成功")
	}

}

func processFile(path string) (int64, int64, error) {
	mp3Name := filepath.Base(path + fileSuffix)
	lycName := filepath.Base(path + lyricsSuffix)
	mp3Bytes, err := readMediaFile(path + fileSuffix)
	if err != nil {
		return 0, 0, err
	}
	lycBytes, err := readMediaFile(path + lyricsSuffix)
	if err != nil {
		return 0, 0, err
	}
	mp3, err := allLogic.FileLogic.Upload(ctx, &prm.UploadFileParam{
		Filename: mp3Name,
		AppId:    0,
		Folders:  []string{"music"},
		File:     mp3Bytes,
	})
	if err != nil {
		return 0, 0, err
	}
	lyc, err := allLogic.FileLogic.Upload(ctx, &prm.UploadFileParam{
		Filename: lycName,
		AppId:    0,
		Folders:  []string{"music"},
		File:     lycBytes,
	})
	if err != nil {
		return 0, 0, err
	}
	return mp3.FileId, lyc.FileId, nil
}
func processArtist(name string, typ api.ArtistType) (int64, error) {
	artist, err := allInfra.ArtistDal.GetByName(name)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}
	if artist == nil {
		id := util.IdGen.Generate()
		err := allInfra.ArtistDal.Save(&mysql.Artist{
			Id:         id,
			Name:       name,
			ArtistType: int32(typ),
			DeletedAt:  gorm.DeletedAt{},
		})
		if err != nil {
			fmt.Println("save artists Error:", err)
			return 0, err
		}
		return id, nil
	} else {
		return artist.Id, nil
	}
}

func readMediaFile(filePath string) ([]byte, error) {
	rawData, err := os.ReadFile(filePath) // Go 1.16+ 推荐
	if err != nil {
		return nil, err
	}

	// 通过扩展名判断类型
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".lrc", ".txt":
		if utf8.Valid(rawData) {
			return rawData, nil
		}
		decoder := simplifiedchinese.GBK.NewDecoder()
		return io.ReadAll(transform.NewReader(bytes.NewReader(rawData), decoder))
	default: // mp3/mp4等二进制文件
		return rawData, nil // 直接返回原始数据
	}
}

func parseArtistFromLyr(path string) (string, string) {
	data, err := readMediaFile(path)
	if err != nil {
		return "", ""
	}

	var lyricist, composer string
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "词：") {
			// 提取词作者
			start := strings.Index(line, "词：") + len("词：")
			lyricist = strings.TrimSpace(line[start:])
		}

		if strings.Contains(line, "曲：") {
			// 提取曲作者
			start := strings.Index(line, "曲：") + len("曲：")
			composer = strings.TrimSpace(line[start:])
		}

		// 如果已经找到两个信息，提前退出循环
		if lyricist != "" && composer != "" {
			break
		}
	}

	return lyricist, composer
}

func findMP3Files(root string) ([]string, error) {
	var mp3Files []string

	// 定义正则表达式，用于匹配 "singer-name" 的文件名结构
	re := regexp.MustCompile(`^.+-.+\.mp3$`)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(info.Name()) == ".mp3" {
			// 使用正则表达式匹配文件名是否符合 "singer-name" 的结构
			if re.MatchString(info.Name()) {
				mp3Files = append(mp3Files, path[:len(path)-4])
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return mp3Files, err
}

func getMp3Meta(file *os.File) int {
	meta, err := tag.ReadFrom(file)
	if err != nil {
		fmt.Println("Error reading tags:", err)
		return 0
	}
	sec := meta.Raw()["length"]
	if sec == nil {
		fmt.Println("No file second found")
		return 0
	}
	return sec.(int)
}
