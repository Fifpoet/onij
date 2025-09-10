package dal

import (
	"context"
	"onij/biz/errdef"
	"onij/util/logs"

	"fmt"
	"log"
	"onij/util/boost/ccmp"
	"onij/util/boost/collection/collext"
	"onij/util/boost/exp"
	"onij/util/cdb"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/schema"
)

func NewMysqlCli() *gorm.DB {
	ip := os.Getenv("ip")
	mysqlPwd := os.Getenv("mypwd")
	dsn := "root:%s@tcp(%s:3306)/footprint?charset=utf8mb4&charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	Db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       fmt.Sprintf(dsn, mysqlPwd, ip), // DSN data source name
		DefaultStringSize:         256,                            // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                           // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                           // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                           // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                          // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 完全关闭日志
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "onij_", // 表名前缀，`User`表为`dp_user`
			SingularTable: true,    // 使用单数表名，启用该选项后，`User` 表将是`user`
		},
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return Db
}

const (
	FieldValue_RootItemRootId              = 0 // 根节点的root_id
	FieldValue_RootItemParentId            = 0 // 根节点的parent_id
	FieldValue_Available                   = 1 // 资源可用标记
	FieldValue_Unavailable                 = 0 // 资源不可用标记
	FieldValue_PaperQuestionIsRelatedFalse = 0 // 题卷题目非关联
)

const (
	maxBatchSize = 500 // 批量上限
)

var (
	timeColumnsSets = map[string]struct{}{"created_at": {}, "updated_at": {}}
)

type idQuerier[Q cdb.Querier] interface {
	Id(any) Q
}

type idDal[U cdb.Updater, Q cdb.Querier, DQ idQuerier[Q]] interface {
	Q() DQ
	U() U
	Delete(ctx context.Context, opts ...cdb.Option) (int64, error)
	Update(ctx context.Context, u U, conditions ...cdb.Option) (int64, error)
}

type txDal[U cdb.Updater] interface {
	U() U
	W(ctx context.Context, opts ...cdb.Option) *gorm.DB
}

func deleteByIds[U cdb.Updater, Q cdb.Querier, DQ idQuerier[Q]](ctx context.Context, dal idDal[U, Q, DQ], ids []int64) (int64, error) {
	if len(ids) == 0 {
		logs.Info("dal.deleteByIds[%T]: empty ids, skip", dal)
		return 0, nil
	}

	affected, err := dal.Delete(ctx, dal.Q().Id(ids).ToOptions()...)
	if err != nil {
		logs.Error("dal.DeleteByIds[%T]: err = %v", dal, err)
		return 0, err
	}

	return affected, nil
}

func updateById[U cdb.Updater, Q cdb.Querier, DQ idQuerier[Q]](ctx context.Context, dal idDal[U, Q, DQ], hook func(U), id int64) (bool, error) {
	if id <= 0 {
		logs.Error("dal.updateById[%T]: invalid id", dal)
		return false, errdef.ErrDalEntityIdInvalid
	}

	u := dal.U()
	hook(u)
	if u.IsEmpty() {
		logs.Info("dal.updateById[%T]: non update items, id = %d", dal, id)
		return false, nil
	}

	n, err := dal.Update(ctx, u, dal.Q().Id(id).ToOptions()...)
	if err != nil {
		logs.Error("dal.updateById[%T]: err = %v, id = %d", dal, err, id)
		return false, err
	}

	return n != 0, nil
}

func updates[M cdb.TableModel, U cdb.Updater](ctx context.Context, dal txDal[U], getId func(*M) int64, hook func(*M, U), models []*M) (affected int64, err error) {
	zero := exp.Zero[M]()
	if len(models) == 0 {
		logs.Info("dal.updates[%T]: empty models, skip", zero)
		return
	}

	var columns []string
	updatableColumns := collext.Sets(exp.Zero[M]().UpdatableColumns())
	updateColumns := make(map[string]struct{})
	for _, v := range models {
		u := dal.U()
		hook(v, u)
		if u.IsEmpty() {
			logs.Info("dal.updates[%T]: non update items", zero)
			continue
		}
		if getId != nil && getId(v) <= 0 {
			err = errdef.ErrDalEntityIdInvalid
			return
		}
		if len(columns) != 0 {
			if !ccmp.ArraysEqual(columns, collext.MapKeys(u.ToMap())) {
				err = errdef.ErrDalUpdateColumnsUnmatched
				return
			}
			continue
		}

		for k := range u.ToMap() {
			if _, ok := updatableColumns[k]; !ok {
				err = errdef.ErrDalEntityUpdatableColumnInvalid
				return
			}
			if _, ok := updateColumns[k]; !ok {
				updateColumns[k] = struct{}{}
				columns = append(columns, k)
			}
		}
	}
	if len(columns) == 0 {
		logs.Info("dal.updates[%T]: empty update columns, skip", zero)
		return
	}

	clauses := clause.OnConflict{DoUpdates: clause.AssignmentColumns(columns)}
	conn := dal.W(ctx).Clauses(clauses).CreateInBatches(models, maxBatchSize)
	affected, err = conn.RowsAffected, conn.Error
	if err != nil {
		logs.Error("dal.updates[%T]: err = %v", zero, err)
	}
	return
}

func saveUpdatable[M cdb.TableModel, U cdb.Updater](ctx context.Context, dal txDal[U], models []*M) (affected int64, err error) {
	zero := exp.Zero[M]()
	if len(models) == 0 {
		logs.Info("dal.saveUpdatable[%T]: empty models, skip", zero)
		return 0, nil
	}

	updatable := zero.UpdatableColumns()
	columns := make([]string, 0, len(updatable))
	for _, v := range updatable {
		if _, ok := timeColumnsSets[v]; ok {
			continue
		}
		columns = append(columns, v)
	}
	if len(columns) == 0 {
		logs.Info("dal.saveUpdatable[%T]: empty update columns, skip", zero)
		return
	}

	clauses := clause.OnConflict{DoUpdates: clause.AssignmentColumns(columns)}
	conn := dal.W(ctx).Clauses(clauses).CreateInBatches(models, maxBatchSize)
	affected, err = conn.RowsAffected, conn.Error
	if err != nil {
		logs.Error("dal.saveUpdatable[%T]: err = %v", zero, err)
	}
	return
}
