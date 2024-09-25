package resq

import (
	"database/sql"
	"mime/multipart"
	"onij/enum"
	"onij/infra/mysql"
	"time"
)

type UpsertRelayReq struct {
	RelayType  int                   `form:"relay_type" binding:"required"`
	Content    string                `form:"content" binding:"required"`
	Password   int                   `form:"password"`
	ExpireType int                   `form:"expire_type" binding:"required"`
	File       *multipart.FileHeader `form:"file"` // 用于接收文件
}

func (u *UpsertRelayReq) ToModel() (*mysql.Relay, *multipart.FileHeader) {
	rl := &mysql.Relay{
		RelayType: u.RelayType,
		Password:  &u.Password,
		Content:   u.Content,
		FileOss:   0,
		Pin:       false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	switch u.ExpireType {
	case enum.ExpireTypeNever:
	case enum.ExpireTypeFiveMinute:
		rl.ExpireAt = sql.NullTime{Time: time.Now().Add(5 * time.Minute), Valid: true}
	case enum.ExpireTypeOneHour:
		rl.ExpireAt = sql.NullTime{Time: time.Now().Add(1 * time.Hour), Valid: true}
	case enum.ExpireTypeOneDay:
		rl.ExpireAt = sql.NullTime{Time: time.Now().Add(24 * time.Hour), Valid: true}
	case enum.ExpireTypeOneWeek:
		rl.ExpireAt = sql.NullTime{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true}
	case enum.ExpireTypeOneMonth:
		rl.ExpireAt = sql.NullTime{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true}
	}

	return rl, u.File
}
