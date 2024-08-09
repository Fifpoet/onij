package mysql

import (
	"gorm.io/gorm"
	"time"
)

type RelayDal interface {
	GetRelayByType(relayType int) ([]*Relay, error)

	Save(relay *Relay) (int, error)
	DelById(id int) error
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

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (r *relayDal) GetRelayByType(relayType int) ([]*Relay, error) {
	var rs []*Relay
	err := r.db.Where("relay_type = ?", relayType).Find(&rs).Error
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
	err := r.db.Delete(&Relay{}, id).Error
	return err
}
