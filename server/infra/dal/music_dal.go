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

type MusicDal interface {
	cdb.Interface[MusicDal]

	Save(ctx context.Context, musics ...*model.Music) (int64, error)
	GetByIds(ctx context.Context, ids ...int64) ([]*model.Music, error)
	Search(ctx context.Context, artistIds []int64, tagTypes []int32, keyword string, pageInfo util.Page) ([]*model.Music, error)
	DeleteById(ctx context.Context, id int64) error
}

type musicDal struct {
	*cdb.Dal[model.Music, model.MusicQuerier, model.MusicUpdater]
}

func NewMusicDal(db *cdb.DefaultProxy) MusicDal {
	return &musicDal{
		cdb.NewDal[model.Music, model.MusicQuerier, model.MusicUpdater](db),
	}
}

func (d *musicDal) With(tx *gorm.DB) MusicDal {
	return &musicDal{d.Dal.With(tx)}
}

func (d *musicDal) Save(ctx context.Context, musics ...*model.Music) (int64, error) {
	return saveUpdatable(ctx, d, musics)
}

func (d *musicDal) GetById(ctx context.Context, id int64) (*model.Music, error) {
	res, err := d.QueryFirst(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("musicDal, GetById error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *musicDal) GetByIds(ctx context.Context, ids ...int64) ([]*model.Music, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	res, err := d.QueryAll(ctx, d.Q().Id(ids).ToOptions()...)
	if err != nil {
		logs.Error("musicDal, GetByIds error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *musicDal) GetByFullName(ctx context.Context, fullName string) ([]*model.Music, error) {
	res, err := d.QueryAll(ctx, d.Q().FullName(fullName).ToOptions()...)
	if err != nil {
		logs.Error("musicDal, GetByFullName error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *musicDal) Search(ctx context.Context, artistIds []int64, tagTypes []int32, keyword string, pageInfo util.Page) ([]*model.Music, error) {
	var conditions []cdb.Option

	if len(artistIds) > 0 {
		orConditions := collext.PickCombine(artistIds, func(id int64) []string {
			return []string{
				fmt.Sprintf("artist_ids LIKE '%%%d%%'", id),
				fmt.Sprintf("composer_ids LIKE '%%%d%%'", id),
				fmt.Sprintf("writer_ids LIKE '%%%d%%'", id),
			}
		})
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where(strings.Join(orConditions, " OR "))
		})
	}

	conditions = append(d.Q().FullName(cdb.LIKE(keyword)).ToOptions())

	// 处理tag_type条件（需要连表查询）
	if len(tagTypes) > 0 {
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN tag ON tag.resource_id = music.id AND tag.resource_type = ?", "music").
				Where("tag.tag_type IN ?", tagTypes)
		})
	}

	res, err := d.QuerySlice(ctx, pageInfo.PageNum(), pageInfo.LimitNum(), conditions...)
	if err != nil {
		logs.Error("musicDal, Search error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *musicDal) GetByTitleArtistPerType(ctx context.Context, title string, artistId, performType int64, offset, limit int32) ([]*model.Music, int32, error) {
	var conditions []cdb.Option

	if title != "" {
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("name LIKE ?", "%"+title+"%")
		})
	}
	if artistId != 0 {
		artistIdStr := fmt.Sprintf(",%d,", artistId)
		conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
			return db.Where("artist_ids LIKE ?", "%"+artistIdStr+"%")
		})
	}
	if performType != 0 {
		conditions = append(conditions, d.Q().PerformType(performType).ToOptions()...)
	}

	cnt, err := d.Count(ctx, conditions...)
	if err != nil {
		logs.Error("musicDal, GetByTitleArtistPerType count error = %v", err)
		return nil, 0, err
	}

	// 添加分页
	conditions = append(conditions, func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(offset)).Limit(int(limit))
	})

	res, err := d.QueryAll(ctx, conditions...)
	if err != nil {
		logs.Error("musicDal, GetByTitleArtistPerType query error = %v", err)
		return nil, 0, err
	}
	return res, int32(cnt), nil
}

func (d *musicDal) DeleteById(ctx context.Context, id int64) error {
	_, err := d.Delete(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("musicDal, DeleteById error = %v", err)
		return err
	}
	return nil
}
