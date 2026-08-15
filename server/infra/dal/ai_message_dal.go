package dal

import (
	"context"
	"onij/model"
	"onij/util/cdb"
	"onij/util/logs"

	"gorm.io/gorm"
)

type AiMessageDal interface {
	cdb.Interface[AiMessageDal]
	Save(ctx context.Context, items ...*model.AiMessage) (int64, error)
	ListBySessionId(ctx context.Context, sessionId int64) ([]*model.AiMessage, error)
	DeleteBySessionId(ctx context.Context, sessionId int64) error
}

type aiMessageDal struct {
	*cdb.Dal[model.AiMessage, model.AiMessageQuerier, model.AiMessageUpdater]
}

func NewAiMessageDal(db *cdb.DefaultProxy) AiMessageDal {
	return &aiMessageDal{cdb.NewDal[model.AiMessage, model.AiMessageQuerier, model.AiMessageUpdater](db)}
}

func (d *aiMessageDal) With(tx *gorm.DB) AiMessageDal {
	return &aiMessageDal{d.Dal.With(tx)}
}

func (d *aiMessageDal) Save(ctx context.Context, items ...*model.AiMessage) (int64, error) {
	return saveUpdatable(ctx, d, items)
}

func (d *aiMessageDal) ListBySessionId(ctx context.Context, sessionId int64) ([]*model.AiMessage, error) {
	if sessionId <= 0 {
		return nil, nil
	}
	res, err := d.QueryAll(ctx, d.Q().SessionId(sessionId).WithAsc(model.AiMessage_Id).ToOptions()...)
	if err != nil {
		logs.Error("aiMessageDal, ListBySessionId error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *aiMessageDal) DeleteBySessionId(ctx context.Context, sessionId int64) error {
	if sessionId <= 0 {
		return nil
	}
	_, err := d.Delete(ctx, d.Q().SessionId(sessionId).ToOptions()...)
	if err != nil {
		logs.Error("aiMessageDal, DeleteBySessionId error = %v", err)
		return err
	}
	return nil
}
