package resq

import (
	"onij/enum"
	"onij/infra/mysql"
	"strconv"
	"time"
)

type UpsertRelayReq struct {
	RelayType  int    `json:"relay_type"`
	Content    string `json:"content"`
	Password   int    `json:"password"`
	ExpireTime int    `json:"expire_time"`
	OssKey     string `json:"oss_key"`
}

func (u *UpsertRelayReq) ToModel() *mysql.Relay {
	expire := time.Now()
	switch u.ExpireTime {
	case enum.ExpireTimeNever:
		expire = time.Time{}
	case enum.ExpireTimeFiveMinute:
		expire = time.Now().Add(5 * time.Minute)
	case enum.ExpireTimeOneHour:
		expire = time.Now().Add(1 * time.Hour)
	case enum.ExpireTimeOneDay:
		expire = time.Now().Add(1 * time.Hour)
	case enum.ExpireTimeOneWeek:
		expire = time.Now().Add(7 * 24 * time.Hour)
	case enum.ExpireTimeOneMonth:
		expire = time.Now().Add(30 * 24 * time.Hour)
	}

	return &mysql.Relay{
		RelayType: u.RelayType,
		Password:  u.Password,
		ExpireAt:  strconv.FormatInt(expire.Unix(), 10),
		Content:   u.Content,
		OssKey:    u.OssKey,
		Pin:       false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
