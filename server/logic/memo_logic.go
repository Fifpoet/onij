package logic

import (
	"context"
	"errors"
	"onij/biz/prm"
	"onij/infra"
	"onij/infra/mysql"
	"onij/util"
)

type MemoLogic interface {
	Upload(ctx context.Context, param *prm.UploadMemoParam) (*prm.UploadMemoResult, error)
	GetDetail(ctx context.Context, param *prm.GetMemoDetailParam) (*prm.GetMemoDetailResult, error)
	GetList(ctx context.Context, param *prm.GetMemoListParam) (*prm.GetMemoListResult, error)
}

type memoLogic struct {
	*infra.AllInfra
}

func NewMemoLogic(i *infra.AllInfra) MemoLogic {
	return &memoLogic{
		AllInfra: i,
	}
}

func (l *memoLogic) Upload(ctx context.Context, param *prm.UploadMemoParam) (*prm.UploadMemoResult, error) {
	id := util.IdGen.Generate()
	err := l.MemoDal.Save(&mysql.Memo{
		Id:             id,
		Title:          param.Title,
		Content:        param.Content,
		OriginAt:       param.OriginAt,
		CoverFileId:    param.CoverFileId,
		ProfileContent: param.ProfileContent,
	})
	if err != nil {
		return nil, err
	}
	return &prm.UploadMemoResult{
		Id: id,
	}, nil
}

func (l *memoLogic) GetDetail(ctx context.Context, param *prm.GetMemoDetailParam) (*prm.GetMemoDetailResult, error) {
	memos, err := l.MemoDal.GetByIds(param.Id)
	if err != nil {
		return nil, err
	}
	if len(memos) == 0 {
		return nil, errors.New("memo not found")
	}
	memo := memos[0]

	return &prm.GetMemoDetailResult{
		Memo: memo,
	}, nil
}

func (l *memoLogic) GetList(ctx context.Context, param *prm.GetMemoListParam) (*prm.GetMemoListResult, error) {
	memos, err := l.MemoDal.GetListByKeyword(param.Keyword, param.Page)
	if err != nil {
		return nil, err
	}

	return &prm.GetMemoListResult{
		Memos: memos,
	}, nil
}
