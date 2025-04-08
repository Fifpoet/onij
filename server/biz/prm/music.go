package prm

import "onij/model/api"

type UploadMusicParam struct {
	Name         string
	Artists      []string
	Mp3FileId    int64
	LyricsFileId int64
	Composer     *string
	Writer       *string
	AlbumId      *int64
	MvUrl        *string
	RootMusicId  *int64
	IssueTime    *int32
}

func NewUploadMusicParam(req *api.UploadMusicReq) *UploadMusicParam {
	return &UploadMusicParam{
		Name:         req.Name,
		Artists:      req.Artists,
		Mp3FileId:    req.Mp3FileId,
		LyricsFileId: req.LyricsFileId,
		Composer:     req.Composer,
		Writer:       req.Writer,
		AlbumId:      req.AlbumId,
		MvUrl:        req.MvUrl,
		RootMusicId:  req.RootMusicId,
		IssueTime:    req.IssueTime,
	}
}

type UploadMusicResult struct {
}

func (p *UploadMusicParam) Response() *UploadMusicResult {
	return &UploadMusicResult{}
}
