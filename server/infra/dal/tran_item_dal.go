package dal

import (
	"context"
	"onij/model"
	"onij/util/cdb"
	"onij/util/logs"

	"gorm.io/gorm"
)

type TranItemDal interface {
	cdb.Interface[TranItemDal]
	Save(ctx context.Context, items ...*model.TranItem) (int64, error)
	GetById(ctx context.Context, id int64) (*model.TranItem, error)
	ListAll(ctx context.Context) ([]*model.TranItem, error)
	DeleteById(ctx context.Context, id int64) error
}

type tranItemDal struct {
	*cdb.Dal[model.TranItem, model.TranItemQuerier, model.TranItemUpdater]
}

func NewTranItemDal(db *cdb.DefaultProxy) TranItemDal {
	return &tranItemDal{
		cdb.NewDal[model.TranItem, model.TranItemQuerier, model.TranItemUpdater](db),
	}
}

func (d *tranItemDal) With(tx *gorm.DB) TranItemDal {
	return &tranItemDal{d.Dal.With(tx)}
}

func (d *tranItemDal) Save(ctx context.Context, items ...*model.TranItem) (int64, error) {
	return saveUpdatable(ctx, d, items)
}

func (d *tranItemDal) GetById(ctx context.Context, id int64) (*model.TranItem, error) {
	if id <= 0 {
		return nil, nil
	}
	res, err := d.QueryFirst(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("tranItemDal, GetById error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *tranItemDal) ListAll(ctx context.Context) ([]*model.TranItem, error) {
	// QueryAll 不允许空条件；用排序作为占位条件拉全量
	res, err := d.QueryAll(ctx, d.Q().WithDesc(model.TranItem_CreatedAt).ToOptions()...)
	if err != nil {
		logs.Error("tranItemDal, ListAll error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *tranItemDal) DeleteById(ctx context.Context, id int64) error {
	_, err := d.Delete(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("tranItemDal, DeleteById error = %v", err)
		return err
	}
	return nil
}
