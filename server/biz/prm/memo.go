package prm

import (
	"onij/infra/mysql"
	api "onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
)

type GetMemoDetailParam struct {
	Id int64
}

func NewGetMemoDetailParam(req *api.GetMemoDetailReq) *GetMemoDetailParam {
	return &GetMemoDetailParam{
		Id: req.MemoId,
	}
}

type GetMemoDetailResult struct {
	Memo *mysql.Memo
	CoverFile *mysql.File
}

func (r *GetMemoDetailResult) Resp() *api.GetMemoDetailResp {
	return &api.GetMemoDetailResp{
		Code: util.BaseCodeOK,
		Message: util.BaseMsgOK,
		MemoDetail: &api.MemoDetail{
			Id: r.Memo.Id,
			Title: r.Memo.Title,
			Content: r.Memo.Content,
			OriginAt: r.Memo.OriginAt,
			CoverUrl: util.DownloadFile(r.CoverFile.StoreKey),
			ProfileContent: r.Memo.ProfileContent,
		},
	}
}

type GetMemoListParam struct {
	Keyword string
	util.Page
}

func NewGetMemoListParam(req *api.GetMemoListReq) *GetMemoListParam {
	return &GetMemoListParam{
		Keyword: req.Keyword,
		Page: util.Page{
			Page: req.Page,
			Limit: req.Limit,
		},
	}
}

type GetMemoListResult struct {
	Memos []*mysql.Memo
	CoverFileMap map[int64]*mysql.File
}

func (r *GetMemoListResult) Resp() *api.GetMemoListResp {
	return &api.GetMemoListResp{
		Code: util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Memos: collext.Pick(r.Memos, func(m *mysql.Memo) *api.MemoProfile {
			return &api.MemoProfile{
				Id: m.Id,
				Title: m.Title,
				CoverUrl: util.DownloadFile(r.CoverFileMap[m.CoverFileId].StoreKey),
				ProfileContent: m.ProfileContent,
				CreatedAt: m.CreatedAt.Unix(),
				UpdatedAt: m.UpdatedAt.Unix(),
				OriginAt: m.OriginAt,
			}
		}),
	}
}

type UploadMemoParam struct {
	Title string
	Content string
	OriginAt int64
	CoverFileId int64
	ProfileContent string
}

func NewUploadMemoParam(req *api.UploadMemoReq) *UploadMemoParam {	
	return &UploadMemoParam{
		Title: req.Title,
		Content: req.Content,
		OriginAt: req.OriginAt,
		CoverFileId: req.CoverFileId,
		ProfileContent: req.ProfileContent,
	}
}

type UploadMemoResult struct {
	Id int64
}

func (r *UploadMemoResult) Resp() *api.UploadMemoResp {
	return &api.UploadMemoResp{
		Code: util.BaseCodeOK,
		Message: util.BaseMsgOK,
		MemoId: r.Id,
	}
}