package logic

import (
	"context"
	"fmt"
	"onij/biz/getter"
	"onij/biz/prm"
	"onij/infra"
	"onij/infra/mysql"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
)

type FileLogic interface {
	Upload(ctx context.Context, prm *prm.UploadFileParam) (*prm.UploadFileResult, error)
	GetList(ctx context.Context, param *prm.GetFileListParam) (*prm.GetFileListResult, error)
	DownloadByIds(ctx context.Context, param *prm.DownloadFileParam) (*prm.DownloadFileResult, error)
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

func (f *fileLogic) Delete(ctx context.Context, param *prm.DeleteFileParam) (*prm.DeleteFileResult, error) {
	fis, err := f.FileDal.GetByIds(param.FileId)
	if err != nil {
		return nil, err
	}
	if len(fis) == 0 {
		return nil, fmt.Errorf("file not found")
	}
	fi := fis[0]
	err = f.FileDal.Delete(param.FileId)
	if err != nil {
		return nil, err
	}
	err = util.DeleteFile(fi.StoreKey)
	if err != nil {
		return nil, err
	}
	return &prm.DeleteFileResult{}, nil
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
	files, total, err := f.FileDal.GetListByParentIdAndKeyword(param.ParentId, param.Keyword, param.Page)
	if err != nil {
		return nil, err
	}

	return &prm.GetFileListResult{
		Files: files,
		Total: total,
	}, nil
}

func (l *fileLogic) Upload(ctx context.Context, param *prm.UploadFileParam) (*prm.UploadFileResult, error) {
	// 处理文件夹创建（空文件表示创建文件夹）
	if len(param.Files) == 1 && len(param.Files[0].StoreKey) == 0 {
		// upload folder
		id := util.IdGen.Generate()
		err := l.FileDal.Save(&mysql.File{
			Id:       id,
			Name:     param.Files[0].Filename,
			Format:   int32(api.FileType_FT_Folder),
			Hash:     param.Files[0].Hash,
			ParentId: param.ParentId,
			OriginAt: param.Files[0].OriginAt,
		})
		if err != nil {
			return nil, err
		}
		return &prm.UploadFileResult{
			FileIds: []int64{id},
		}, nil
	}

	fis := collext.Pick(param.Files, func(f *api.UploadFileReq_FileInfo) *mysql.File {
		return &mysql.File{
			Id:       util.IdGen.Generate(),
			Name:     f.Filename,
			Format:   int32(f.Format),
			StoreKey: f.StoreKey,
			ParentId: param.ParentId,
			Hash:     f.Hash,
			OriginAt: f.OriginAt,
		}
	})
	err := l.FileDal.Save(fis...)
	if err != nil {
		return nil, err
	}

	return &prm.UploadFileResult{
		FileIds: collext.Pick(fis, getter.FileId),
	}, nil
}
