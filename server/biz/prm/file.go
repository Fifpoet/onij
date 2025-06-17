package prm

import (
	"onij/infra/mysql"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
	"slices"
)

type UploadFileParam struct {
	Filename string
	Folders  []string
	File     []byte
	OriginAt int64
}

func NewUploadFileParam(req *api.UploadFileReq) *UploadFileParam {
	return &UploadFileParam{
		Filename: req.Filename,
		Folders:  req.Folders,
		File:     req.File,
		OriginAt: req.OriginAt,
	}
}

type UploadFileResult struct {
	FileId  int64
	FileUrl string
}

func (r *UploadFileResult) Resp() *api.UploadFileResp {
	return &api.UploadFileResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		FileId:  r.FileId,
		FileUrl: r.FileUrl,
	}
}

type GetFileListParam struct {
	ParentId int64
	Page     util.Page
}

func NewGetFileListParam(req *api.GetFileListReq) *GetFileListParam {
	return &GetFileListParam{
		ParentId: req.ParentId,
		Page: util.Page{
			Page:  req.Page,
			Limit: req.Limit,
		},
	}
}

type GetFileListResult struct {
	Files []*mysql.File
	Total int32
}

func (r *GetFileListResult) Resp() *api.GetFileListResp {
	res := collext.Pick(r.Files, func(f *mysql.File) *api.FileDetail {
		var originAt int64 = 0
		if f.OriginAt != nil {
			originAt = f.OriginAt.Unix()
		}

		return &api.FileDetail{
			Id:        f.Id,
			Name:      f.Name,
			Format:    api.FileType(f.Format),
			ParentId:  f.ParentId,
			OriginAt:  originAt,
			CreatedAt: f.CreatedAt.Unix(),
			UpdatedAt: f.UpdatedAt.Unix(),
		}
	})
	slices.SortFunc(res, func(a, b *api.FileDetail) int {
		return int(b.OriginAt - a.OriginAt)
	})
	return &api.GetFileListResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Files:   res,
		Total:   r.Total,
	}
}

type DownloadFileParam struct {
	FileIds []int64
}

func NewDownloadFileParam(req *api.DownloadFileReq) *DownloadFileParam {
	return &DownloadFileParam{
		FileIds: req.FileIds,
	}
}

type DownloadFileResult struct {
	Urls []string
}

func (r *DownloadFileResult) Resp() *api.DownloadFileResp {
	return &api.DownloadFileResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Urls:    r.Urls,
	}
}
