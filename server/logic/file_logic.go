package logic

import (
	"context"
	"onij/biz/prm"
	"onij/infra"
	"onij/infra/mysql"
	"onij/util"
	"onij/util/boost/collection/collext"
	"onij/util/boost/crypto"
)

type FileLogic interface {
	Upload(ctx context.Context, prm *prm.UploadFileParam) (*prm.UploadFileResult, error)
	GetList(ctx context.Context, param *prm.GetFileListParam) (*prm.GetFileListResult, error)
	DownloadByIds(ctx context.Context, param *prm.DownloadFileParam) (*prm.DownloadFileResult, error)
}

type fileLogic struct {
	*infra.AllInfra
}

func NewFileLogic(i *infra.AllInfra) FileLogic {
	return &fileLogic{
		AllInfra: i,
	}
}

func (f *fileLogic) DownloadByIds(ctx context.Context, param *prm.DownloadFileParam) (*prm.DownloadFileResult, error) {
	files, err := f.FileDal.GetByIds(param.FileIds...)
	if err != nil {
		return nil, err
	}
	return &prm.DownloadFileResult{
		Urls: collext.Pick(files, func(f *mysql.File) string {
			return util.DownloadFile(f.StoreKey)
		}),
	}, nil
}

func (f *fileLogic) GetList(ctx context.Context, param *prm.GetFileListParam) (*prm.GetFileListResult, error) {
	files, err := f.FileDal.GetByParentId(param.ParentId, param.Page)
	if err != nil {
		return nil, err
	}

	// 获取总数
	total, err := f.FileDal.CountByParentId(param.ParentId)
	if err != nil {
		return nil, err
	}

	return &prm.GetFileListResult{
		Files: files,
		Total: total,
	}, nil
}

func (l *fileLogic) Upload(ctx context.Context, param *prm.UploadFileParam) (*prm.UploadFileResult, error) {
	// check hash
	fi, err := l.FileDal.GetByHash(crypto.Md5(param.File))
	if err != nil {
		return nil, err
	}
	if fi != nil {
		return &prm.UploadFileResult{
			FileId:  fi.Id,
			FileUrl: util.DownloadFile(fi.StoreKey),
		}, nil
	}

	key, err := util.UploadFile(ctx, util.UploadInfo{
		Name:      param.Filename,
		Bytes:     param.File,
		OssFolder: param.Folders,
	})
	if err != nil {
		return nil, err
	}
	fi = &mysql.File{
		Id:       util.IdGen.Generate(),
		Name:     param.Filename,
		Format:   int32(util.GetFileType(param.Filename)),
		StoreKey: key,
		Hash:     crypto.Md5(param.File),
	}
	err = l.FileDal.Save(fi)
	if err != nil {
		return nil, err
	}
	return &prm.UploadFileResult{
		FileId:  fi.Id,
		FileUrl: util.DownloadFile(fi.StoreKey),
	}, nil
}
