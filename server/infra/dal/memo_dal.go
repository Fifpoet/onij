package dal

import (
	"onij/util"
	"time"

	"gorm.io/gorm"
)

type MemoDal interface {
	Save(memo *Memo) error
	GetByIds(ids ...int64) ([]*Memo, error)
	GetListByKeyword(keyword string, page util.Page) ([]*Memo, error)
}

type memoDal struct {
	db *gorm.DB
}

func NewMemoDal(db *gorm.DB) MemoDal {
	return &memoDal{db: db}
}

type Memo struct {
	Id             int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	RootId         int64          `json:"root_id"`
	Title          string         `json:"title"`
	Content        string         `json:"content"`
	ProfileContent string         `json:"profile_content"`
	CoverFileId    int64          `json:"cover_file_id"`
	RelateFileId   int64          `json:"relate_file_id"`
	OriginAt       int64          `json:"origin_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at"`
}

func (m *memoDal) Save(memo *Memo) error {
	return m.db.Save(memo).Error
}

func (m *memoDal) GetByIds(ids ...int64) ([]*Memo, error) {
	var memos []*Memo
	err := m.db.Where("id IN (?)", ids).Find(&memos).Error
	if err != nil {
		return nil, err
	}
	return memos, nil
}

func (m *memoDal) GetListByKeyword(keyword string, page util.Page) ([]*Memo, error) {
	var memos []*Memo
	err := m.db.Where("title LIKE ?", "%"+keyword+"%").Offset(page.OffsetNum()).Limit(page.LimitNum()).Find(&memos).Error
	if err != nil {
		return nil, err
	}
	return memos, nil
}
