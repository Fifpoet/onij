package dal

import (
	"context"
	"gorm.io/gorm"
	"onij/model"
	"onij/util/cdb"
	"onij/util/logs"
)

type AlbumMusicDal interface {
	cdb.Interface[AlbumMusicDal]

	Save(ctx context.Context, albums ...*model.AlbumMusic) (int64, error)
	GetById(ctx context.Context, id int64) (*model.AlbumMusic, error)
	GetByAlbumId(ctx context.Context, id int64) ([]*model.AlbumMusic, error)
}
type albumMusicDal struct {
	*cdb.Dal[model.AlbumMusic, model.AlbumMusicQuerier, model.AlbumMusicUpdater]
}

func NewAlbumMusicDal(db *cdb.DefaultProxy) AlbumMusicDal {
	return &albumMusicDal{
		cdb.NewDal[model.AlbumMusic, model.AlbumMusicQuerier, model.AlbumMusicUpdater](db),
	}
}

func (d *albumMusicDal) With(tx *gorm.DB) AlbumMusicDal {
	return &albumMusicDal{d.Dal.With(tx)}
}

func (d *albumMusicDal) Save(ctx context.Context, albums ...*model.AlbumMusic) (int64, error) {
	return saveUpdatable(ctx, d, albums)
}

func (d *albumMusicDal) GetById(ctx context.Context, id int64) (*model.AlbumMusic, error) {
	res, err := d.QueryFirst(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("albumMusicDal, GetById error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *albumMusicDal) GetByAlbumId(ctx context.Context, id int64) ([]*model.AlbumMusic, error) {
	res, err := d.QueryAll(ctx, d.Q().AlbumId(id).ToOptions()...)
	if err != nil {
		logs.Error("albumMusicDal, GetByAlbumId error = %v", err)
		return nil, err
	}
	return res, nil
}
