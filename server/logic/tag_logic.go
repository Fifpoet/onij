package logic

import (
	"context"
	"onij/biz/prm"
	"onij/infra"
	"onij/model"
	"onij/util"
)

type TagLogic interface {
	Upload(ctx context.Context, param *prm.UploadTagParam) (*prm.UploadTagResult, error)
	Delete(ctx context.Context, param *prm.DeleteTagParam) (*prm.DeleteTagResult, error)
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
	err := l.TagDal.DeleteByResourceAndType(ctx, param.ResourceId, param.TagType)
	if err != nil {
		return nil, err
	}
	return &prm.DeleteTagResult{}, nil
}

func (l *tagLogic) Upload(ctx context.Context, param *prm.UploadTagParam) (*prm.UploadTagResult, error) {
	tag := &model.Tag{
		Id:           util.IdGen.Generate(),
		ResourceId:   param.ResourceId,
		ResourceType: param.ResourceType,
		TagBiz:       param.TagBiz,
		TagGroup:     param.TagGroup,
		TagType:      param.TagType,
		TargetId:     param.TargetId,
		TargetType:   param.TargetType,
		Extra:        param.Extra,
	}

	_, err := l.TagDal.Save(ctx, tag)
	if err != nil {
		return nil, err
	}

	return &prm.UploadTagResult{}, nil
}
