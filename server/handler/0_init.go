package handler

import "onij/logic"

type AllHandler struct {
	*BaseHandler
}

type BaseHandler struct {
	*logic.AllLogic
}

func NewBaseHandler(allLogic *logic.AllLogic) *BaseHandler {
	return &BaseHandler{
		AllLogic: allLogic,
	}
}
