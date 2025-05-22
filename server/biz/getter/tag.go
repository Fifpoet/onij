package getter

import (
	"onij/infra/mysql"
	"onij/model/api"
)

func TagResourceId(t *mysql.Tag) int64 { return t.ResourceId }
func TagToDetail(t *mysql.Tag) *api.TagDetail {
	return &api.TagDetail{
		ResourceId:   t.ResourceId,
		ResourceType: t.ResourceType,
		TagBiz:       api.TagBiz(t.TagBiz),
		TagGroup:     api.TagGroup(t.TagGroup),
		TagType:      api.TagType(t.TagType),
		TargetId:     &t.TargetId,
		TargetType:   &t.TargetType,
		Extra:        t.Extra,
		ListShow:     t.ListShow,
	}
}
