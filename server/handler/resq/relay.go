package resq

import "onij/infra/mysql"

type UpsertRelayReq struct {
}

func (u *UpsertRelayReq) ToModel() *mysql.Relay {
	return nil
}
