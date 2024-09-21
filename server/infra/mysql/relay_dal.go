package mysql

import (
	"gorm.io/gorm"
	"time"
)

type RelayDal interface {
	GetByIds(id []int) ([]*Relay, error)
	GetRelayByType(relayType int) ([]*Relay, error)
	GetRelayByPwd(password int) (*Relay, error)

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

type Relay struct {
	Id        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	RelayType int    `json:"relay_type"`
	Password  int    `json:"password"`
	ExpireAt  string `json:"expire_at"`
	Content   string `json:"content"`
	OssKey    string `json:"oss_key"`
	Pin       bool   `json:"pin"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (r *relayDal) GetRelayByType(relayType int) ([]*Relay, error) {
	now := time.Now().Unix()
	var rs []*Relay
	err := r.db.Where("relay_type = ?", relayType).Where("expire_at > ?", now).
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
	err := r.db.Where("id in ?", id).Where("expire_at > ?", now).
		Order("pin desc").Order("created_at desc").
		Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *relayDal) DelByType(relayType int) error {
	return r.db.Where("relay_type = ?", relayType).Delete(&Relay{}).Error
}

func (r *relayDal) GetRelayByPwd(password int) (*Relay, error) {
	res := &Relay{}
	err := r.db.Where("password = ?", password).First(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
