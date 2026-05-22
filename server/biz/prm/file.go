package prm

import (
	"onij/model"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
	"slices"
)

type UploadFileParam struct {
	ParentId int64
	Files    []*api.UploadFileReq_FileInfo
	Url      *string
}

func NewUploadFileParam(req *api.UploadFileReq) *UploadFileParam {
	return &UploadFileParam{
		ParentId: req.ParentId,
		Files:    req.Files,
		Url:      req.Url,
	}
}

type UploadFileResult struct {
	FileIds []int64
}

func (r *UploadFileResult) Resp() *api.UploadFileResp {
	return &api.UploadFileResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		FileIds: r.FileIds,
	}
}

type GetFileListParam struct {
	ParentId *int64
	Keyword  *string
	Page     util.Page
}

func NewGetFileListParam(req *api.GetFileListReq) *GetFileListParam {
	return &GetFileListParam{
		ParentId: req.ParentId,
		Keyword:  req.Keyword,
		Page: util.Page{
			Page:  req.Page,
			Limit: req.Limit,
		},
	}
}

type GetFileListResult struct {
	Files []*model.File
	Total int32
}

func (r *GetFileListResult) Resp() *api.GetFileListResp {
	res := collext.Pick(r.Files, func(f *model.File) *api.File {
		return &api.File{
			Id:        f.Id,
			Name:      f.Name,
			Format:    api.FileType(f.Format),
			ParentId:  f.ParentId,
			CreatedAt: f.CreatedAt.Unix(),
			UpdatedAt: f.UpdatedAt.Unix(),
			Url:       util.DownloadFile(f.StoreKey),
		}
	})
	slices.SortFunc(res, func(a, b *api.File) int {
		if (a.Format == api.FileType_FT_Folder || b.Format == api.FileType_FT_Folder) &&
			a.Format != b.Format {
			return int(b.Format - a.Format) // format大的在前面
		}
		return int(a.Id - b.Id)
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

type DownloadFileResult struct {
	Urls []string
}

func (r *DownloadFileResult) Resp() map[string]interface{} {
	urls := r.Urls
	if urls == nil {
		urls = []string{}
	}
	return map[string]interface{}{
		"code":    util.BaseCodeOK,
		"message": util.BaseMsgOK,
		"urls":    urls,
	}
}

type DeleteFileParam struct {
	FileId int64
}

func NewDeleteFileParam(req *api.DeleteFileReq) *DeleteFileParam {
	return &DeleteFileParam{
		FileId: req.FileId,
	}
}

type DeleteFileResult struct {
}

func (r *DeleteFileResult) Resp() *api.DeleteFileResp {
	return &api.DeleteFileResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
	}
}
