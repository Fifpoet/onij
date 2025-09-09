package prm

import (
	"onij/biz/getter"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
)

type UploadAlbumParam struct {
	Name           string
	ArtistId       int64
	CoverFileId    int64
	IssueTime      int32
	MusicId        *int64
	RelatedAlbumId *int64
}

func NewUploadAlbumParam(req *api.UploadAlbumReq) *UploadAlbumParam {
	return &UploadAlbumParam{
		Name:           req.Name,
		ArtistId:       req.ArtistId,
		CoverFileId:    req.CoverFileId,
		IssueTime:      req.IssueTime,
		MusicId:        req.MusicId,
		RelatedAlbumId: req.RelatedAlbumId,
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
	Album           *model.Album
	Musics          []*model.Music
	CoverUrl        string
	Artist          *model.Artist
	MusicArtistsMap map[int64][]*model.Artist
	MusicTagsMap    map[int64][]*model.Tag
}

func (r *GetAlbumDetailResult) Resp() *api.GetAlbumDetailResp {
	return &api.GetAlbumDetailResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Detail: &api.AlbumDetail{
			Id:           r.Album.Id,
			Name:         r.Album.Name,
			ArtistId:     r.Artist.Id,
			ArtistName:   r.Artist.Name,
			IssueTime:    r.Album.IssueTime,
			CoverFileUrl: r.CoverUrl,
			MvUrl:        r.Album.MvUrl,
			RelatedMusics: collext.Pick(r.Musics, func(music *model.Music) *api.MusicProfile {
				return &api.MusicProfile{
					Id:             music.Id,
					Name:           music.Name,
					SingerProfiles: collext.Pick(r.MusicArtistsMap[music.Id], getter.ArtistToProfile),
					Tags:           collext.Pick(r.MusicTagsMap[music.Id], getter.TagToDetail),
				}
			}),
		},
	}
}
