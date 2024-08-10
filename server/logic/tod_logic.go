package logic

import (
	"onij/infra"
	"onij/logic/prm"
)

type TodLogic interface {
	GetWeeklyTodList() (res [][]prm.TodPrime, err error)
}

type todLogic struct {
	infra *infra.AllInfra
}

func NewTodLogic(i *infra.AllInfra) TodLogic {
	return &todLogic{infra: i}
}

func (t *todLogic) GetWeeklyTodList() (res [][]prm.TodPrime, err error) {

	return nil, nil
}
