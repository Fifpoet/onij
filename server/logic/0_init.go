package logic

import "onij/inject"

type AllLogic struct {
	TodLogic   TodLogic
	RelayLogic RelayLogic
}

var app *inject.App

func Init() {
	app = inject.InitializeApp()
}
