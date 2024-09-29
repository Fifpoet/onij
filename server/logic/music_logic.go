package logic

import (
	"mime/multipart"
	"onij/boost/collection/collext"
	"onij/enum"
	"onij/handler/resq"
	"onij/infra/mysql"
	"onij/util"
	"strings"
)

type MusicLogic interface {
	Save(music *mysql.Music, cover, mp, lyric, sheet *multipart.FileHeader) (int, error)
	DelById(id int) error

	ListByCond(req *resq.ListMusicReq) ([]*resq.ListMusicResp, error)
	GetMusic(id int) (*resq.GetMusicResp, error)
}

type musicLogic struct {
}

func NewMusicLogic() MusicLogic {
	return &musicLogic{}
}

// Save .
// 这里如果原始文件已存在, 则fileDal校验hash后返回原id
func (m *musicLogic) Save(music *mysql.Music, cover, mp, lyric, sheet *multipart.FileHeader) (int, error) {
	// 如果上传文件成功, 则为覆盖场景, 更新model
	var err error
	cov, err := app.FileDal.CreateFormFile(cover, enum.BizMusic)
	if err != nil {
		return 0, err
	}
	if cov != 0 {
		music.CoverOss = cov
	}
	mpo, err := app.FileDal.CreateFormFile(mp, enum.BizMusic)
	if err != nil {
		return 0, err
	}
	if mpo != 0 {
		music.MpOss = mpo
	}
	lrc, err := app.FileDal.CreateFormFile(lyric, enum.BizMusic)
	if err != nil {
		return 0, err
	}
	if lrc != 0 {
		music.LyricOss = lrc
	}
	sht, err := app.FileDal.CreateFormFile(sheet, enum.BizMusic)
	if err != nil {
		return 0, err
	}
	if sht != 0 {
		music.SheetOss = sht
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

func (m *musicLogic) GetByTitle(title string) ([]*mysql.Music, error) {
	return app.MusicDal.GetByTitle(title)
}

func (m *musicLogic) GetByArtist(artist int) ([]*mysql.Music, error) {
	return app.MusicDal.GetByArtist(artist)
}

func (m *musicLogic) GetMusic(id int) (*resq.GetMusicResp, error) {
	mu, err := app.MusicDal.GetById(id)
	if err != nil {
		return nil, err
	}

	singerNames, composer, writer, err := getNameFormMusic(mu)
	if err != nil {
		return nil, err
	}

	urls, err := app.FileDal.GetUrlByIds(mu.CoverOss, mu.MpOss, mu.LyricOss, mu.SheetOss)
	if err != nil {
		return nil, err
	}

	return &resq.GetMusicResp{
		Title:    mu.Title,
		Artist:   singerNames,
		Composer: composer,
		Writer:   writer,
		Concert:  mu.Concert,
		MvUrl:    mu.MvUrl,
		CoverUrl: urls[0],
		MpUrl:    urls[1],
		LyricUrl: urls[2],
		SheetUrl: urls[3],
	}, nil
}

func (m *musicLogic) ListByCond(req *resq.ListMusicReq) ([]*resq.ListMusicResp, error) {
	mus, err := app.MusicDal.GetByTitleArtistPerType(req.Title, req.Artist, req.PerformType, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	res := collext.Select(mus, func(mu *mysql.Music) (*resq.ListMusicResp, bool) {
		singerNames, composer, writer, err := getNameFormMusic(mu)
		if err != nil {
			return nil, false
		}
		return &resq.ListMusicResp{
			Id:       mu.Id,
			Title:    mu.Title,
			Artist:   singerNames,
			Composer: composer,
			Writer:   writer,
			Concert:  mu.Concert,
			Sequence: mu.Sequence,
			MvUrl:    mu.MvUrl,
		}, true
	})
	return res, nil
}

func getNameFormMusic(mu *mysql.Music) (string, string, string, error) {
	per, err := app.PerformerDal.GetByIds(append(util.DbToList(mu.ArtistIds), mu.Composer, mu.Writer)...)
	if err != nil {
		return "", "", "", err
	}
	singerNames := collext.Select(per, func(p *mysql.Performer) (string, bool) {
		return p.Name, p.PerformerType == 1 //TODO
	})
	singer := strings.Join(singerNames, "&")
	composer := collext.SelectOne(per, func(p *mysql.Performer) (string, bool) {
		return p.Name, p.PerformerType == 2
	})
	writer := collext.SelectOne(per, func(p *mysql.Performer) (string, bool) {
		return p.Name, p.PerformerType == 3
	})
	return singer, composer, writer, nil
}
