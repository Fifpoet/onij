package logic

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
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
	err := f.FileDal.Delete(param.FileId)
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
	log.Printf("UploadFileParam: %+v\n", param)

	// 处理文件夹创建（空文件表示创建文件夹）
	if len(param.Files) == 1 && len(param.Files[0].File) == 0 {
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

	// 对每个文件进行base64解码
	decodedFiles := make([]*api.UploadFileReq_FileInfo, len(param.Files))
	for i, fileInfo := range param.Files {
		// 解码base64字符串为原始字节数组
		fileBytes, err := base64.StdEncoding.DecodeString(fileInfo.File)
		if err != nil {
			return nil, fmt.Errorf("base64 decode failed for file %s: %v", fileInfo.Filename, err)
		}

		decodedFile := &api.UploadFileReq_FileInfo{
			Filename: fileInfo.Filename,
			File:     string(fileBytes), // 现在是原始字节数组
			OriginAt: fileInfo.OriginAt,
		}
		decodedFiles[i] = decodedFile
	}

	// check hash
	fis, err := l.FileDal.GetByParentAndHash(param.ParentId, collext.Pick(decodedFiles, func(f *api.UploadFileReq_FileInfo) string {
		return crypto.Md5([]byte(f.File))
	})...)
	if err != nil {
		return nil, err
	}
	hashFileMap := collext.Map(fis, getter.FileHash)

	// 获取文件夹路径, 提取不存在的文件并上传
	folders, err := l.FileDal.GetFolderPathByParentId(param.ParentId)
	if err != nil {
		return nil, err
	}
	toUploads := collext.Select(decodedFiles, func(f *api.UploadFileReq_FileInfo) (*api.UploadFileReq_FileInfo, bool) {
		return f, hashFileMap[crypto.Md5([]byte(f.File))] == nil
	})
	resIds, err := concurrent.Go(ctx, toUploads, func(ctx context.Context, f *api.UploadFileReq_FileInfo) (int64, error) {
		key, err := util.UploadFile(ctx, util.UploadInfo{
			Name:      f.Filename,
			Bytes:     []byte(f.File),
			OssFolder: folders,
		})
		if err != nil {
			return 0, err
		}
		id := util.IdGen.Generate()
		err = l.FileDal.Save(&mysql.File{
			Id:       id,
			Name:     f.Filename,
			Format:   int32(util.GetFileType(f.Filename)),
			StoreKey: key,
			ParentId: param.ParentId,
			Hash:     crypto.Md5([]byte(f.File)),
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
