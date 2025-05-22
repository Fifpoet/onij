package prm

import (
	"onij/infra/mysql"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
)

type UploadArtistParam struct {
	Name       string
	ArtistType api.ArtistType
}

func NewUploadArtistParam(req *api.UploadArtistReq) *UploadArtistParam {
	return &UploadArtistParam{
		Name:       req.Name,
		ArtistType: req.ArtistType,
	}
}

type UploadArtistResult struct {
	ArtistId int64
}

func (r *UploadArtistResult) Resp() *api.UploadArtistResp {
	return &api.UploadArtistResp{
		Code:     util.BaseCodeOK,
		Message:  util.BaseMsgOK,
		ArtistId: r.ArtistId,
	}
}

type SearchArtistParam struct {
	Keyword string
	Page    int32
	Limit   int32
}

func NewSearchArtistParam(req *api.SearchArtistReq) *SearchArtistParam {
	return &SearchArtistParam{
		Keyword: req.Keyword,
		Page:    req.Page,
		Limit:   req.Limit,
	}
}

type SearchArtistResult struct {
	Artists []*mysql.Artist
}
func (r *SearchArtistResult) Resp() *api.SearchArtistResp {
	return &api.SearchArtistResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Artists: collext.Pick(r.Artists, func(at *mysql.Artist) *api.ArtistProfile {
			return &api.ArtistProfile{
				Id: at.Id,
				Name: at.Name,
				ArtistType: api.ArtistType(at.ArtistType),
			}
		}),
	}
}