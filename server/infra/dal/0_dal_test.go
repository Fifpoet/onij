package dal

import (
	"context"
	"onij/model"
	"onij/util"
	"onij/util/cdb"
	"testing"
)

func init() {
}

func TestSaveUpdateAll(t *testing.T) {
	db := NewMysqlCli()
	proxy := cdb.NewDefaultProxy(db)
	dal := NewArtistDal(proxy)

	ctx := context.Background()
	_, err := dal.Save(ctx, &model.Artist{
		Id:   util.IdGen.Generate(),
		Name: "test",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestFirst(t *testing.T) {
	db := NewMysqlCli()
	proxy := cdb.NewDefaultProxy(db)
	dal := NewFileDal(proxy)

	ctx := context.Background()
	res, err := dal.GetByParentAndHash(ctx, 0, "aaa")
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}
