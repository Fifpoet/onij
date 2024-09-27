package infra

import (
	"onij/infra/mysql"
)

// AllInfra 定义上层需要的Dal对象
type AllInfra struct {
	mysql.TagDal
	mysql.RelayDal
	mysql.FileDal
	mysql.MusicDal
	mysql.MetaDal
}
