//go:build wireinject

package inject

import (
	"github.com/google/wire"
	"onij/infra"
	"onij/infra/dal"
	"onij/logic"
	"onij/util/cdb"
)

var infraSet = wire.NewSet(
	dal.NewMysqlCli,
	cdb.NewDefaultProxy,

	dal.NewTagDal,
	dal.NewFileDal,
	dal.NewMusicDal,
	dal.NewAlbumDal,
	dal.NewAlbumMusicDal,
	dal.NewArtistDal,
	dal.NewPracticeDal,
	wire.Struct(new(infra.AllInfra), "*"),
)

var logicSet = wire.NewSet(
	logic.NewMusicLogic,
	logic.NewFileLogic,
	logic.NewAlbumLogic,
	logic.NewArtistLogic,
	logic.NewTagLogic,
	logic.NewPracticeLogic,
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
func InitLogicForTest() *logic.AllLogic {
	wire.Build(allSet)
	return &logic.AllLogic{}
}

func InitializeApp() *App {
	wire.Build(allSet)
	return &App{}
}
