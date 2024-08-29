// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package inject

import (
	"github.com/google/wire"
	"onij/infra"
	"onij/infra/mysql"
)

// Injectors from wire.go:

func InitializeApp() *App {
	db := mysql.NewMysqlCli()
	todDal := mysql.NewTodDal(db)
	tagDal := mysql.NewTagDal()
	relayDal := mysql.NewRelayDal(db)
	allInfra := &infra.AllInfra{
		TodDal:   todDal,
		TagDal:   tagDal,
		RelayDal: relayDal,
	}
	app := &App{
		AllInfra: allInfra,
	}
	return app
}

// wire.go:

var infraSet = wire.NewSet(mysql.NewMysqlCli, mysql.NewTagDal, mysql.NewTodDal, mysql.NewRelayDal, wire.Struct(new(infra.AllInfra), "*"))

var allSet = wire.NewSet(
	infraSet, wire.Struct(new(App), "*"),
)

type App struct {
	*infra.AllInfra
}
