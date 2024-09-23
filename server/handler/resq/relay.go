package resq

import (
	"mime/multipart"
	"onij/enum"
	"onij/infra/mysql"
	"strconv"
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
	expire := time.Now()
	switch u.ExpireType {
	case enum.ExpireTypeNever:
		expire = time.Time{}
	case enum.ExpireTypeFiveMinute:
		expire = time.Now().Add(5 * time.Minute)
	case enum.ExpireTypeOneHour:
		expire = time.Now().Add(1 * time.Hour)
	case enum.ExpireTypeOneDay:
		expire = time.Now().Add(24 * time.Hour)
	case enum.ExpireTypeOneWeek:
		expire = time.Now().Add(7 * 24 * time.Hour)
	case enum.ExpireTypeOneMonth:
		expire = time.Now().Add(30 * 24 * time.Hour)
	}

	return &mysql.Relay{
		RelayType: u.RelayType,
		Password:  u.Password,
		ExpireAt:  strconv.FormatInt(expire.Unix(), 10),
		Content:   u.Content,
		OssKey:    "",
		Pin:       false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, u.File
}
