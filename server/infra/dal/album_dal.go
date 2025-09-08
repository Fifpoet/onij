package dal

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"onij/model"
	"onij/util/cdb"
)

type AlbumDal interface {
	cdb.Interface[AlbumDal]

	Save(albums ...*model.Album) error
	GetById(id int64) (*model.Album, error)
	GetByRelatedId(id int64) ([]*model.Album, error)
	GetByMusicId(id int64) ([]*model.Album, error)
	GetByArtistId(id int64, offset, limit int32) ([]*model.Album, int32, error)
}
type albumDal struct {
	*cdb.Dal[ model.Album, model.AlbumQuerier,  model.AlbumUpdater]
}

func NewAlbumDal(db *cdb.DefaultProxy) AlbumDal {
	return &albumDal{
		cdb.NewDal[model.Album, model.AlbumQuerier, model.AlbumUpdater](db),
	}
}

func (r *albumDal) With(tx *gorm.DB) AlbumDal {
	return &albumDal{r.Dal.With(tx)}
}

func (f *albumDal) Save(albums ...*model.Album) error {
	if err := f.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&album).Error; err != nil {
		log.Printf("save album err: %v", err)
		return err
	}
	return nil
}

func (f *albumDal) GetById(id int64) (*model.Album, error) {
	res, err :=
}
func (f *albumDal) GetByRelatedId(id int64) ([]*model.Album, error) {
	var models []*model.Album
	if err := f.db.Where("related_album_id = ?", id).Find(&models).Error; err != nil {
		log.Printf("get album by music id err: %v", err)
		return nil, err
	}
	return models, nil

}
func (f *albumDal) GetByMusicId(id int64) ([]*model.Album, error) {
	var models []*model.Album
	if err := f.db.Where("music_id = ?", id).Find(&models).Error; err != nil {
		log.Printf("get album by music id err: %v", err)
		return nil, err
	}
	return models, nil
}
func (f *albumDal) GetByArtistId(id int64, offset, limit int32) ([]*model.Album, int32, error) {
	var models []*model.Album
	var total int64
	if err := f.db.Where("artist_id = ?", id).Model(&Album{}).Count(&total).Error; err != nil {
		log.Printf("count album by artist id err: %v", err)
		return nil, 0, err
	}
	if err := f.db.Where("artist_id = ?", id).Offset(int(offset)).Limit(int(limit)).Find(&models).Error; err != nil {
		log.Printf("get album by artist id err: %v", err)
		return nil, 0, err
	}
	return models, int32(total), nil
}
