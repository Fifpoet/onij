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
	mysql.NewTodDal,
	mysql.NewRelayDal,

	wire.Struct(new(infra.AllInfra), "*"),
)

var logicSet = wire.NewSet(
	logic.NewTodLogic,
	logic.NewRelayLogic,

	wire.Struct(new(logic.AllLogic), "*"),
)

var handlerSet = wire.NewSet(
	handler.NewRelayHandler,

	wire.Struct(new(handler.AllHandler), "*"),
)

var allSet = wire.NewSet(
	infraSet,
	logicSet,
	handlerSet,
	wire.Struct(new(App), "*"),
)

type App struct {
	*infra.AllInfra
}

func InitializeApp() *App {
	wire.Build(allSet)
	return &App{}
}
