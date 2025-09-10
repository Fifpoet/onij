package biz

import (
	"onij/model"
	"onij/model/api"
	"onij/util/boost/exp"
	"time"
)

type TagPrime struct {
	Id           int64
	ResourceId   *int64
	ResourceType *api.ResourceType
	TagBiz       *api.TagBiz
	TagGroup     *api.TagGroup
	TagType      *api.TagType
	TargetId     *int64
	Extra        *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

func (p *TagPrime) Model() *model.Tag {
	return &model.Tag{
		Id:           p.Id,
		ResourceId:   exp.ValueOrZero(p.ResourceId),
		ResourceType: int32(exp.ValueOrZero(p.ResourceType)),
		TagBiz:       int32(exp.ValueOrZero(p.TagBiz)),
		TagGroup:     int32(exp.ValueOrZero(p.TagGroup)),
		TagType:      int32(exp.ValueOrZero(p.TagType)),
		TargetId:     exp.ValueOrZero(p.TargetId),
		Extra:        exp.ValueOrZero(p.Extra),
	}
}

func Tag(p *TagPrime) *api.Tag {
	return &api.Tag{
		Id:           p.Id,
		ResourceId:   exp.ValueOrZero(p.ResourceId),
		ResourceType: exp.ValueOrZero(p.ResourceType),
		TagBiz:       exp.ValueOrZero(p.TagBiz),
		TagGroup:     exp.ValueOrZero(p.TagGroup),
		TagType:      exp.ValueOrZero(p.TagType),
		TargetId:     exp.ValueOrZero(p.TargetId),
		Extra:        exp.ValueOrZero(p.Extra),
	}
}

func TagToBiz(t *model.Tag) *TagPrime {
	if t == nil {
		return nil
	}
	return &TagPrime{
		Id:           t.Id,
		ResourceId:   exp.Ptr(t.ResourceId),
		ResourceType: (*api.ResourceType)(exp.Ptr(t.ResourceType)),
		TagBiz:       (*api.TagBiz)(exp.Ptr(t.TagBiz)),
		TagGroup:     (*api.TagGroup)(exp.Ptr(t.TagGroup)),
		TagType:      (*api.TagType)(exp.Ptr(t.TagType)),
		TargetId:     exp.Ptr(t.TargetId),
		Extra:        exp.Ptr(t.Extra),
		CreatedAt:    exp.Ptr(t.CreatedAt),
		UpdatedAt:    exp.Ptr(t.UpdatedAt),
	}
}
