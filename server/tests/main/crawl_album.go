package main

import (
	"context"
	"onij/biz/errdef"
	"onij/biz/getter"
	"onij/biz/prm"
	"onij/inject"
	"onij/model"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
	"onij/util/boost/exp"
	"onij/util/boost/tool"
	"os"
)

var ctx = context.Background()
var app *inject.App

const ArtistAvatarFolderId = 10000000000002
const AlbumCoverFolderId = 10000000000003

func init() {
	app = inject.InitializeApp()
}

func FetchAlbum(thirdAlbumIds ...int64) (*AlbumResponse, error) {
	file, err := os.ReadFile("/Users/asen/Documents/gopath/src/owner/onij/server/tests/main/jsn/album.json")
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
	_, err = processAlbum(album.Album, album.Songs, artistIds)
	if err != nil {
		return nil, err
	}

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
		ParentId: AlbumCoverFolderId,
		Url:      exp.Ptr(album.PicURL),
	})
	if err != nil {
		return 0, err
	}

	// album.songs
	albumResp, err := app.AlbumLogic.Upload(ctx, &prm.UploadAlbumParam{
		Name:        album.Name,
		Profile:     album.Description,
		ArtistIds:   artistIds,
		CoverFileId: fileUploadRes.FileIds[0],
		IssueTime:   int32(album.PublishTime / 1000),
		AlbumType:   util.NameToAlbumType(album.Type),
		LiveUrl:     "",
		ThirdId:     exp.Ptr(int64(album.ID)),
		AlbumMusics: collext.Pick(songs, func(song Song) *api.UploadAlbumReq_AlbumMusic {
			return &api.UploadAlbumReq_AlbumMusic{
				Name:        song.Name,
				ArtistName:  collext.Pick(song.Ar, func(a Artist) string { return a.Name }),
				TimeLength:  int32(song.Dt),
				MusicId:     0,
				IsAvailable: false,
				ThirdId:     exp.Ptr(int64(song.ID)),
			}
		}),
	})
	if err != nil {
		return 0, err
	}
	return albumResp.AlbumId, nil
}

func FetchMusic(thirdMusicIds ...int64) error {
	file, err := os.ReadFile("/Users/asen/Documents/gopath/src/owner/onij/server/tests/main/jsn/lyric.json")
	if err != nil {
		return err
	}

	musics, err := app.MusicDal.GetByIds(ctx, thirdMusicIds...)
	if err != nil {
		return err
	}
	if len(musics) < len(thirdMusicIds) {
		return errdef.ErrThirdMusicIdNotFound
	}
	musicMap := collext.Map(musics, getter.MusicThirdId)
	for _, thirdId := range thirdMusicIds {
		lyc := tool.LoadJson[LyricResponse](string(file), true)
		if len(lyc.Lrc.Lyric) == 0 {
			return errdef.ErrResponseNoLyric
		}
		_, err = app.MusicDal.UpdateById(ctx, func(u model.MusicUpdater) {
			u.LyricContent(lyc.Lrc.Lyric)
		}, musicMap[thirdId].Id)
		if err != nil {
			return err
		}
	}
	return nil
}
