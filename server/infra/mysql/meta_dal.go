package mysql

import (
	"gorm.io/gorm"
	"time"
)

type Meta struct {
	Id       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Category int    `json:"category" gorm:"uniqueIndex:uni_meta"`
	Value    int    `json:"value" gorm:"uniqueIndex:uni_meta"`
	Name     string `json:"name" `

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
