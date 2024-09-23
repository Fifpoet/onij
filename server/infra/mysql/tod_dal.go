package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"onij/util"
	"time"
)

type TodDal interface {
}

type todDal struct {
	db *gorm.DB
}

func NewTodDal(db *gorm.DB) TodDal {
	return &todDal{db: db}
}

type Tod struct {
	Id      int    `json:"id" gorm:"primaryKey"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Score   int    `json:"score"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// GetWeeklyTodList 本周的Tod
func (t *todDal) GetWeeklyTodList() (res [][]Tod, err error) {
	start := util.GetThisWeekDates()[0]
	score := util.GetDayScore(start)
	res = make([][]Tod, 7)

	var ts []Tod
	err = t.db.Where("score BETWEEN ? AND ?", score, score+6).Find(&ts).Error
	if err != nil {
		return nil, err
	}

	for _, t := range ts {
		res[t.Score-score] = append(res[t.Score-score], t)
	}
	return res, nil
}

// UpsertTod 插入或更新tod
func (t *todDal) UpsertTod(tod Tod) error {
	err := t.db.Clauses(clause.OnConflict{
		UpdateAll: true, // 更新所有列
	}).Create(&tod).Error
	if err != nil {
		log.Printf("UpsertTod, insert failed: {%v}\n", err)
	}
	return nil
}

// DelTod 删除tod
func (t *todDal) DelTod(id int) error {
	err := t.db.Delete(Tod{Id: id}).Error
	if err != nil {
		log.Printf("DelTod, del failed: {%v}\n", err)
	}
	return nil
}
