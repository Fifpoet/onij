package getter

import (
	"onij/model"
	"onij/model/api"
)

func TagResourceId(t *model.Tag) int64 { return t.ResourceId }
func TagToDetail(t *model.Tag) *api.Tag {
	if t == nil {
		return nil
	}
	return &api.Tag{
		ResourceId:   t.ResourceId,
		ResourceType: api.ResourceType(t.ResourceType),
		TagBiz:       api.TagBiz(t.TagBiz),
		TagGroup:     api.TagGroup(t.TagGroup),
		TagType:      api.TagType(t.TagType),
		TargetId:     t.TargetId,
		TargetType:   t.TargetType,
		Extra:        t.Extra,
	}
}
