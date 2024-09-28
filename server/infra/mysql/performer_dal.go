package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type PerformerDal interface {
	GetByIds(id []int) ([]*Performer, error)
	GetByName(name string) ([]*Performer, error)
	GetByNameAndType(name string, performType int) ([]*Performer, error)
	Save(performer *Performer) (int, error)
	DelById(id int) error
}

type performerDal struct {
	db *gorm.DB
}

func NewPerformerDal(db *gorm.DB) PerformerDal {
	return &performerDal{db: db}
}

type Performer struct {
	Id            int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string         `json:"name"`
	PerformerType int            `json:"performer_type"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}

func (p *performerDal) GetByIds(id []int) ([]*Performer, error) {
	var res []*Performer
	err := p.db.Where("id IN ?", id).Order("created_at desc").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *performerDal) GetByName(name string) ([]*Performer, error) {
	var res []*Performer
	err := p.db.Where("name LIKE ?", "%"+name+"%").Order("created_at desc").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *performerDal) GetByNameAndType(name string, performType int) ([]*Performer, error) {
	var performers []*Performer
	err := p.db.Where("name LIKE ? AND performer_type = ?", "%"+name+"%", performType).Find(&performers).Error
	if err != nil {
		return nil, err
	}
	return performers, nil
}

func (p *performerDal) Save(performer *Performer) (int, error) {
	err := p.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(performer).Error
	if err != nil {
		return 0, err
	}
	return performer.Id, nil
}

func (p *performerDal) DelById(id int) error {
	err := p.db.Where("id = ?", id).Delete(&Performer{}).Error
	return err
}
