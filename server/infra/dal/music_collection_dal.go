package dal

import (
	"context"
	"onij/model"
	"onij/util/cdb"
	"onij/util/logs"

	"gorm.io/gorm"
)

type MusicCollectionDal interface {
	cdb.Interface[MusicCollectionDal]
	Save(ctx context.Context, items ...*model.MusicCollection) (int64, error)
	GetById(ctx context.Context, id int64) (*model.MusicCollection, error)
	ListByKeyword(ctx context.Context, keyword string) ([]*model.MusicCollection, error)
	DeleteById(ctx context.Context, id int64) error
}

type musicCollectionDal struct {
	*cdb.Dal[model.MusicCollection, model.MusicCollectionQuerier, model.MusicCollectionUpdater]
}

func NewMusicCollectionDal(db *cdb.DefaultProxy) MusicCollectionDal {
	return &musicCollectionDal{
		cdb.NewDal[model.MusicCollection, model.MusicCollectionQuerier, model.MusicCollectionUpdater](db),
	}
}

func (d *musicCollectionDal) With(tx *gorm.DB) MusicCollectionDal {
	return &musicCollectionDal{d.Dal.With(tx)}
}

func (d *musicCollectionDal) Save(ctx context.Context, items ...*model.MusicCollection) (int64, error) {
	return saveUpdatable(ctx, d, items)
}

func (d *musicCollectionDal) GetById(ctx context.Context, id int64) (*model.MusicCollection, error) {
	if id <= 0 {
		return nil, nil
	}
	res, err := d.QueryFirst(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("musicCollectionDal, GetById error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *musicCollectionDal) ListByKeyword(ctx context.Context, keyword string) ([]*model.MusicCollection, error) {
	opts := d.Q().WithDesc(model.MusicCollection_UpdatedAt).ToOptions()
	if keyword != "" {
		opts = append(opts, func(db *gorm.DB) *gorm.DB {
			return db.Where("search_text LIKE ?", "%"+keyword+"%")
		})
	}
	res, err := d.QueryAll(ctx, opts...)
	if err != nil {
		logs.Error("musicCollectionDal, ListByKeyword error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *musicCollectionDal) DeleteById(ctx context.Context, id int64) error {
	_, err := d.Delete(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("musicCollectionDal, DeleteById error = %v", err)
		return err
	}
	return nil
}
