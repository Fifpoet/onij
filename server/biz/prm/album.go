package prm

import (
	"onij/model"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
)

type UploadAlbumParam struct {
	Name        string
	Profile     string
	ArtistIds   []int64
	CoverFileId int64
	IssueTime   int32
	AlbumType   api.AlbumType
	LiveUrl     string
	AlbumMusics []*api.UploadAlbumReq_AlbumMusic
	AlbumId     *int64
}

func NewUploadAlbumParam(req *api.UploadAlbumReq) *UploadAlbumParam {
	return &UploadAlbumParam{
		Name:        req.Name,
		Profile:     req.Profile,
		ArtistIds:   req.ArtistIds,
		CoverFileId: req.CoverFileId,
		IssueTime:   req.IssueTime,
		AlbumType:   req.AlbumType,
		LiveUrl:     req.LiveUrl,
		AlbumMusics: req.AlbumMusics,
		AlbumId:     req.AlbumId,
	}
}

type UploadAlbumResult struct {
	AlbumId int64
}

func (r *UploadAlbumResult) Resp() *api.UploadAlbumResp {
	return &api.UploadAlbumResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		AlbumId: r.AlbumId,
	}
}

type GetAlbumDetailParam struct {
	AlbumId int64
}

func NewGetAlbumParam(req *api.GetAlbumDetailReq) *GetAlbumDetailParam {
	return &GetAlbumDetailParam{
		AlbumId: req.AlbumId,
	}
}

type GetAlbumDetailResult struct {
	Album        *model.Album
	CoverUrl     string
	Artists      []*model.Artist
	MusicTagsMap map[int64][]*model.Tag
	AlbumMusics  []*model.AlbumMusic
}

func (r *GetAlbumDetailResult) Resp() *api.GetAlbumDetailResp {
	return &api.GetAlbumDetailResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Album: &api.Album{
			Id:        r.Album.Id,
			Name:      r.Album.Name,
			Profile:   r.Album.Profile,
			AlbumType: api.AlbumType(r.Album.AlbumType),
			//Artists:   ,
			IssueTime:    int32(r.Album.IssueTime.Unix()),
			CoverFileUrl: r.CoverUrl,
			LiveUrl:      r.Album.LiveUrl,
		},
		AlbumMusics: collext.Pick(r.AlbumMusics, func(am *model.AlbumMusic) *api.AlbumMusic {
			return &api.AlbumMusic{
				Name:        am.Name,
				TimeLength:  am.TimeLength,
				MusicId:     am.Id,
				AlbumId:     r.Album.Id,
				IsAvailable: util.IntToBool(am.IsAvailable),
			}
		}),
	}
}
