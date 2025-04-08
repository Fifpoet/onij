package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArtistDal interface {
	GetByIds(id ...int) ([]*Artist, error)
	GetByName(name string) ([]*Artist, error)
	GetByNameAndType(name string, performType int) ([]*Artist, error)
	Save(performer *Artist) (int, error)
	DelById(id int) error
}

type artistDal struct {
	db *gorm.DB
}

func NewArtistDal(db *gorm.DB) ArtistDal {
	return &artistDal{db: db}
}

type Artist struct {
	Id            int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string         `json:"name"`
	PerformerType int            `json:"performer_type"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}

func (p *artistDal) GetByIds(id ...int) ([]*Artist, error) {
	var res []*Artist
	err := p.db.Where("id IN ?", id).Order("created_at desc").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *artistDal) GetByName(name string) ([]*Artist, error) {
	var res []*Artist
	err := p.db.Where("name LIKE ?", "%"+name+"%").Order("created_at desc").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *artistDal) GetByNameAndType(name string, performType int) ([]*Artist, error) {
	var performers []*Artist
	err := p.db.Where("name LIKE ? AND performer_type = ?", "%"+name+"%", performType).Find(&performers).Error
	if err != nil {
		return nil, err
	}
	return performers, nil
}

func (p *artistDal) Save(performer *Artist) (int, error) {
	err := p.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(performer).Error
	if err != nil {
		return 0, err
	}
	return performer.Id, nil
}

func (p *artistDal) DelById(id int) error {
	err := p.db.Where("id = ?", id).Delete(&Artist{}).Error
	return err
}
