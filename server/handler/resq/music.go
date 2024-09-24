package resq

import (
	"mime/multipart"
	"onij/infra/mysql"
)

type UpsertMusicReq struct {
	MusicId int                   `json:"music_id"`
	Cover   string                `json:"cover"`
	File    *multipart.FileHeader `json:"file"`
}

func (u *UpsertMusicReq) ToModel() (*mysql.Music, *multipart.FileHeader) {

	return nil, nil
}
