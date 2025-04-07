//go:build wireinject

package inject

import (
	"github.com/google/wire"
	"onij/handler"
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
	mysql.NewPerformerDal,

	wire.Struct(new(infra.AllInfra), "*"),
)

var logicSet = wire.NewSet(
	logic.NewMusicLogic,
	wire.Struct(new(logic.AllLogic), "*"),
)

var handleSet = wire.NewSet(
	handler.NewBaseHandler,
	wire.Struct(new(handler.AllHandler), "*"),
)

var allSet = wire.NewSet(
	infraSet,
	logicSet,
	handleSet,
	wire.Struct(new(App), "*"),
)

type App struct {
	*infra.AllInfra
}

func InitializeApp() *App {
	wire.Build(allSet)
	return &App{}
}
