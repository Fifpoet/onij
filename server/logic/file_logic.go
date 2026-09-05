package logic

import (
	"context"
	"fmt"
	"onij/biz/biz"
	"onij/biz/errdef"
	"onij/biz/getter"
	"onij/biz/prm"
	"onij/infra"
	"onij/model"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
	"onij/util/logs"
	"time"
)

type FileLogic interface {
	Upload(ctx context.Context, param *prm.UploadFileParam) (*prm.UploadFileResult, error)
	GetList(ctx context.Context, param *prm.GetFileListParam) (*prm.GetFileListResult, error)
	PreviewRedirect(ctx context.Context, fileId int64, attname string) (string, error)
	PreviewContent(ctx context.Context, fileId int64) (name string, data []byte, err error)
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
		return &prm.DeleteFileResult{}, nil
	}
	fi := fis[0]

	err = l.FileDal.DeleteById(ctx, param.FileId)
	if err != nil {
		return nil, err
	}

	if delErr := l.FileVideoMarkerDal.DeleteByFileId(ctx, param.FileId); delErr != nil {
		logs.Error("fileLogic Delete markers file_id=%d: %v", param.FileId, delErr)
	}

	if len(fi.StoreKey) > 0 {
		_ = util.DeleteFile(fi.StoreKey)
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

func (l *fileLogic) PreviewRedirect(ctx context.Context, fileId int64, attname string) (string, error) {
	fis, err := l.FileDal.GetByIds(ctx, fileId)
	if err != nil {
		return "", err
	}
	if len(fis) == 0 || len(fis[0].StoreKey) == 0 {
		return "", fmt.Errorf("file not found")
	}
	return util.BrowserFileURL(fis[0].StoreKey, attname), nil
}

func (l *fileLogic) PreviewContent(ctx context.Context, fileId int64) (string, []byte, error) {
	fis, err := l.FileDal.GetByIds(ctx, fileId)
	if err != nil {
		return "", nil, err
	}
	if len(fis) == 0 || len(fis[0].StoreKey) == 0 {
		return "", nil, fmt.Errorf("file not found")
	}
	fi := fis[0]
	data, err := util.GetFile(fi.StoreKey)
	if err != nil {
		return "", nil, err
	}
	return fi.Name, data, nil
}

func (l *fileLogic) Upload(ctx context.Context, param *prm.UploadFileParam) (*prm.UploadFileResult, error) {
	if param.Url == nil {
		fis := make([]*model.File, 0, len(param.Files))
		for _, f := range param.Files {
			storeKey := f.StoreKey
			if f.Format == api.FileType_FT_Folder && len(storeKey) == 0 {
				folder, err := l.FileDal.GetFolderPathByParentId(ctx, param.ParentId)
				if err != nil {
					return nil, err
				}
				key, err := util.UploadFile(ctx, util.UploadInfo{
					Name:      f.Filename + "/.keep",
					Bytes:     []byte(" "),
					OssFolder: folder,
				})
				if err != nil {
					return nil, err
				}
				storeKey = key
			}

			originAt := biz.ParseOriginAtFromExifField(f.Exif)
			fi := &model.File{
				Id:       util.IdGen.Generate(),
				Name:     f.Filename,
				Format:   int32(f.Format),
				StoreKey: storeKey,
				ParentId: param.ParentId,
				Hash:     f.Hash,
				Size:     f.Size,
				Extra:    biz.BuildFileExtra(originAt, f.Exif),
			}
			if originAt > 0 {
				fi.CreatedAt = time.Unix(originAt, 0)
			}
			fis = append(fis, fi)
		}
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
