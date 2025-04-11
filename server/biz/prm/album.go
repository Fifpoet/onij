package prm

import "onij/model/api"

type UploadAlbumParam struct {
	Name        string
	ArtistId    int64
	CoverFileId string
	IssueTime   int32
	MusicId     *int64
}

func NewUploadAlbumParam(req *api.UploadAlbumReq) *UploadAlbumParam {
	return &UploadAlbumParam{
		Name:        req.Name,
		ArtistId:    req.ArtistId,
		CoverFileId: req.CoverFileId,
		IssueTime:   req.IssueTime,
		MusicId:     req.MusicId,
	}
}

type UploadAlbumResult struct {
}

func (r *UploadAlbumResult) resp() *api.UploadAlbumResp {
	return &api.UploadAlbumResp{}
}
