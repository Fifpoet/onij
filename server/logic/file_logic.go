package logic

import (
	"context"
	"onij/handler/prm"
	"onij/infra"
	"onij/infra/mysql"
	"onij/util"
	"onij/util/boost/crypto"
)

type FileLogic interface {
	Upload(ctx context.Context, prm *prm.UploadFileParam) (*prm.UploadFileResult, error)
}

type fileLogic struct {
	*infra.AllInfra
}

func NewFileLogic(i *infra.AllInfra) FileLogic {
	return &fileLogic{
		AllInfra: i,
	}
}

func (l *fileLogic) Upload(ctx context.Context, param *prm.UploadFileParam) (*prm.UploadFileResult, error) {
	key, err := util.UploadFile(ctx, util.UploadInfo{
		Name:      param.Filename,
		Bytes:     param.File,
		OssFolder: param.Folders,
	})
	if err != nil {
		return nil, err
	}
	err = l.FileDal.Save(&mysql.File{
		Name:     param.Filename,
		Format:   0,
		StoreKey: key,
		Hash:     crypto.Md5(param.File),
	})
	if err != nil {
		return nil, err
	}
	return &prm.UploadFileResult{}, nil
}
