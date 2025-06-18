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
	"onij/util/boost/concurrent"
	"onij/util/boost/crypto"
	"onij/util/boost/exp"
	"time"
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
	// 处理文件夹创建（空文件表示创建文件夹）
	if len(param.Files) == 1 && len(param.Files[0].StoreKey) == 0 {
		// upload folder
		id := util.IdGen.Generate()
		err := l.FileDal.Save(&mysql.File{
			Id:       id,
			Name:     param.Files[0].Filename,
			Format:   int32(api.FileType_FT_Folder),
			Hash:     crypto.Md5([]byte(param.Files[0].Filename)),
			ParentId: param.ParentId,
			OriginAt: exp.Ptr(time.Unix(param.Files[0].OriginAt, 0)),
		})
		if err != nil {
			return nil, err
		}
		return &prm.UploadFileResult{
			FileIds: []int64{id},
		}, nil
	}

	// 检查文件是否已存在（基于hash）
	fis, err := l.FileDal.GetByParentAndHash(param.ParentId, collext.Pick(param.Files, func(f *api.UploadFileReq_FileInfo) string {
		return f.Hash
	})...)
	if err != nil {
		return nil, err
	}
	hashFileMap := collext.Map(fis, getter.FileHash)

	// 过滤出不存在的文件
	toUploads := collext.Select(param.Files, func(f *api.UploadFileReq_FileInfo) (*api.UploadFileReq_FileInfo, bool) {
		return f, hashFileMap[f.Hash] == nil
	})

	// 保存文件信息到数据库
	resIds, err := concurrent.Go(ctx, toUploads, func(ctx context.Context, f *api.UploadFileReq_FileInfo) (int64, error) {
		id := util.IdGen.Generate()
		err := l.FileDal.Save(&mysql.File{
			Id:       id,
			Name:     f.Filename,
			Format:   int32(f.Format),
			StoreKey: f.StoreKey,
			ParentId: param.ParentId,
			Hash:     f.Hash,
			OriginAt: exp.Ptr(time.Unix(f.OriginAt, 0)),
		})
		if err != nil {
			return 0, err
		}
		return id, nil
	})
	if err != nil {
		return nil, err
	}
	return &prm.UploadFileResult{
		FileIds: resIds,
	}, nil
}
