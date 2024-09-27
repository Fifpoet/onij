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

// Save .
// 这里如果原始文件已存在, 则fileDal校验hash后返回原id
func (m *musicLogic) Save(music *mysql.Music, cover, mp, lyric, sheet *multipart.FileHeader) (int, error) {
	// 处理文件
	var err error
	cov, err := app.FileDal.CreateFileFromForm(cover, enum.BizMusic)
	if err != nil {
		return 0, err
	}
	music.CoverOss = cov
	mpo, err := app.FileDal.CreateFileFromForm(mp, enum.BizMusic)
	if err != nil {
		return 0, err
	}
	music.MpOss = mpo
	lrc, err := app.FileDal.CreateFileFromForm(lyric, enum.BizMusic)
	if err != nil {
		return 0, err
	}
	music.LyricOss = lrc
	sht, err := app.FileDal.CreateFileFromForm(sheet, enum.BizMusic)
	if err != nil {
		return 0, err
	}
	music.SheetOss = sht

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
