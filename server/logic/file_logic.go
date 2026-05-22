package logic

import (
	"context"
	"fmt"
	"onij/biz/errdef"
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
	Download(ctx context.Context, param *prm.DownloadFileParam) (*prm.DownloadFileResult, error)
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

func (l *fileLogic) Download(ctx context.Context, param *prm.DownloadFileParam) (*prm.DownloadFileResult, error) {
	files, err := l.FileDal.GetByIds(ctx, param.FileIds...)
	if err != nil {
		return nil, err
	}
	urls := make([]string, 0, len(files))
	for _, f := range files {
		if f == nil || len(f.StoreKey) == 0 {
			continue
		}
		urls = append(urls, util.DownloadFile(f.StoreKey))
	}
	return &prm.DownloadFileResult{Urls: urls}, nil
}

func (l *fileLogic) GetList(ctx context.Context, param *prm.GetFileListParam) (*prm.GetFileListResult, error) {
	files, total, err := l.FileDal.GetListByParentIdAndKeyword(ctx, param.ParentId, param.Keyword, param.Page)
	if err != nil {
		return nil, err
	}

	return &prm.GetFileListResult{
		Files: files,
		Total: total,
	}, nil
}

func (l *fileLogic) Upload(ctx context.Context, param *prm.UploadFileParam) (*prm.UploadFileResult, error) {
	if param.Url == nil {
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

	url := *param.Url
	name := util.GetUrlResourceName(url)
	ext := util.GetFileSuffix(name)
	if len(ext) == 0 {
		return nil, errdef.ErrUrlFileExtUnknown
	}
	bytes, err := util.GET(url, nil)
	if err != nil {
		return nil, err
	}
	hash, err := util.CalcFileHash(bytes)
	if err != nil {
		return nil, err
	}

	// check hash
	fil, err := l.FileDal.GetByParentAndHash(ctx, param.ParentId, hash)
	if err != nil {
		return nil, err
	}
	if fil != nil {
		return &prm.UploadFileResult{
			FileIds: []int64{fil.Id},
		}, nil
	}

	folder, err := l.FileDal.GetFolderPathByParentId(ctx, param.ParentId)
	if err != nil {
		return nil, err
	}
	key, err := util.UploadFile(ctx, util.UploadInfo{
		Name:      name,
		Bytes:     bytes,
		OssFolder: folder,
	})
	if err != nil {
		return nil, err
	}

	id := util.IdGen.Generate()
	_, err = l.FileDal.Save(ctx, &model.File{
		Id:       id,
		Name:     name,
		Format:   int32(util.GetFileType(name)),
		Size:     int64(len(bytes)),
		StoreKey: key,
		Hash:     hash,
		ParentId: param.ParentId,
		Extra:    "",
	})
	if err != nil {
		return nil, err
	}
	return &prm.UploadFileResult{
		FileIds: []int64{id},
	}, nil
}
