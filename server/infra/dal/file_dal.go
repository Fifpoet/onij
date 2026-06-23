package dal

import (
	"context"
	"onij/model"
	"onij/util"
	"onij/util/cdb"
	"onij/util/logs"

	"gorm.io/gorm"
)

type FileDal interface {
	cdb.Interface[FileDal]

	Save(ctx context.Context, files ...*model.File) (int64, error)
	GetByIds(ctx context.Context, ids ...int64) ([]*model.File, error)
	GetByParentAndHash(ctx context.Context, parentId int64, hashes ...string) (*model.File, error)
	GetListByParentIdAndKeyword(ctx context.Context, parentId *int64, keyword *string, page util.Page) ([]*model.File, int32, error)
	GetFolderPathByParentId(ctx context.Context, parentId int64) ([]string, error)
	DeleteByIds(ctx context.Context, ids ...int64) ([]*model.File, error)
	DeleteById(ctx context.Context, id int64) error
}

type fileDal struct {
	*cdb.Dal[model.File, model.FileQuerier, model.FileUpdater]
}

func NewFileDal(db *cdb.DefaultProxy) FileDal {
	return &fileDal{
		cdb.NewDal[model.File, model.FileQuerier, model.FileUpdater](db),
	}
}

func (d *fileDal) With(tx *gorm.DB) FileDal {
	return &fileDal{d.Dal.With(tx)}
}

func (d *fileDal) Save(ctx context.Context, files ...*model.File) (int64, error) {
	return saveUpdatable(ctx, d, files)
}

func (d *fileDal) GetByIds(ctx context.Context, ids ...int64) ([]*model.File, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	res, err := d.QueryAll(ctx, d.Q().Id(ids).ToOptions()...)
	if err != nil {
		logs.Error("fileDal, GetByIds error = %v", err)
		return nil, err
	}
	return res, nil
}

func (d *fileDal) GetByParentAndHash(ctx context.Context, parentId int64, hashes ...string) (*model.File, error) {
	if len(hashes) == 0 {
		return nil, nil
	}
	res, err := d.QueryAll(ctx, d.Q().ParentId(parentId).Hash(hashes).ToOptions()...)
	if err != nil {
		logs.Error("fileDal, GetByParentAndHash error = %v", err)
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res[0], nil
}

func (d *fileDal) GetListByParentIdAndKeyword(ctx context.Context, parentId *int64, keyword *string, page util.Page) ([]*model.File, int32, error) {
	pid := int64(0)
	if parentId != nil {
		pid = *parentId
	}

	q := d.Q().ParentId(pid)
	if keyword != nil && len(*keyword) > 0 {
		q = q.Name(cdb.LIKE("%" + *keyword + "%"))
	}
	opts := q.ToOptions()

	cnt, err := d.Count(ctx, opts...)
	if err != nil {
		logs.Error("fileDal, GetListByParentIdAndKeyword count error = %v", err)
		return nil, 0, err
	}

	// 添加分页
	opts = append(opts, func(db *gorm.DB) *gorm.DB {
		return db.Offset(page.OffsetNum()).Limit(page.LimitNum())
	})

	res, err := d.QueryAll(ctx, opts...)
	if err != nil {
		logs.Error("fileDal, GetListByParentIdAndKeyword query error = %v", err)
		return nil, 0, err
	}
	return res, int32(cnt), nil
}

func (d *fileDal) GetFolderPathByParentId(ctx context.Context, parentId int64) ([]string, error) {
	// 递归查询file, 直到parentId为0
	var folders []string
	for parentId != 0 {
		res, err := d.QueryFirst(ctx, d.Q().Id(parentId).ToOptions()...)
		if err != nil {
			logs.Error("fileDal, GetFolderPathByParentId error = %v", err)
			return nil, err
		}
		if res == nil {
			break
		}
		parentId = res.ParentId
		folders = append([]string{res.Name}, folders...)
	}
	return folders, nil
}

func (d *fileDal) DeleteByIds(ctx context.Context, ids ...int64) ([]*model.File, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// 先获取要删除的文件
	res, err := d.GetByIds(ctx, ids...)
	if err != nil {
		return nil, err
	}

	// 执行删除
	_, err = d.Delete(ctx, d.Q().Id(ids).ToOptions()...)
	if err != nil {
		logs.Error("fileDal, DeleteByIds error = %v", err)
		return res, err
	}
	return res, nil
}

func (d *fileDal) DeleteById(ctx context.Context, id int64) error {
	_, err := d.Delete(ctx, d.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("fileDal, DeleteById error = %v", err)
		return err
	}
	return nil
}
