//go:build wireinject

package inject

import (
	"github.com/google/wire"
	"onij/infra"
	"onij/infra/mysql"
)

var infraSet = wire.NewSet(
	mysql.NewMysqlCli,

	mysql.NewTagDal,
	mysql.NewRelayDal,
	mysql.NewFileDal,
	mysql.NewMusicDal,
	mysql.NewMetaDal,
	mysql.NewPerformerDal,

	wire.Struct(new(infra.AllInfra), "*"),
)

var allSet = wire.NewSet(
	infraSet,
	wire.Struct(new(App), "*"),
)

type App struct {
	*infra.AllInfra
}

func InitializeApp() *App {
	wire.Build(allSet)
	return &App{}
}
