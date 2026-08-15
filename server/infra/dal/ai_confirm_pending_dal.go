package dal

import (
	"context"
	"errors"
	"onij/model"
	"onij/util/cdb"
	"onij/util/logs"
	"time"

	"gorm.io/gorm"
)

type AiConfirmPendingDal interface {
	cdb.Interface[AiConfirmPendingDal]
	Save(ctx context.Context, items ...*model.AiConfirmPending) (int64, error)
	GetById(ctx context.Context, id int64) (*model.AiConfirmPending, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type aiConfirmPendingDal struct {
	*cdb.Dal[model.AiConfirmPending, model.AiConfirmPendingQuerier, model.AiConfirmPendingUpdater]
}

func NewAiConfirmPendingDal(db *cdb.DefaultProxy) AiConfirmPendingDal {
	return &aiConfirmPendingDal{cdb.NewDal[model.AiConfirmPending, model.AiConfirmPendingQuerier, model.AiConfirmPendingUpdater](db)}
}

func (d *aiConfirmPendingDal) With(tx *gorm.DB) AiConfirmPendingDal {
	return &aiConfirmPendingDal{d.Dal.With(tx)}
}

func (d *aiConfirmPendingDal) Save(ctx context.Context, items ...*model.AiConfirmPending) (int64, error) {
	return saveUpdatable(ctx, d, items)
}

func (d *aiConfirmPendingDal) GetById(ctx context.Context, id int64) (*model.AiConfirmPending, error) {
	if id <= 0 {
		return nil, nil
	}
	res, err := d.QueryFirst(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logs.Error("aiConfirmPendingDal, GetById error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *aiConfirmPendingDal) UpdateStatus(ctx context.Context, id int64, status string) error {
	row, err := d.GetById(ctx, id)
	if err != nil {
		return err
	}
	if row == nil {
		return nil
	}
	row.Status = status
	row.UpdatedAt = time.Now()
	_, err = d.Save(ctx, row)
	return err
}
