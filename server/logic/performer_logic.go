package logic

import (
	"onij/handler/resq"
	"onij/infra/mysql"
	"onij/util/boost/collection/collext"
)

type PerformerLogic interface {
	Save(performer *mysql.Artist) (int, error)
	DelById(id int) error
	GetByNameAndType(name string, performType int) ([]*resq.GetPerformerModel, error)
}

type performerLogic struct {
}

func NewPerformerLogic() PerformerLogic {
	return &performerLogic{}
}

func (p *performerLogic) Save(performer *mysql.Artist) (int, error) {
	return app.PerformerDal.Save(performer)
}

func (p *performerLogic) GetByNameAndType(name string, performType int) ([]*resq.GetPerformerModel, error) {
	dbModel, err := app.PerformerDal.GetByNameAndType(name, performType)
	if err != nil {
		return nil, err
	}
	return collext.Pick(dbModel, func(db *mysql.Artist) *resq.GetPerformerModel {
		return &resq.GetPerformerModel{
			Id:   db.Id,
			Name: db.Name,
		}
	}), nil
}

func (p *performerLogic) DelById(id int) error {
	return app.PerformerDal.DelById(id)
}
