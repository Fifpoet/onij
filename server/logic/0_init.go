package logic

import "onij/inject"

type AllLogic struct {
	TodLogic   TodLogic
	RelayLogic RelayLogic
	MusicLogic MusicLogic
}

var app *inject.App

func Init() {
	app = inject.InitializeApp()
}
