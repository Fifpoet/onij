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
}

type tagLogic struct {
	*infra.AllInfra
}

func NewTagLogic(i *infra.AllInfra) TagLogic {
	return &tagLogic{
		AllInfra: i,
	}
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
