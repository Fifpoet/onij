package logic

import (
	"context"
	"onij/biz/prm"
	"onij/infra"
	"onij/infra/mysql"
	"onij/util"
	"onij/util/boost/exp"
)

type TagLogic interface {
	Upload(ctx context.Context, prm *prm.UploadTagParam) (*prm.UploadTagResult, error)
	GetList(ctx context.Context, prm *prm.GetTagListParam) (*prm.GetTagListResult, error)
	Delete(ctx context.Context, prm *prm.DeleteTagParam) (*prm.DeleteTagResult, error)
}

type tagLogic struct {
	*infra.AllInfra
}

func NewTagLogic(i *infra.AllInfra) TagLogic {
	return &tagLogic{
		AllInfra: i,
	}
}

func (l *tagLogic) Delete(ctx context.Context, param *prm.DeleteTagParam) (*prm.DeleteTagResult, error) {
	err := l.TagDal.Delete(param.ResourceId, param.TagType)
	if err!= nil {
		return nil, err
	}
	return &prm.DeleteTagResult{}, nil
}

func (l *tagLogic) GetList(ctx context.Context, param *prm.GetTagListParam) (*prm.GetTagListResult, error) {
	tags, err := l.TagDal.GetByResourceAndGroupType(param.TagGroups, param.TagTypes, param.ResourceIds...)
	if err != nil {
		return nil, err
	}
	return &prm.GetTagListResult{
		Tags: tags,
	}, nil
}

func (l *tagLogic) Upload(ctx context.Context, param *prm.UploadTagParam) (*prm.UploadTagResult, error) {
	err := l.TagDal.Save(&mysql.Tag{
		Id:           util.IdGen.Generate(),
		ResourceId:   param.ResourceId,
		ResourceType: param.ResourceType,
		TagBiz:       int32(param.TagBiz),
		TagGroup:     int32(param.TagGroup),
		TagType:      int32(param.TagType),
		TargetId:     exp.ValueOrZero(param.TargetId),
		TargetType:   exp.ValueOrZero(param.TargetType),
		ListShow:     param.ListShow,
		Extra:        param.Extra,
	})
	if err != nil {
		return nil, err
	}
	return &prm.UploadTagResult{
	}, nil
}
