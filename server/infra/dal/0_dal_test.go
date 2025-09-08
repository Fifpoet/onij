package dal

import (
	"onij/util"
	"testing"
)

func init() {
}

func TestSaveUpdateAll(t *testing.T) {
	db := NewMysqlCli()
	dal := NewArtistDal(db)

	err := dal.Save(&Artist{
		Id:         util.IdGen.Generate(),
		Name:       "test",
		ArtistType: 1,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestFirst(t *testing.T) {
	db := NewMysqlCli()
	dal := NewFileDal(db)
	res, err := dal.GetByParentAndHash(0, "aaa")
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}
