package prm

import "onij/model/api"

type UploadMusicParam struct {
}

func NewUploadMusicParam(req *api.UploadMusicReq) *UploadMusicParam {
	return &UploadMusicParam{}
}

type UploadMusicResult struct {
}

func (p *UploadMusicParam) Response() *UploadMusicResult {
	return &UploadMusicResult{}
}
