package dal

import (
	"context"
	"gorm.io/gorm"
	"onij/model"
	"onij/util/cdb"
	"onij/util/logs"
)

type AlbumDal interface {
	cdb.Interface[AlbumDal]

	Save(ctx context.Context, albums ...*model.Album) (int64, error)
	GetByIds(ctx context.Context, id ...int64) ([]*model.Album, error)
	GetByArtistId(ctx context.Context, id int64, offset, limit int32) ([]*model.Album, int32, error)
}
type albumDal struct {
	*cdb.Dal[model.Album, model.AlbumQuerier, model.AlbumUpdater]
}

func NewAlbumDal(db *cdb.DefaultProxy) AlbumDal {
	return &albumDal{
		cdb.NewDal[model.Album, model.AlbumQuerier, model.AlbumUpdater](db),
	}
}

func (d *albumDal) With(tx *gorm.DB) AlbumDal {
	return &albumDal{d.Dal.With(tx)}
}

func (d *albumDal) Save(ctx context.Context, albums ...*model.Album) (int64, error) {
	return saveUpdatable(ctx, d, albums)
}

func (d *albumDal) GetByIds(ctx context.Context, id ...int64) ([]*model.Album, error) {
	res, err := d.QueryAll(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("albumDal, GetByIds error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *albumDal) GetByArtistId(ctx context.Context, id int64, offset, limit int32) ([]*model.Album, int32, error) {
	opts := d.Q().ArtistIds(cdb.LIKE(id)).ToOptions()
	cnt, err := d.Count(ctx, opts...)
	if err != nil {
		logs.Error("albumDal, GetById error = %v", err)
		return nil, 0, err
	}
	res, err := d.QuerySlice(ctx, int(offset), int(limit), opts...)
	if err != nil {
		logs.Error("albumDal, GetById error = %v", err)
		return nil, 0, err
	}
	return res, int32(cnt), nil
}
