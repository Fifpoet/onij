package dal

import (
	"context"
	"fmt"
	"onij/model"
	"onij/util"
	"onij/util/boost/collection/collext"
	"onij/util/cdb"
	"onij/util/logs"
	"strings"

	"gorm.io/gorm"
)

type AlbumDal interface {
	cdb.Interface[AlbumDal]

	Save(ctx context.Context, albums ...*model.Album) (int64, error)
	GetByIds(ctx context.Context, id ...int64) ([]*model.Album, error)
	GetByNameAndArtistId(ctx context.Context, keyword, name string, artistIds []int64, page util.Page) ([]*model.Album, int32, error)
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

func (d *albumDal) GetByNameAndArtistId(ctx context.Context, keyword, name string, artistIds []int64, page util.Page) ([]*model.Album, int32, error) {
	q := d.Q()
	if len(keyword) > 0 {
		q = d.Q().Name(cdb.LIKE(keyword))
	}
	if len(name) > 0 {
		q = d.Q().Name(name)
	}
	opts := q.ToOptions()
	if len(artistIds) > 0 {
		ors := collext.Pick(artistIds, func(a int64) string {
			return fmt.Sprintf("artist_ids LIKE '%%%d%%'", a)
		})
		opts = append(opts, func(db *gorm.DB) *gorm.DB {
			return db.Where(strings.Join(ors, " OR "))
		})
	}
	cnt, err := d.Count(ctx, opts...)
	if err != nil {
		logs.Error("albumDal, GetByKeywordAndArtistId count error = %v", err)
		return nil, 0, err
	}
	res, err := d.QuerySlice(ctx, page.OffsetNum(), page.LimitNum(), opts...)
	if err != nil {
		logs.Error("albumDal, GetByKeywordAndArtistId query error = %v", err)
		return nil, 0, err
	}
	return res, int32(cnt), nil
}
