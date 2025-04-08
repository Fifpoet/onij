//go:build wireinject

package inject

import (
	"github.com/google/wire"
	"onij/infra"
	"onij/infra/mysql"
	"onij/logic"
)

var infraSet = wire.NewSet(
	mysql.NewMysqlCli,

	mysql.NewTagDal,
	mysql.NewRelayDal,
	mysql.NewFileDal,
	mysql.NewMusicDal,
	mysql.NewArtistDal,

	wire.Struct(new(infra.AllInfra), "*"),
)

var logicSet = wire.NewSet(
	logic.NewMusicLogic,
	logic.NewFileLogic,
	wire.Struct(new(logic.AllLogic), "*"),
)

var allSet = wire.NewSet(
	infraSet,
	logicSet,
	wire.Struct(new(App), "*"),
)

type App struct {
	*infra.AllInfra
	*logic.AllLogic
}

func InitializeApp() *App {
	wire.Build(allSet)
	return &App{}
}
