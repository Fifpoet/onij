package biz

import (
	"onij/model/api"
	"onij/util"
	"onij/util/boost/exp"
)

type FilePrime struct {
	Id        int64
	Name      *string
	Size      *int32
	StoreKey  *string
	FileType  *api.FileType
	ParentId  *int64
	extra     *string
	CreatedAt *int64
	UpdatedAt *int64
}

func (p *FilePrime) Url() string {
	if p.StoreKey == nil {
		return ""
	}
	return util.DownloadFile(exp.ValueOrZero(p.StoreKey))
}
