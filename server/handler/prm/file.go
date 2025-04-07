package prm

import "onij/model/api"

type UploadFileParam struct {
	Filename string
	AppId    int32
	Folders  []string
	File     []byte
}

func NewUploadFileParam(req *api.UploadFileReq) *UploadFileParam {
	return &UploadFileParam{
		Filename: req.Filename,
		AppId:    int32(req.AppId),
		Folders:  req.Folders,
		File:     req.File,
	}
}

type UploadFileResult struct {
}
