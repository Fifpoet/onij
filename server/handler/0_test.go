package handler

import (
	"context"
	"onij/biz/prm"
	"onij/inject"
	"onij/util/boost/tool"
	"testing"
)

var ctx = context.Background()

func init() {
	ser = inject.InitializeApp()
}

func TestGetMusicDetail(t *testing.T) {
	detail, err := ser.MusicLogic.GetDetail(ctx, &prm.GetMusicDetailParam{
		Id: 10318921068806,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(tool.ToJson(detail))
}

func TestGetMusicList(t *testing.T) {
	list, err := ser.MusicLogic.GetList(ctx, &prm.GetMusicListParam{
		SortType:    1,
		Keywords:    []string{},
		AlbumId:     nil,
		ArtistIds:   nil,
		WriterIds:   []int64{10321609615619},
		// ComposerIds: []int64{10321609615633},
		TagTypes:    []int32{},
		Page:        1,
		Limit:       20,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(tool.ToJson(list.Musics))
	t.Log(tool.ToJson(list.MusicArtistsMap))
}
