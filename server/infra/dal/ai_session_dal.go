package dal

import (
	"context"
	"errors"
	"onij/model"
	"onij/util/cdb"
	"onij/util/logs"

	"gorm.io/gorm"
)

type AiSessionDal interface {
	cdb.Interface[AiSessionDal]
	Save(ctx context.Context, items ...*model.AiSession) (int64, error)
	GetById(ctx context.Context, id int64) (*model.AiSession, error)
	GetByKind(ctx context.Context, kind int32) (*model.AiSession, error)
	ListRecent(ctx context.Context, limit int) ([]*model.AiSession, error)
	DeleteById(ctx context.Context, id int64) error
}

type aiSessionDal struct {
	*cdb.Dal[model.AiSession, model.AiSessionQuerier, model.AiSessionUpdater]
}

func NewAiSessionDal(db *cdb.DefaultProxy) AiSessionDal {
	return &aiSessionDal{cdb.NewDal[model.AiSession, model.AiSessionQuerier, model.AiSessionUpdater](db)}
}

func (d *aiSessionDal) With(tx *gorm.DB) AiSessionDal {
	return &aiSessionDal{d.Dal.With(tx)}
}

func (d *aiSessionDal) Save(ctx context.Context, items ...*model.AiSession) (int64, error) {
	return saveUpdatable(ctx, d, items)
}

func (d *aiSessionDal) GetById(ctx context.Context, id int64) (*model.AiSession, error) {
	if id <= 0 {
		return nil, nil
	}
	res, err := d.QueryFirst(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logs.Error("aiSessionDal, GetById error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *aiSessionDal) DeleteById(ctx context.Context, id int64) error {
	if id <= 0 {
		return nil
	}
	_, err := d.Delete(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("aiSessionDal, DeleteById error = %v", err)
		return err
	}
	return nil
}

func (d *aiSessionDal) GetByKind(ctx context.Context, kind int32) (*model.AiSession, error) {
	res, err := d.QueryFirst(ctx, d.Q().Kind(kind).ToOptions()...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logs.Error("aiSessionDal, GetByKind error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *aiSessionDal) ListRecent(ctx context.Context, limit int) ([]*model.AiSession, error) {
	if limit <= 0 {
		limit = 50
	}
	opts := d.Q().WithDesc(model.AiSession_UpdatedAt).ToOptions()
	opts = append(opts, func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	})
	res, err := d.QueryAll(ctx, opts...)
	if err != nil {
		logs.Error("aiSessionDal, ListRecent error = %v", err)
		return nil, err
	}
	return res, nil
}
