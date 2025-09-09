package dal

import (
	"context"
	"onij/model"
	"onij/util/cdb"
	"onij/util/logs"

	"gorm.io/gorm"
)

type ArtistDal interface {
	cdb.Interface[ArtistDal]

	Save(ctx context.Context, artists ...*model.Artist) (int64, error)
	GetById(ctx context.Context, id int64) (*model.Artist, error)
	GetByIds(ctx context.Context, ids ...int64) ([]*model.Artist, error)
	GetByName(ctx context.Context, name string) (*model.Artist, error)
	GetByKeyword(ctx context.Context, keyword string, offset, limit int32) ([]*model.Artist, int32, error)
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
	res, err := d.QueryAll(ctx, d.Q().Id(cdb.IN(ids)).ToOptions()...)
	if err != nil {
		logs.Error("artistDal, GetByIds error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *artistDal) GetByName(ctx context.Context, name string) (*model.Artist, error) {
	res, err := d.QueryFirst(ctx, d.Q().Name(name).ToOptions()...)
	if err != nil {
		logs.Error("artistDal, GetByName error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *artistDal) GetByKeyword(ctx context.Context, keyword string, offset, limit int32) ([]*model.Artist, int32, error) {
	opts := d.Q().Name(cdb.LIKE("%" + keyword + "%")).ToOptions()
	cnt, err := d.Count(ctx, opts...)
	if err != nil {
		logs.Error("artistDal, GetByKeyword count error = %v", err)
		return nil, 0, err
	}
	res, err := d.QuerySlice(ctx, int(offset), int(limit), opts...)
	if err != nil {
		logs.Error("artistDal, GetByKeyword query error = %v", err)
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
