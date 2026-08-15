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

type CollectionLogic interface {
	List(ctx context.Context, param *prm.ListCollectionParam) (*prm.ListCollectionResult, error)
	Create(ctx context.Context, param *prm.CreateCollectionParam) (*prm.MutationCollectionResult, error)
	Update(ctx context.Context, param *prm.UpdateCollectionParam) (*prm.MutationCollectionResult, error)
	Detail(ctx context.Context, param *prm.CollectionIdParam) (*prm.CollectionDetailResult, error)
	Delete(ctx context.Context, param *prm.CollectionIdParam) (*prm.OkCollectionResult, error)
	AddSongs(ctx context.Context, param *prm.AddCollectionSongsParam) (*prm.MutationCollectionResult, error)
	RemoveSong(ctx context.Context, param *prm.RemoveCollectionSongParam) (*prm.MutationCollectionResult, error)
}

type collectionLogic struct {
	*infra.AllInfra
}

func NewCollectionLogic(i *infra.AllInfra) CollectionLogic {
	return &collectionLogic{AllInfra: i}
}

func (l *collectionLogic) List(ctx context.Context, param *prm.ListCollectionParam) (*prm.ListCollectionResult, error) {
	keyword := strings.TrimSpace(param.Keyword)
	rows, err := l.MusicCollectionDal.ListByKeyword(ctx, keyword)
	if err != nil {
		return nil, err
	}
	items := make([]*prm.CollectionDTO, 0, len(rows))
	for _, row := range rows {
		ids, _ := l.songIdsOf(ctx, row.Id)
		items = append(items, prm.ToCollectionDTO(row, ids))
	}
	return &prm.ListCollectionResult{Items: items}, nil
}

func (l *collectionLogic) Create(ctx context.Context, param *prm.CreateCollectionParam) (*prm.MutationCollectionResult, error) {
	name := strings.TrimSpace(param.Name)
	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}
	now := time.Now()
	row := &model.MusicCollection{
		Id:         util.IdGen.Generate(),
		Name:       name,
		CoverUrl:   strings.TrimSpace(param.CoverUrl),
		SearchText: name,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if _, err := l.MusicCollectionDal.Save(ctx, row); err != nil {
		return nil, err
	}
	return &prm.MutationCollectionResult{Item: prm.ToCollectionDTO(row, []int64{})}, nil
}

func (l *collectionLogic) Update(ctx context.Context, param *prm.UpdateCollectionParam) (*prm.MutationCollectionResult, error) {
	row, err := l.MusicCollectionDal.GetById(ctx, param.Id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, fmt.Errorf("collection not found")
	}
	if param.Name != nil {
		name := strings.TrimSpace(*param.Name)
		if name == "" {
			return nil, fmt.Errorf("name is empty")
		}
		row.Name = name
	}
	if param.CoverUrl != nil {
		row.CoverUrl = strings.TrimSpace(*param.CoverUrl)
	}
	items, err := l.MusicCollectionItemDal.ListByCollectionId(ctx, row.Id)
	if err != nil {
		return nil, err
	}
	row.SearchText = prm.BuildCollectionSearchText(row.Name, items)
	row.UpdatedAt = time.Now()
	if _, err = l.MusicCollectionDal.Save(ctx, row); err != nil {
		return nil, err
	}
	ids := songIdsFromItems(items)
	return &prm.MutationCollectionResult{Item: prm.ToCollectionDTO(row, ids)}, nil
}

func (l *collectionLogic) Detail(ctx context.Context, param *prm.CollectionIdParam) (*prm.CollectionDetailResult, error) {
	row, err := l.MusicCollectionDal.GetById(ctx, param.Id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, fmt.Errorf("collection not found")
	}
	ids, err := l.songIdsOf(ctx, row.Id)
	if err != nil {
		return nil, err
	}
	return &prm.CollectionDetailResult{Item: prm.ToCollectionDTO(row, ids)}, nil
}

func (l *collectionLogic) Delete(ctx context.Context, param *prm.CollectionIdParam) (*prm.OkCollectionResult, error) {
	_ = l.MusicCollectionItemDal.DeleteByCollectionId(ctx, param.Id)
	if err := l.MusicCollectionDal.DeleteById(ctx, param.Id); err != nil {
		return nil, err
	}
	return &prm.OkCollectionResult{}, nil
}

func (l *collectionLogic) AddSongs(ctx context.Context, param *prm.AddCollectionSongsParam) (*prm.MutationCollectionResult, error) {
	row, err := l.MusicCollectionDal.GetById(ctx, param.Id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, fmt.Errorf("collection not found")
	}
	existing, err := l.MusicCollectionItemDal.ListByCollectionId(ctx, row.Id)
	if err != nil {
		return nil, err
	}
	nextOrder := int32(len(existing))
	now := time.Now()
	for _, meta := range param.Songs {
		if meta.SongId <= 0 {
			continue
		}
		found, err := l.MusicCollectionItemDal.GetByCollectionAndSong(ctx, row.Id, meta.SongId)
		if err != nil {
			return nil, err
		}
		if found != nil {
			continue
		}
		item := &model.MusicCollectionItem{
			Id:           util.IdGen.Generate(),
			CollectionId: row.Id,
			SongId:       meta.SongId,
			SongName:     strings.TrimSpace(meta.SongName),
			ArtistNames:  strings.TrimSpace(meta.ArtistNames),
			SortOrder:    nextOrder,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		nextOrder++
		if _, err = l.MusicCollectionItemDal.Save(ctx, item); err != nil {
			return nil, err
		}
	}
	return l.refreshAndReturn(ctx, row.Id)
}

func (l *collectionLogic) RemoveSong(ctx context.Context, param *prm.RemoveCollectionSongParam) (*prm.MutationCollectionResult, error) {
	row, err := l.MusicCollectionDal.GetById(ctx, param.Id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, fmt.Errorf("collection not found")
	}
	if err = l.MusicCollectionItemDal.DeleteByCollectionAndSong(ctx, param.Id, param.SongId); err != nil {
		return nil, err
	}
	return l.refreshAndReturn(ctx, row.Id)
}

func (l *collectionLogic) refreshAndReturn(ctx context.Context, id int64) (*prm.MutationCollectionResult, error) {
	row, err := l.MusicCollectionDal.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, fmt.Errorf("collection not found")
	}
	items, err := l.MusicCollectionItemDal.ListByCollectionId(ctx, id)
	if err != nil {
		return nil, err
	}
	row.SearchText = prm.BuildCollectionSearchText(row.Name, items)
	row.UpdatedAt = time.Now()
	if row.CoverUrl == "" {
		// leave empty; frontend may set first song cover later
	}
	if _, err = l.MusicCollectionDal.Save(ctx, row); err != nil {
		return nil, err
	}
	return &prm.MutationCollectionResult{Item: prm.ToCollectionDTO(row, songIdsFromItems(items))}, nil
}

func (l *collectionLogic) songIdsOf(ctx context.Context, collectionId int64) ([]int64, error) {
	items, err := l.MusicCollectionItemDal.ListByCollectionId(ctx, collectionId)
	if err != nil {
		return nil, err
	}
	return songIdsFromItems(items), nil
}

func songIdsFromItems(items []*model.MusicCollectionItem) []int64 {
	ids := make([]int64, 0, len(items))
	for _, it := range items {
		if it != nil && it.SongId > 0 {
			ids = append(ids, it.SongId)
		}
	}
	return ids
}
