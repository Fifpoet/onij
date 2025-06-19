package infra

import (
	"onij/infra/mysql"
)

// AllInfra 定义上层需要的Dal对象
type AllInfra struct {
	mysql.TagDal
	mysql.FileDal
	mysql.MusicDal
	mysql.AlbumDal
	mysql.ArtistDal
	mysql.MemoDal
}
