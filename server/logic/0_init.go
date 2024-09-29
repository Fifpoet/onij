package logic

import (
	"errors"
	"onij/boost/collection/collext"
	"onij/enum"
	"onij/handler/resq"
	"onij/infra/mysql"
	"onij/inject"
)

var app *inject.App

func Init() {
	app = inject.InitializeApp()
}

type LocalLogic interface {
	SaveMusicFromDir(music []*mysql.Music, mps, lyrics []string) error
	SaveMeta(metas []*mysql.Meta) error

	GetMeta() (*resq.GetMetaResp, error)
}

type localLogic struct {
}

func NewLocalLogic() LocalLogic {
	return &localLogic{}
}

// SaveMusicFromDir 本地上传music
func (m *localLogic) SaveMusicFromDir(music []*mysql.Music, mps, lyrics []string) error {
	if len(music) != len(mps) || len(music) != len(lyrics) {
		return errors.New("music, mp, lyric length not match")
	}
	for i := 0; i < len(music); i++ {
		fid, err := app.FileDal.CreateLocalFile(mps[i], enum.BizMusic)
		if err != nil {
			return err
		}
		music[i].MpOss = fid
		if lyrics[i] != "" {
			fid, err = app.FileDal.CreateLocalFile(lyrics[i], enum.BizMusic)
			if err != nil {
				return err
			}
			music[i].LyricOss = fid
		}

		_, err = app.MusicDal.Save(music[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *localLogic) SaveMeta(metas []*mysql.Meta) error {
	err := app.MetaDal.Save(metas)
	return err
}

func (m *localLogic) GetMeta() (*resq.GetMetaResp, error) {
	metaCodes, err := app.MetaDal.GetByMetaEnumCode([]int{1})
	if err != nil {
		return nil, err
	}
	cods := collext.Pick(metaCodes, func(meta *mysql.Meta) int { return meta.Value })

	metas, err := app.MetaDal.GetByMetaEnumCode(cods)
	return &resq.GetMetaResp{
		MetaEnumCode: 1,
		MetaList: collext.Pick(metas, func(meta *mysql.Meta) resq.MetaModel {
			return resq.MetaModel{
				Value: meta.Value,
				Name:  meta.Name,
			}
		}),
	}, err
}
