package dal

import (
	"context"
	"errors"
	"onij/model"
	"onij/util/cdb"
	"onij/util/logs"

	"gorm.io/gorm"
)

type MusicCollectionItemDal interface {
	cdb.Interface[MusicCollectionItemDal]
	Save(ctx context.Context, items ...*model.MusicCollectionItem) (int64, error)
	ListByCollectionId(ctx context.Context, collectionId int64) ([]*model.MusicCollectionItem, error)
	GetByCollectionAndSong(ctx context.Context, collectionId, songId int64) (*model.MusicCollectionItem, error)
	DeleteByCollectionAndSong(ctx context.Context, collectionId, songId int64) error
	DeleteByCollectionId(ctx context.Context, collectionId int64) error
}

type musicCollectionItemDal struct {
	*cdb.Dal[model.MusicCollectionItem, model.MusicCollectionItemQuerier, model.MusicCollectionItemUpdater]
}

func NewMusicCollectionItemDal(db *cdb.DefaultProxy) MusicCollectionItemDal {
	return &musicCollectionItemDal{
		cdb.NewDal[model.MusicCollectionItem, model.MusicCollectionItemQuerier, model.MusicCollectionItemUpdater](db),
	}
}

func (d *musicCollectionItemDal) With(tx *gorm.DB) MusicCollectionItemDal {
	return &musicCollectionItemDal{d.Dal.With(tx)}
}

func (d *musicCollectionItemDal) Save(ctx context.Context, items ...*model.MusicCollectionItem) (int64, error) {
	return saveUpdatable(ctx, d, items)
}

func (d *musicCollectionItemDal) ListByCollectionId(ctx context.Context, collectionId int64) ([]*model.MusicCollectionItem, error) {
	if collectionId <= 0 {
		return nil, nil
	}
	res, err := d.QueryAll(ctx, d.Q().CollectionId(collectionId).WithAsc(model.MusicCollectionItem_SortOrder, model.MusicCollectionItem_Id).ToOptions()...)
	if err != nil {
		logs.Error("musicCollectionItemDal, ListByCollectionId error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *musicCollectionItemDal) GetByCollectionAndSong(ctx context.Context, collectionId, songId int64) (*model.MusicCollectionItem, error) {
	if collectionId <= 0 || songId <= 0 {
		return nil, nil
	}
	res, err := d.QueryFirst(ctx, d.Q().CollectionId(collectionId).SongId(songId).ToOptions()...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logs.Error("musicCollectionItemDal, GetByCollectionAndSong error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *musicCollectionItemDal) DeleteByCollectionAndSong(ctx context.Context, collectionId, songId int64) error {
	_, err := d.Delete(ctx, d.Q().CollectionId(collectionId).SongId(songId).ToOptions()...)
	if err != nil {
		logs.Error("musicCollectionItemDal, DeleteByCollectionAndSong error = %v", err)
		return err
	}
	return nil
}

func (d *musicCollectionItemDal) DeleteByCollectionId(ctx context.Context, collectionId int64) error {
	_, err := d.Delete(ctx, d.Q().CollectionId(collectionId).ToOptions()...)
	if err != nil {
		logs.Error("musicCollectionItemDal, DeleteByCollectionId error = %v", err)
		return err
	}
	return nil
}
