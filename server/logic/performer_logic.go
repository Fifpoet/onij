package logic

import (
	"onij/infra/mysql"
)

type PerformerLogic interface {
	Save(performer *mysql.Performer) (int, error)
	DelById(id int) error
	GetByNameAndType(name string, performType int) ([]*mysql.Performer, error)
}

type performerLogic struct {
}

func NewPerformerLogic() PerformerLogic {
	return &performerLogic{}
}

func (p *performerLogic) Save(performer *mysql.Performer) (int, error) {
	return app.PerformerDal.Save(performer)
}

func (p *performerLogic) GetByNameAndType(name string, performType int) ([]*mysql.Performer, error) {
	return app.PerformerDal.GetByNameAndType(name, performType)
}

func (p *performerLogic) DelById(id int) error {
	return app.PerformerDal.DelById(id)
}
