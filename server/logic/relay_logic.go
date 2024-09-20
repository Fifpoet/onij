package logic

import (
	"onij/infra/mysql"
)

type RelayLogic interface {
	GetRelayByType(relayType int) (res []*mysql.Relay, err error)
	DelById(id int) (int, error)
	DelByType(relayType int) error
	Save(relay *mysql.Relay) (int, error)
	PinRelay(id int) (int, error)
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

func (r *relayLogic) DelById(id int) (int, error) {
	return 1, app.DelById(id)
}

func (r *relayLogic) Save(relay *mysql.Relay) (int, error) {
	return app.Save(relay)
}

func (r *relayLogic) PinRelay(id int) (int, error) {
	ori, err := app.GetById([]int{id})
	if err != nil {
		return 0, err
	}
	ori[0].Pin = !ori[0].Pin
	return app.Save(ori[0])
}

func (r *relayLogic) DelByType(relayType int) error {
	return app.DelByType(relayType)
}
