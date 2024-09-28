package logic

import (
	"errors"
	"onij/enum"
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

	GetMeta(metaEnumCode int) ([]*mysql.Meta, error)
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
		fid, err := app.FileDal.CreateFileFormLocal(mps[i], enum.BizMusic)
		if err != nil {
			return err
		}
		music[i].MpOss = fid
		if lyrics[i] != "" {
			fid, err = app.FileDal.CreateFileFormLocal(lyrics[i], enum.BizMusic)
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

func (m *localLogic) GetMeta(metaEnumCode int) ([]*mysql.Meta, error) {
	return app.MetaDal.GetByMetaEnumCode(metaEnumCode)
}
