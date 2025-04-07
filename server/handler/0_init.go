package handler

import "onij/inject"

var ser *inject.App

func InitApp() {
	ser = inject.InitializeApp()
}
