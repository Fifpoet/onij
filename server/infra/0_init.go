package infra

import (
	"onij/infra/dal"
)

// AllInfra 定义上层需要的Dal对象
type AllInfra struct {
	dal.TagDal
	dal.FileDal
	dal.MusicDal
	dal.AlbumDal
	dal.AlbumMusicDal
	dal.ArtistDal
	dal.PracticeDal
	dal.TranItemDal
	dal.AiSessionDal
	dal.AiMessageDal
	dal.AiConfirmPendingDal
	dal.MusicCollectionDal
	dal.MusicCollectionItemDal
	dal.FileVideoMarkerDal
}
