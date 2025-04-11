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
	mysql.NewFileDal,
	mysql.NewMusicDal,
	mysql.NewAlbumDal,
	mysql.NewArtistDal,

	wire.Struct(new(infra.AllInfra), "*"),
)

var logicSet = wire.NewSet(
	logic.NewMusicLogic,
	logic.NewFileLogic,
	logic.NewAlbumLogic,
	logic.NewArtistLogic,
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

func InitDalForTest() *infra.AllInfra {
	wire.Build(infraSet)
	return &infra.AllInfra{}
}

func InitializeApp() *App {
	wire.Build(allSet)
	return &App{}
}
