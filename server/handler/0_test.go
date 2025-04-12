package handler

import (
	"context"
	"onij/biz/prm"
	"onij/inject"
	"onij/util/boost/exp"
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
		SortType: 1,
		Keyword:  exp.Ptr("张"),
		AlbumId:  nil,
		ArtistId: nil,
		Page:     1,
		Limit:    10,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(tool.ToJson(list))
}
