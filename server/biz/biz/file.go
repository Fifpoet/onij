package biz

import (
	"onij/model"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/exp"
)

type FilePrime struct {
	Id        int64
	Name      *string
	Size      *int64
	StoreKey  *string
	FileType  *api.FileType
	Hash      *string
	ParentId  *int64
	extra     *string
	CreatedAt *int64
	UpdatedAt *int64
}

func (p *FilePrime) Model() *model.File {
	if p == nil {
		return nil
	}
	return &model.File{
		Id:       p.Id,
		Name:     exp.ValueOrZero(p.Name),
		Size:     exp.ValueOrZero(p.Size),
		StoreKey: exp.ValueOrZero(p.StoreKey),
		Format:   int32(exp.ValueOrZero(p.FileType)),
		Hash:     exp.ValueOrZero(p.Hash),
		ParentId: exp.ValueOrZero(p.ParentId),
		Extra:    exp.ValueOrZero(p.extra),
	}
}

func (p *FilePrime) Url() string {
	if p == nil || p.StoreKey == nil {
		return ""
	}
	return util.DownloadFile(exp.ValueOrZero(p.StoreKey))
}
func Url(p *FilePrime) string {
	if p == nil || p.StoreKey == nil {
		return ""
	}
	return util.DownloadFile(exp.ValueOrZero(p.StoreKey))
}

func FileToBiz(f *model.File) *FilePrime {
	if f == nil {
		return &FilePrime{}
	}
	return &FilePrime{
		Id:        f.Id,
		Name:      exp.Ptr(f.Name),
		Size:      exp.Ptr(f.Size),
		StoreKey:  exp.Ptr(f.StoreKey),
		FileType:  (*api.FileType)(exp.Ptr(f.Format)),
		Hash:      exp.Ptr(f.Hash),
		ParentId:  exp.Ptr(f.ParentId),
		extra:     exp.Ptr(f.Extra),
		CreatedAt: exp.Ptr(f.CreatedAt.Unix()),
		UpdatedAt: exp.Ptr(f.UpdatedAt.Unix()),
	}
}
