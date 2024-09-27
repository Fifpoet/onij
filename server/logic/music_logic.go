package logic

import (
	"mime/multipart"
	"onij/enum"
	"onij/infra/mysql"
)

type MusicLogic interface {
	Save(music *mysql.Music, cover, mp, lyric, sheet *multipart.FileHeader) (int, error)
	DelById(id int) error
}

type musicLogic struct {
}

func NewMusicLogic() MusicLogic {
	return &musicLogic{}
}

func (m *musicLogic) Save(music *mysql.Music, cover, mp, lyric, sheet *multipart.FileHeader) (int, error) {
	// 如果已有文件，删除
	if cover != nil {
		if music.CoverOss != 0 {
			_, err := app.FileDal.DelByIds([]int{music.CoverOss, music.MpOss, music.LyricOss, music.SheetOss})
			if err != nil {
				return 0, err
			}
		}
		music.CoverOss, _ = app.FileDal.CreateFileFromForm(cover, enum.BizMusic)
		music.MpOss, _ = app.FileDal.CreateFileFromForm(mp, enum.BizMusic)
		music.LyricOss, _ = app.FileDal.CreateFileFromForm(lyric, enum.BizMusic)
		music.SheetOss, _ = app.FileDal.CreateFileFromForm(sheet, enum.BizMusic)

	}

	return app.MusicDal.Save(music)
}

func (m *musicLogic) DelById(id int) error {
	mus, err := app.MusicDal.DelById(id)
	if err != nil {
		return err
	}

	// del file
	_, err = app.FileDal.DelByIds([]int{mus.CoverOss, mus.MpOss, mus.LyricOss, mus.SheetOss})
	if err != nil {
		return err
	}
	return nil
}
