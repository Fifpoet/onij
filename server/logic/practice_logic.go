package logic

import (
	"context"
	"fmt"
	"onij/biz/prm"
	"onij/infra"
	"onij/model"
	"onij/util"
)

type PracticeLogic interface {
	Upload(ctx context.Context, param *prm.UploadPracticeParam) (*prm.UploadPracticeResult, error)
	Delete(ctx context.Context, param *prm.DeletePracticeParam) (*prm.DeletePracticeResult, error)
	GetMonthly(ctx context.Context, param *prm.GetMonthlyPracticeParam) (*prm.GetMonthlyPracticeResult, error)
}

type practiceLogic struct {
	*infra.AllInfra
}

func NewPracticeLogic(i *infra.AllInfra) PracticeLogic {
	return &practiceLogic{AllInfra: i}
}

func (l *practiceLogic) Upload(ctx context.Context, param *prm.UploadPracticeParam) (*prm.UploadPracticeResult, error) {
	if param.PracticeAt <= 0 {
		return nil, fmt.Errorf("invalid practice_at")
	}

	var practice *model.Practice
	if param.Id != nil && *param.Id > 0 {
		existing, err := l.PracticeDal.GetById(ctx, *param.Id)
		if err != nil {
			return nil, err
		}
		if existing == nil {
			return nil, fmt.Errorf("practice not found")
		}
		practice = existing
	} else {
		practice = &model.Practice{
			Id: util.IdGen.Generate(),
		}
	}

	practice.PracticeType = param.PracticeType
	practice.Content = param.Content
	practice.PracticeAt = param.PracticeAt
	practice.Duration = param.Duration

	_, err := l.PracticeDal.Save(ctx, practice)
	if err != nil {
		return nil, err
	}
	return &prm.UploadPracticeResult{PracticeId: practice.Id}, nil
}

func (l *practiceLogic) Delete(ctx context.Context, param *prm.DeletePracticeParam) (*prm.DeletePracticeResult, error) {
	existing, err := l.PracticeDal.GetById(ctx, param.Id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("practice not found")
	}
	if err = l.PracticeDal.DeleteById(ctx, param.Id); err != nil {
		return nil, err
	}
	return &prm.DeletePracticeResult{}, nil
}

func (l *practiceLogic) GetMonthly(ctx context.Context, param *prm.GetMonthlyPracticeParam) (*prm.GetMonthlyPracticeResult, error) {
	if param.Year <= 0 || param.Month < 1 || param.Month > 12 {
		return nil, fmt.Errorf("invalid year or month")
	}
	start, end := prm.MonthRange(param.Year, param.Month, nil)
	practices, err := l.PracticeDal.GetByPracticeAtRange(ctx, start, end)
	if err != nil {
		return nil, err
	}
	return &prm.GetMonthlyPracticeResult{
		Days: prm.BuildMonthlyPracticeDays(practices, nil),
	}, nil
}
