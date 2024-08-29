package logic

import (
	"onij/infra/mysql"
)

type RelayLogic interface {
	GetRelayByType(relayType int) (res []*mysql.Relay, err error)
	DelById(id int) error
	Save(relay *mysql.Relay) (int, error)
}

type relayLogic struct {
}

func NewRelayLogic() RelayLogic {
	return &relayLogic{}
}

func (r *relayLogic) GetRelayByType(relayType int) (res []*mysql.Relay, err error) {
	rs, err := app.GetRelayByType(relayType)
	return rs, err
}

func (r *relayLogic) DelById(id int) error {
	return app.DelById(id)
}

func (r *relayLogic) Save(relay *mysql.Relay) (int, error) {
	return app.Save(relay)
}
