// 单独运行：cd server && go run ./cmd/ensure_cloud_folder
// 前置：数据库可连、环境变量 QINIU_SK（或 sk）已设置
package main

import (
	"context"
	"fmt"
	"onij/biz/prm"
	"onij/inject"
	"onij/model/api"
	"onij/util"
)

const cloudFolderName = "cloud"

func main() {
	app := inject.InitializeApp()
	ctx := context.Background()

	parentId := int64(0)
	files, _, err := app.FileDal.GetListByParentIdAndKeyword(ctx, &parentId, nil, util.Page{Page: 1, Limit: 200})
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if f.Format == int32(api.FileType_FT_Folder) && f.Name == cloudFolderName {
			fmt.Printf("cloud folder already exists: id=%d\n", f.Id)
			return
		}
	}

	res, err := app.FileLogic.Upload(ctx, &prm.UploadFileParam{
		ParentId: 0,
		Files: []*api.UploadFileReq_FileInfo{{
			Filename: cloudFolderName,
			StoreKey: "",
			Hash:     "",
			Format:   api.FileType_FT_Folder,
			Size:     0,
		}},
	})
	if err != nil {
		panic(err)
	}
	if len(res.FileIds) == 0 {
		panic("create cloud folder failed")
	}
	fmt.Printf("cloud folder created: id=%d\n", res.FileIds[0])
}
