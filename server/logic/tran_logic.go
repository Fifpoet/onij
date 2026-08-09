package logic

import (
	"context"
	"fmt"
	"onij/biz/prm"
	"onij/infra"
	"onij/model"
	"onij/util"
	"strings"
	"time"
)

type TranLogic interface {
	List(ctx context.Context, param *prm.ListTranParam) (*prm.ListTranResult, error)
	AddText(ctx context.Context, param *prm.AddTranTextParam) (*prm.AddTranResult, error)
	AddFile(ctx context.Context, param *prm.AddTranFileParam) (*prm.AddTranResult, error)
	TogglePin(ctx context.Context, param *prm.PinTranParam) (*prm.PinTranResult, error)
	Delete(ctx context.Context, param *prm.DeleteTranParam) (*prm.DeleteTranResult, error)
}

type tranLogic struct {
	*infra.AllInfra
	fileLogic FileLogic
}

func NewTranLogic(i *infra.AllInfra, fileLogic FileLogic) TranLogic {
	return &tranLogic{AllInfra: i, fileLogic: fileLogic}
}

func (l *tranLogic) List(ctx context.Context, _ *prm.ListTranParam) (*prm.ListTranResult, error) {
	rows, err := l.TranItemDal.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*prm.TranItemDTO, 0, len(rows))
	for _, row := range rows {
		items = append(items, prm.ToTranItemDTO(row))
	}
	prm.SortTranDTOs(items)
	return &prm.ListTranResult{Items: items}, nil
}

func (l *tranLogic) AddText(ctx context.Context, param *prm.AddTranTextParam) (*prm.AddTranResult, error) {
	content := strings.TrimSpace(param.Content)
	if content == "" {
		return nil, fmt.Errorf("content is empty")
	}
	now := time.Now()
	row := &model.TranItem{
		Id:        util.IdGen.Generate(),
		Kind:      model.TranKindText,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := l.TranItemDal.Save(ctx, row); err != nil {
		return nil, err
	}
	return &prm.AddTranResult{Item: prm.ToTranItemDTO(row)}, nil
}

func (l *tranLogic) AddFile(ctx context.Context, param *prm.AddTranFileParam) (*prm.AddTranResult, error) {
	if param.FileId <= 0 {
		return nil, fmt.Errorf("invalid file_id")
	}
	name := strings.TrimSpace(param.Name)
	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}
	now := time.Now()
	row := &model.TranItem{
		Id:        util.IdGen.Generate(),
		Kind:      model.TranKindFile,
		FileId:    param.FileId,
		Name:      name,
		Size:      param.Size,
		Format:    param.Format,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := l.TranItemDal.Save(ctx, row); err != nil {
		return nil, err
	}
	return &prm.AddTranResult{Item: prm.ToTranItemDTO(row)}, nil
}

func (l *tranLogic) TogglePin(ctx context.Context, param *prm.PinTranParam) (*prm.PinTranResult, error) {
	row, err := l.TranItemDal.GetById(ctx, param.Id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, fmt.Errorf("tran item not found")
	}
	if row.Pinned {
		row.Pinned = false
		row.PinnedAt = 0
	} else {
		row.Pinned = true
		row.PinnedAt = time.Now().UnixMilli()
	}
	if _, err = l.TranItemDal.Save(ctx, row); err != nil {
		return nil, err
	}
	return &prm.PinTranResult{Item: prm.ToTranItemDTO(row)}, nil
}

func (l *tranLogic) Delete(ctx context.Context, param *prm.DeleteTranParam) (*prm.DeleteTranResult, error) {
	row, err := l.TranItemDal.GetById(ctx, param.Id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return &prm.DeleteTranResult{}, nil
	}
	if row.Kind == model.TranKindFile && row.FileId > 0 {
		_, _ = l.fileLogic.Delete(ctx, &prm.DeleteFileParam{FileId: row.FileId})
	}
	if err = l.TranItemDal.DeleteById(ctx, param.Id); err != nil {
		return nil, err
	}
	return &prm.DeleteTranResult{}, nil
}
