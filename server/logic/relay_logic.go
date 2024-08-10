package logic

import (
	"onij/infra"
	"onij/infra/mysql"
)

type RelayLogic interface {
	GetRelayByType(relayType int) (res []*mysql.Relay, err error)
	DelById(id int) error
	Save(relay *mysql.Relay) (int, error)
}

type relayLogic struct {
	infra *infra.AllInfra
}

func NewRelayLogic(i *infra.AllInfra) RelayLogic {
	return &relayLogic{infra: i}
}

func (r *relayLogic) GetRelayByType(relayType int) (res []*mysql.Relay, err error) {
	rs, err := r.infra.GetRelayByType(relayType)
	return rs, err
}

func (r *relayLogic) DelById(id int) error {
	return r.infra.DelById(id)
}

func (r *relayLogic) Save(relay *mysql.Relay) (int, error) {
	return r.infra.Save(relay)
}
