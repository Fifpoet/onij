package dal

import (
	"context"
	"onij/model"
	"onij/util/cdb"
	"onij/util/logs"

	"gorm.io/gorm"
)

type PracticeDal interface {
	cdb.Interface[PracticeDal]

	Save(ctx context.Context, practices ...*model.Practice) (int64, error)
	GetById(ctx context.Context, id int64) (*model.Practice, error)
	GetByPracticeAtRange(ctx context.Context, start, end int64) ([]*model.Practice, error)
	DeleteById(ctx context.Context, id int64) error
}

type practiceDal struct {
	*cdb.Dal[model.Practice, model.PracticeQuerier, model.PracticeUpdater]
}

func NewPracticeDal(db *cdb.DefaultProxy) PracticeDal {
	return &practiceDal{
		cdb.NewDal[model.Practice, model.PracticeQuerier, model.PracticeUpdater](db),
	}
}

func (d *practiceDal) With(tx *gorm.DB) PracticeDal {
	return &practiceDal{d.Dal.With(tx)}
}

func (d *practiceDal) Save(ctx context.Context, practices ...*model.Practice) (int64, error) {
	return saveUpdatable(ctx, d, practices)
}

func (d *practiceDal) GetById(ctx context.Context, id int64) (*model.Practice, error) {
	if id <= 0 {
		return nil, nil
	}
	res, err := d.QueryFirst(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("practiceDal, GetById error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *practiceDal) GetByPracticeAtRange(ctx context.Context, start, end int64) ([]*model.Practice, error) {
	opts := d.Q().WithAsc(model.Practice_PracticeAt).ToOptions()
	opts = append(opts, func(db *gorm.DB) *gorm.DB {
		return db.Where("practice_at >= ? AND practice_at < ?", start, end)
	})
	res, err := d.QueryAll(ctx, opts...)
	if err != nil {
		logs.Error("practiceDal, GetByPracticeAtRange error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *practiceDal) DeleteById(ctx context.Context, id int64) error {
	_, err := d.Delete(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("practiceDal, DeleteById error = %v", err)
		return err
	}
	return nil
}
