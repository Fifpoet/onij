package resq

import (
	"database/sql"
	"mime/multipart"
	"onij/infra/mysql"
	"time"
)

type UpsertRelayReq struct {
	Id         int                   `form:"id"`
	RelayType  int                   `form:"relay_type" binding:"required"`
	Content    string                `form:"content" binding:"required"`
	Password   int                   `form:"password"`
	ExpireType int                   `form:"expire_type" binding:"required"`
	File       *multipart.FileHeader `form:"file"` // 用于接收文件
}

func (u *UpsertRelayReq) ToModel() (*mysql.Relay, *multipart.FileHeader) {
	var pwd *int
	if u.Password == 0 {
		pwd = nil
	} else {
		pwd = &u.Password
	}
	rl := &mysql.Relay{
		Id:        u.Id,
		RelayType: u.RelayType,
		Password:  pwd,
		Content:   u.Content,
		FileOss:   0,
		Pin:       false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	rl.ExpireAt = sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true}

	return rl, u.File
}
