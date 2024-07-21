package infra

import (
	"gorm.io/gorm"
	"onij/infra/mysql"
)

type AllInfra struct {
	db *gorm.DB
}

func NewAllInfra() *AllInfra {
	mysql.InitMysql()
	return &AllInfra{db: mysql.Db}
}
