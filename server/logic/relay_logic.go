package logic

import (
	"mime/multipart"
	"onij/enum"
	"onij/infra/mysql"
)

type RelayLogic interface {
	GetRelayByType(relayType int) (res []*mysql.Relay, err error)
	DelById(id int) (int, error)
	Save(relay *mysql.Relay, file *multipart.FileHeader) (int, error)
	PinRelay(id int) (int, error)
	GetRelayByPwd(pwd int) (res []*mysql.Relay, err error)
	GetRelayByPwdAndId(pwd, id int) (res *mysql.Relay, err error)
}

type relayLogic struct {
}

func NewRelayLogic() RelayLogic {
	return &relayLogic{}
}

func (r *relayLogic) GetRelayByType(relayType int) (res []*mysql.Relay, err error) {
	rs, err := app.GetRelayByType(relayType)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (r *relayLogic) DelById(id int) (int, error) {
	return 1, app.DelById(id)
}

func (r *relayLogic) Save(relay *mysql.Relay, file *multipart.FileHeader) (int, error) {
	// upload
	oss, err := app.FileDal.CreateFileFromForm(file, enum.BizRelay)
	if err != nil {
		return 0, err
	}
	relay.OssKey = oss

	return app.Save(relay)
}

func (r *relayLogic) PinRelay(id int) (int, error) {
	ori, err := app.GetByIds([]int{id})
	if err != nil || len(ori) == 0 {
		return 0, err
	}
	ori[0].Pin = !ori[0].Pin
	return app.Save(ori[0])
}

func (r *relayLogic) GetRelayByPwd(pwd int) (res []*mysql.Relay, err error) {
	return app.GetRelayByPwd(pwd)
}

func (r *relayLogic) GetRelayByPwdAndId(pwd, id int) (res *mysql.Relay, err error) {
	return app.GetRelayByPwdAndId(pwd, id)
}
