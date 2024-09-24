package mysql

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type RelayDal interface {
	GetByIds(id []int) ([]*Relay, error)
	GetRelayByType(relayType int) ([]*Relay, error)
	GetRelayByPwd(password int) ([]*Relay, error)
	GetRelayByPwdAndId(password, id int) (*Relay, error)

	Save(relay *Relay) (int, error)
	DelById(id int) error
	DelByType(relayType int) error
}

type relayDal struct {
	db *gorm.DB
}

func NewRelayDal(db *gorm.DB) RelayDal {
	return &relayDal{db: db}
}

// Relay .
// RelayType 为文本, Oss则为其附件. 为文件则Content为标题
// Password 默认为0, 非0则返回时隐藏内容 只显示内容提示, 由id+pwd获取内容
// ExpireAt 为0则永不过期, 否则为过期时间戳
// Pin 排序时Pin 为true的优先, 删除时跳过
type Relay struct {
	Id        int          `json:"id" gorm:"primaryKey;autoIncrement"`
	RelayType int          `json:"relay_type"`
	Password  int          `json:"password"`
	ExpireAt  sql.NullTime `json:"expire_at"`
	Content   string       `json:"content"`
	OssKey    string       `json:"oss_key"`
	Pin       bool         `json:"pin"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (r *relayDal) GetRelayByType(relayType int) ([]*Relay, error) {
	now := time.Now().Unix()
	var rs []*Relay
	err := r.db.Where("relay_type = ?", relayType).Where("expire_at > ? or expire_at = 0", now).
		Order("pin desc").Order("created_at desc").
		Find(&rs).Error
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (r *relayDal) Save(relay *Relay) (int, error) {
	err := r.db.Save(relay).Error
	if err != nil {
		return 0, err
	}
	return relay.Id, nil
}

func (r *relayDal) DelById(id int) error {
	err := r.db.Where("id = ? AND pin = ?", id, false).Delete(&Relay{}).Error
	return err
}

func (r *relayDal) GetByIds(id []int) ([]*Relay, error) {
	now := time.Now().Unix()
	var res []*Relay
	err := r.db.Where("id in ?", id).Where("expire_at > ? or expire_at = 0", now).
		Order("pin desc").Order("created_at desc").
		Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *relayDal) DelByType(relayType int) error {
	return r.db.Where("relay_type = ? AND pin = ?", relayType, false).Delete(&Relay{}).Error
}

func (r *relayDal) GetRelayByPwd(password int) ([]*Relay, error) {
	var res []*Relay
	now := time.Now().Unix()
	err := r.db.Where("password = ?", password).Where("expire_at > ? or expire_at = 0", now).
		Order("pin desc").Order("created_at desc").
		Find(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *relayDal) GetRelayByPwdAndId(password, id int) (*Relay, error) {
	res := &Relay{}
	err := r.db.Where("password = ? AND id = ?", password, id).First(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
