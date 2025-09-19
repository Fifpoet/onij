package dal

import (
	"context"
	"onij/model"
	"onij/util"
	"onij/util/cdb"
	"onij/util/logs"

	"gorm.io/gorm"
)

type ArtistDal interface {
	cdb.Interface[ArtistDal]

	Save(ctx context.Context, artists ...*model.Artist) (int64, error)
	GetById(ctx context.Context, id int64) (*model.Artist, error)
	GetByIds(ctx context.Context, ids ...int64) ([]*model.Artist, error)
	GetByName(ctx context.Context, name string, keyword string, page util.Page) ([]*model.Artist, int32, error)
	DeleteById(ctx context.Context, id int64) error
}

type artistDal struct {
	*cdb.Dal[model.Artist, model.ArtistQuerier, model.ArtistUpdater]
}

func NewArtistDal(db *cdb.DefaultProxy) ArtistDal {
	return &artistDal{
		cdb.NewDal[model.Artist, model.ArtistQuerier, model.ArtistUpdater](db),
	}
}

func (d *artistDal) With(tx *gorm.DB) ArtistDal {
	return &artistDal{d.Dal.With(tx)}
}

func (d *artistDal) Save(ctx context.Context, artists ...*model.Artist) (int64, error) {
	return saveUpdatable(ctx, d, artists)
}

func (d *artistDal) GetById(ctx context.Context, id int64) (*model.Artist, error) {
	res, err := d.QueryFirst(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("artistDal, GetById error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *artistDal) GetByIds(ctx context.Context, ids ...int64) ([]*model.Artist, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	res, err := d.QueryAll(ctx, d.Q().Id(ids).ToOptions()...)
	if err != nil {
		logs.Error("artistDal, GetByIds error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *artistDal) GetByName(ctx context.Context, name string, keyword string, page util.Page) ([]*model.Artist, int32, error) {
	q := d.Q()
	if len(name) > 0 {
		q = q.Name(name)
	}
	if len(keyword) > 0 {
		q = q.Name(cdb.LIKE(keyword))
	}
	opts := q.ToOptions()
	cnt, err := d.Count(ctx, opts...)
	if err != nil {
		logs.Error("artistDal, GetByName count error = %v", err)
		return nil, 0, err
	}
	res, err := d.QuerySlice(ctx, page.PageNum(), page.LimitNum(), opts...)
	if err != nil {
		logs.Error("artistDal, GetByName query error = %v", err)
		return nil, 0, err
	}
	return res, int32(cnt), nil
}

func (d *artistDal) DeleteById(ctx context.Context, id int64) error {
	_, err := d.Delete(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("artistDal, DeleteById error = %v", err)
		return err
	}
	return nil
}
