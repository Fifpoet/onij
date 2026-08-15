package handler

import (
	"context"
	"fmt"

	"onij/biz/prm"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type listCollectionReq struct {
	Keyword string `json:"keyword"`
}

type createCollectionReq struct {
	Name     string `json:"name"`
	CoverUrl string `json:"cover_url"`
}

type updateCollectionReq struct {
	Id       int64   `json:"id,string"`
	Name     *string `json:"name"`
	CoverUrl *string `json:"cover_url"`
}

type collectionIdReq struct {
	Id int64 `json:"id,string"`
}

type collectionSongMetaReq struct {
	SongId      int64  `json:"song_id"`
	SongName    string `json:"song_name"`
	ArtistNames string `json:"artist_names"`
}

type addCollectionSongsReq struct {
	Id    int64                   `json:"id,string"`
	Songs []collectionSongMetaReq `json:"songs"`
}

type removeCollectionSongReq struct {
	Id     int64 `json:"id,string"`
	SongId int64 `json:"song_id"`
}

func ListCollection(ctx context.Context, c *app.RequestContext) {
	var req listCollectionReq
	_ = c.BindAndValidate(&req)
	resp, err := ser.CollectionLogic.List(ctx, &prm.ListCollectionParam{Keyword: req.Keyword})
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

func CreateCollection(ctx context.Context, c *app.RequestContext) {
	var req createCollectionReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp, err := ser.CollectionLogic.Create(ctx, &prm.CreateCollectionParam{
		Name:     req.Name,
		CoverUrl: req.CoverUrl,
	})
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

func UpdateCollection(ctx context.Context, c *app.RequestContext) {
	var req updateCollectionReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Id <= 0 {
		c.String(consts.StatusBadRequest, "invalid id")
		return
	}
	resp, err := ser.CollectionLogic.Update(ctx, &prm.UpdateCollectionParam{
		Id:       req.Id,
		Name:     req.Name,
		CoverUrl: req.CoverUrl,
	})
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

func DetailCollection(ctx context.Context, c *app.RequestContext) {
	var req collectionIdReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	if req.Id <= 0 {
		c.String(consts.StatusBadRequest, fmt.Sprintf("invalid id: %d", req.Id))
		return
	}
	resp, err := ser.CollectionLogic.Detail(ctx, &prm.CollectionIdParam{Id: req.Id})
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

func DeleteCollection(ctx context.Context, c *app.RequestContext) {
	var req collectionIdReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp, err := ser.CollectionLogic.Delete(ctx, &prm.CollectionIdParam{Id: req.Id})
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

func AddCollectionSongs(ctx context.Context, c *app.RequestContext) {
	var req addCollectionSongsReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	songs := make([]prm.CollectionSongMeta, 0, len(req.Songs))
	for _, s := range req.Songs {
		songs = append(songs, prm.CollectionSongMeta{
			SongId:      s.SongId,
			SongName:    s.SongName,
			ArtistNames: s.ArtistNames,
		})
	}
	resp, err := ser.CollectionLogic.AddSongs(ctx, &prm.AddCollectionSongsParam{
		Id:    req.Id,
		Songs: songs,
	})
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}

func RemoveCollectionSong(ctx context.Context, c *app.RequestContext) {
	var req removeCollectionSongReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp, err := ser.CollectionLogic.RemoveSong(ctx, &prm.RemoveCollectionSongParam{
		Id:     req.Id,
		SongId: req.SongId,
	})
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	c.JSON(consts.StatusOK, resp.Resp())
}
