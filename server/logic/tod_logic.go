package logic

import (
	"onij/logic/prm"
)

type TodLogic interface {
	GetWeeklyTodList() (res [][]prm.TodPrime, err error)
}

type todLogic struct {
}

func NewTodLogic() TodLogic {
	return &todLogic{}
}

func (t *todLogic) GetWeeklyTodList() (res [][]prm.TodPrime, err error) {

	return nil, nil
}
