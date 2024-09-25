package main

import (
	"fmt"
	"github.com/dhowden/tag"
	"onij/infra/mysql"
	"onij/logic"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const targetDir = "/Users/asen/Documents/gopath/src/owner/onij/server/static/kuwo"
const fileSuffix = ".mp3"
const lyricsSuffix = ".lrc"

func main() {
	logic.Init()

	mp3Files, err := findMP3Files(targetDir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Found MP3 files: %v \n", mp3Files)

	var models []*mysql.Music
	var lyricsFiles []string
	var f *os.File
	for _, mp3File := range mp3Files {
		f, err = os.Open(mp3File)
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}

		fileName := filepath.Base(mp3File)
		split := strings.Split(fileName, "-")
		singer := split[0]
		songName := split[1]

		// 构建歌词文件路径
		lyricsFile := strings.Replace(mp3File, ".mp3", ".lrc", -1)
		if _, err = os.Stat(lyricsFile); os.IsNotExist(err) {
			lyricsFile = ""
		}
		lyricsFiles = append(lyricsFiles, lyricsFile)

		models = append(models, &mysql.Music{
			RootId:      0,
			Title:       songName,
			ArtistIds:   singer,
			Composer:    0,
			Writer:      0,
			Length:      getMp3Meta(f),
			IssueYear:   0,
			Language:    "",
			PerformType: 0,
			Instrument:  "",
			Concert:     "",
			ConcertYear: 0,
			Sequence:    0,
			MvUrl:       "",
		})
	}
	defer f.Close()

	err = logic.NewMusicLogic().SaveFromDir(models, mp3Files, lyricsFiles)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// findMP3Files 遍历目录并返回符合条件的 .mp3 文件路径
func findMP3Files(root string) ([]string, error) {
	var mp3Files []string

	// 定义正则表达式，用于匹配 "singer-name" 的文件名结构
	re := regexp.MustCompile(`^[\p{Han}a-zA-Z]+-[\p{Han}a-zA-Z]+\.mp3$`)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".mp3" {
			// 使用正则表达式匹配文件名是否符合 "singer-name" 的结构
			if re.MatchString(info.Name()) {
				mp3Files = append(mp3Files, path)
			}
		}
		return nil
	})

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
