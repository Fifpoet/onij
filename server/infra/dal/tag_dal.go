package dal

import (
	"context"
	"onij/model"
	"onij/util/cdb"
	"onij/util/logs"

	"gorm.io/gorm"
)

type TagDal interface {
	cdb.Interface[TagDal]

	Save(ctx context.Context, tags ...*model.Tag) (int64, error)
	GetByResource(ctx context.Context, resourceIds ...int64) ([]*model.Tag, error)
	GetByGroupType(ctx context.Context, tagGroup, tagType *int32) ([]*model.Tag, error)
	GetByResourceAndGroupType(ctx context.Context, tagGroups, tagTypes []int32, resourceIds ...int64) ([]*model.Tag, error)
	DeleteById(ctx context.Context, id int64) error
	DeleteByResourceAndType(ctx context.Context, resourceId int64, tagType int32) error
}

type tagDal struct {
	*cdb.Dal[model.Tag, model.TagQuerier, model.TagUpdater]
}

func NewTagDal(db *cdb.DefaultProxy) TagDal {
	return &tagDal{
		cdb.NewDal[model.Tag, model.TagQuerier, model.TagUpdater](db),
	}
}

func (d *tagDal) With(tx *gorm.DB) TagDal {
	return &tagDal{d.Dal.With(tx)}
}

func (d *tagDal) Save(ctx context.Context, tags ...*model.Tag) (int64, error) {
	return saveUpdatable(ctx, d, tags)
}

func (d *tagDal) GetByResource(ctx context.Context, resourceIds ...int64) ([]*model.Tag, error) {
	if len(resourceIds) == 0 {
		return nil, nil
	}
	res, err := d.QueryAll(ctx, d.Q().ResourceId(cdb.IN(resourceIds)).ToOptions()...)
	if err != nil {
		logs.Error("tagDal, GetByResource error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *tagDal) GetByGroupType(ctx context.Context, tagGroup, tagType *int32) ([]*model.Tag, error) {
	if tagGroup == nil && tagType == nil {
		return nil, nil
	}

	q := d.Q()
	if tagGroup != nil {
		q = q.TagGroup(*tagGroup)
	}
	if tagType != nil {
		q = q.TagType(*tagType)
	}

	res, err := d.QueryAll(ctx, q.ToOptions()...)
	if err != nil {
		logs.Error("tagDal, GetByGroupType error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *tagDal) GetByResourceAndGroupType(ctx context.Context, tagGroups, tagTypes []int32, resourceIds ...int64) ([]*model.Tag, error) {
	q := d.Q()

	if len(resourceIds) > 0 {
		q = q.ResourceId(cdb.IN(resourceIds))
	}
	if len(tagGroups) > 0 {
		q = q.TagGroup(cdb.IN(tagGroups))
	}
	if len(tagTypes) > 0 {
		q = q.TagType(cdb.IN(tagTypes))
	}

	res, err := d.QueryAll(ctx, q.ToOptions()...)
	if err != nil {
		logs.Error("tagDal, GetByResourceAndGroupType error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *tagDal) DeleteById(ctx context.Context, id int64) error {
	_, err := d.Delete(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("tagDal, DeleteById error = %v", err)
		return err
	}
	return nil
}

func (d *tagDal) DeleteByResourceAndType(ctx context.Context, resourceId int64, tagType int32) error {
	q := d.Q().ResourceId(resourceId).TagType(tagType)
	_, err := d.Delete(ctx, q.ToOptions()...)
	if err != nil {
		logs.Error("tagDal, DeleteByResourceAndType error = %v", err)
		return err
	}
	return nil
}
