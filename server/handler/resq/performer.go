package resq

import "onij/infra/mysql"

type UpsertPerformerReq struct {
	Id            int    `json:"id" form:"id"`
	Name          string `json:"name" form:"name"`
	PerformerType int    `json:"performer_type" form:"performer_type"`
}

func (u *UpsertPerformerReq) ToModel() *mysql.Performer {
	return &mysql.Performer{Id: u.Id, Name: u.Name, PerformerType: u.PerformerType}
}
