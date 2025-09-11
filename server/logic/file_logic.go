package logic

import (
	"context"
	"fmt"
	"onij/biz/getter"
	"onij/biz/prm"
	"onij/infra"
	"onij/model"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
)

type FileLogic interface {
	Upload(ctx context.Context, param *prm.UploadFileParam) (*prm.UploadFileResult, error)
	GetList(ctx context.Context, param *prm.GetFileListParam) (*prm.GetFileListResult, error)
	Delete(ctx context.Context, param *prm.DeleteFileParam) (*prm.DeleteFileResult, error)
}

type fileLogic struct {
	*infra.AllInfra
}

func NewFileLogic(i *infra.AllInfra) FileLogic {
	return &fileLogic{
		AllInfra: i,
	}
}

func (l *fileLogic) Delete(ctx context.Context, param *prm.DeleteFileParam) (*prm.DeleteFileResult, error) {
	fis, err := l.FileDal.GetByIds(ctx, param.FileId)
	if err != nil {
		return nil, err
	}
	if len(fis) == 0 {
		return nil, fmt.Errorf("file not found")
	}
	fi := fis[0]

	err = l.FileDal.DeleteById(ctx, param.FileId)
	if err != nil {
		return nil, err
	}

	err = util.DeleteFile(fi.StoreKey)
	if err != nil {
		return nil, err
	}

	return &prm.DeleteFileResult{}, nil
}

func (l *fileLogic) GetList(ctx context.Context, param *prm.GetFileListParam) (*prm.GetFileListResult, error) {
	parentId := int64(0)
	if param.ParentId != nil {
		parentId = *param.ParentId
	}

	files, total, err := l.FileDal.GetListByParentIdAndKeyword(ctx, parentId, param.Keyword, param.Page)
	if err != nil {
		return nil, err
	}

	return &prm.GetFileListResult{
		Files: files,
		Total: total,
	}, nil
}

func (l *fileLogic) Upload(ctx context.Context, param *prm.UploadFileParam) (*prm.UploadFileResult, error) {
	fis := collext.Pick(param.Files, func(f *api.UploadFileReq_FileInfo) *model.File {
		return &model.File{
			Id:       util.IdGen.Generate(),
			Name:     f.Filename,
			Format:   int32(f.Format),
			StoreKey: f.StoreKey,
			ParentId: param.ParentId,
			Hash:     f.Hash,
			Size:     f.Size,
		}
	})

	_, err := l.FileDal.Save(ctx, fis...)
	if err != nil {
		return nil, err
	}

	return &prm.UploadFileResult{
		FileIds: collext.Pick(fis, getter.FileId),
	}, nil
}
