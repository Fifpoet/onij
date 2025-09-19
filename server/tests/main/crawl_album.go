package main

import (
	"context"
	"onij/biz/errdef"
	"onij/biz/prm"
	"onij/inject"
	"onij/util"
	"onij/util/boost/exp"
	"onij/util/boost/tool"
	"os"
)

var ctx = context.Background()
var app *inject.App

const ArtistAvatarFolderId = 11111111111
const AlbumCoverFolderId = 11111111111

func init() {
	app = inject.InitializeApp()
}

func FetchAlbum(albumIds ...int64) (*AlbumResponse, error) {
	file, err := os.ReadFile("album.json")
	if err != nil {
		return nil, err
	}
	album := tool.LoadJson[AlbumResponse](string(file), true)

	// 处理album.artists
	var artistIds []int64
	for _, albumArtist := range album.Album.Artists {
		aId, err := processArtist(albumArtist.Name, albumArtist.Img1v1URL)
		if err != nil {
			return nil, err
		}
		artistIds = append(artistIds, aId)
	}

	// 处理album
	albumId, err := processAlbum(album.Album, album.Songs, artistIds)
	if err != nil {
		return nil, err
	}

	// 处理album.songs

	return album, nil
}

func processArtist(name string, avatarUrl string) (int64, error) {
	artistsResp, err := app.ArtistLogic.GetList(ctx, &prm.GetArtistListParam{
		Name:  name,
		Page:  1,
		Limit: 1,
	})
	if err != nil {
		return 0, err
	}
	if artistsResp.Total > 0 {
		return 0, nil
	}

	fileUploadRes, err := app.FileLogic.Upload(ctx, &prm.UploadFileParam{
		ParentId: ArtistAvatarFolderId,
		Url:      exp.Ptr(avatarUrl),
	})
	if err != nil {
		return 0, err
	}
	_, err = app.AllLogic.ArtistLogic.Upload(ctx, &prm.UploadArtistParam{
		Name:         name,
		AvatarFileId: exp.Ptr(fileUploadRes.FileIds[0]),
		ArtistId:     nil,
	})
	if err != nil {
		return 0, err
	}
	return fileUploadRes.FileIds[0], nil
}

func processAlbum(album Album, songs []Song, artistIds []int64) (int64, error) {
	_, cnt, err := app.AlbumDal.GetByNameAndArtistId(ctx, "", album.Name, artistIds, util.Page{1, 1})
	if cnt > 0 {
		return 0, errdef.ErrAlbumExisted
	}
	fileUploadRes, err := app.FileLogic.Upload(ctx, &prm.UploadFileParam{
		ParentId: ArtistAvatarFolderId,
		Url:      exp.Ptr(album.PicURL),
	})
	if err != nil {
		return 0, err
	}

	// album.songs
	for _, song := range songs {

	}

	albumResp, err := app.AlbumLogic.Upload(ctx, &prm.UploadAlbumParam{
		Name:        album.Name,
		Profile:     album.Description,
		ArtistIds:   artistIds,
		CoverFileId: fileUploadRes.FileIds[0],
		IssueTime:   int32(album.PublishTime / 1000),
		AlbumType:   util.NameToAlbumType(album.Type),
		LiveUrl:     "",
		AlbumMusics: nil,
	})
	if err != nil {
		return 0, err
	}
	return albumResp.AlbumId, nil
}
