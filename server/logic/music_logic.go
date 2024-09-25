package logic

import (
	"mime/multipart"
	"onij/enum"
	"onij/infra/mysql"
)

type MusicLogic interface {
	Save(music *mysql.Music, cover, mp, lyric, sheet *multipart.FileHeader) (int, error)
}

type musicLogic struct {
}

func NewMusicLogic() MusicLogic {
	return &musicLogic{}
}

func (m *musicLogic) Save(music *mysql.Music, cover, mp, lyric, sheet *multipart.FileHeader) (int, error) {
	if cover != nil {
		if music.CoverOss != 0 {
			err := app.FileDal.DelById(music.CoverOss)
			if err != nil {
				return 0, err
			}
		}
		music.CoverOss, _ = app.FileDal.CreateFileFromForm(cover, enum.BizMusic)
	}
	if mp != nil {
		if music.MpOss != 0 {
			err := app.FileDal.DelById(music.MpOss)
			if err != nil {
				return 0, err
			}
		}
		music.MpOss, _ = app.FileDal.CreateFileFromForm(mp, enum.BizMusic)
	}
	if lyric != nil {
		if music.LyricOss != 0 {
			err := app.FileDal.DelById(music.LyricOss)
			if err != nil {
				return 0, err
			}
		}
		music.LyricOss, _ = app.FileDal.CreateFileFromForm(lyric, enum.BizMusic)
	}
	if sheet != nil {
		if music.SheetOss != 0 {
			err := app.FileDal.DelById(music.SheetOss)
			if err != nil {
				return 0, err
			}
		}
		music.SheetOss, _ = app.FileDal.CreateFileFromForm(sheet, enum.BizMusic)
	}

	return app.MusicDal.Save(music)
}
