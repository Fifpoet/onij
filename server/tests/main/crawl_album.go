package main

import (
	"context"
	"onij/inject"
	"onij/util/boost/tool"
	"onij/util/logs"
	"os"
)
var ctx = context.Background()
var app *inject.App

func init () {
	app = inject.InitializeApp()
}

// http请求, 获取album json, 反序列化到album结构体
func GetOriginAlbum(albumIds ...int64) (*AlbumResponse, error) {
	file, err := os.ReadFile("album.json")
	if err != nil {
		return nil, err
	}
	album := tool.LoadJson[AlbumResponse](string(file), true)
	
	// 处理album.artists
	album.Album.Artists

	return album, nil
}
