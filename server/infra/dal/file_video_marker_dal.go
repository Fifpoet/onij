package dal

import (
	"context"
	"onij/model"
	"onij/util/cdb"
	"onij/util/logs"

	"gorm.io/gorm"
)

type FileVideoMarkerDal interface {
	cdb.Interface[FileVideoMarkerDal]
	Save(ctx context.Context, items ...*model.FileVideoMarker) (int64, error)
	ListByFileId(ctx context.Context, fileId int64) ([]*model.FileVideoMarker, error)
	DeleteById(ctx context.Context, id int64) error
	DeleteByFileId(ctx context.Context, fileId int64) error
}

type fileVideoMarkerDal struct {
	*cdb.Dal[model.FileVideoMarker, model.FileVideoMarkerQuerier, model.FileVideoMarkerUpdater]
}

func NewFileVideoMarkerDal(db *cdb.DefaultProxy) FileVideoMarkerDal {
	return &fileVideoMarkerDal{
		cdb.NewDal[model.FileVideoMarker, model.FileVideoMarkerQuerier, model.FileVideoMarkerUpdater](db),
	}
}

func (d *fileVideoMarkerDal) With(tx *gorm.DB) FileVideoMarkerDal {
	return &fileVideoMarkerDal{d.Dal.With(tx)}
}

func (d *fileVideoMarkerDal) Save(ctx context.Context, items ...*model.FileVideoMarker) (int64, error) {
	return saveUpdatable(ctx, d, items)
}

func (d *fileVideoMarkerDal) ListByFileId(ctx context.Context, fileId int64) ([]*model.FileVideoMarker, error) {
	if fileId <= 0 {
		return nil, nil
	}
	res, err := d.QueryAll(ctx, d.Q().FileId(fileId).WithAsc(model.FileVideoMarker_TimeMs).ToOptions()...)
	if err != nil {
		logs.Error("fileVideoMarkerDal, ListByFileId error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *fileVideoMarkerDal) DeleteById(ctx context.Context, id int64) error {
	if id <= 0 {
		return nil
	}
	_, err := d.Delete(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("fileVideoMarkerDal, DeleteById error = %v", err)
		return err
	}
	return nil
}

func (d *fileVideoMarkerDal) DeleteByFileId(ctx context.Context, fileId int64) error {
	if fileId <= 0 {
		return nil
	}
	_, err := d.Delete(ctx, d.Q().FileId(fileId).ToOptions()...)
	if err != nil {
		logs.Error("fileVideoMarkerDal, DeleteByFileId error = %v", err)
		return err
	}
	return nil
}
