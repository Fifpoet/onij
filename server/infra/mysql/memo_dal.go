package mysql

import "gorm.io/gorm"

type MemoDal interface {

}

type memoDal struct {
	db *gorm.DB
}

func NewMemoDal(db *gorm.DB) MemoDal {
	return &memoDal{db: db}
}

type Memo struct {
	
}