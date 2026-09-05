//go:build wireinject

package inject

import (
	"onij/infra"
	"onij/infra/dal"
	"onij/infra/uvr"
	"onij/logic"
	"onij/util/cdb"

	"github.com/google/wire"
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
	dal.NewTranItemDal,
	dal.NewAiSessionDal,
	dal.NewAiMessageDal,
	dal.NewAiConfirmPendingDal,
	dal.NewMusicCollectionDal,
	dal.NewMusicCollectionItemDal,
	dal.NewFileVideoMarkerDal,
	wire.Struct(new(infra.AllInfra), "*"),
)

var logicSet = wire.NewSet(
	uvr.NewClient,
	logic.NewMusicLogic,
	logic.NewFileLogic,
	logic.NewAlbumLogic,
	logic.NewArtistLogic,
	logic.NewTagLogic,
	logic.NewPracticeLogic,
	logic.NewUvrLogic,
	logic.NewTranLogic,
	logic.NewAiLogic,
	logic.NewCollectionLogic,
	logic.NewFileVideoLogic,
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
