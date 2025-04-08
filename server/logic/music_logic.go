package logic

import (
	"code.chenji.com/pkg/boost/collection/collext"
	"context"
	"gorm.io/gorm"
	"mime/multipart"
	"onij/biz/convert"
	"onij/biz/prm"
	"onij/handler/resq"
	"onij/infra"
	"onij/infra/mysql"
	"onij/model/api"
	"onij/util"
	"strings"
)

type MusicLogic interface {
	Upload(ctx context.Context, param *prm.UploadMusicParam) (*prm.UploadMusicResult, error)
	Save(music *mysql.Music, cover, mp, lyric, sheet *multipart.FileHeader) (int, error)
	DelById(id int) error

	ListByCond(req *resq.ListMusicReq) ([]*resq.ListMusicResp, error)
	GetMusic(id int) (*resq.GetMusicResp, error)
}

type musicLogic struct {
	*infra.AllInfra
}

func NewMusicLogic(i *infra.AllInfra) MusicLogic {
	return &musicLogic{
		AllInfra: i,
	}
}

func (l *musicLogic) Upload(ctx context.Context, param *prm.UploadMusicParam) (*prm.UploadMusicResult, error) {
	arts := collext.Pick(param.Artists, convert.SingerNameToArtist)
	if param.Composer != nil {
		arts = append(arts, convert.NameToArtist(*param.Composer, api.ArtistType_AT_Composer))
	}
	if param.Writer != nil {
		arts = append(arts, convert.NameToArtist(*param.Writer, api.ArtistType_AT_Writer))
	}
	err := l.ArtistDal.Save(arts...)
	if err != nil {
		return nil, err
	}
	err = l.MusicDal.Save(&mysql.Music{
		Id:          0,
		RootId:      0,
		Title:       "",
		ArtistIds:   "",
		Composer:    0,
		Writer:      0,
		Length:      0,
		IssueYear:   0,
		Language:    0,
		PerformType: 0,
		Concert:     "",
		ConcertYear: 0,
		Sequence:    0,
		MvUrl:       "",
		CoverOss:    0,
		MpOss:       0,
		LyricOss:    0,
		SheetOss:    0,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		DeletedAt:   gorm.DeletedAt{},
	})
	return &prm.UploadMusicResult{}, nil
}

func (m *musicLogic) DelById(id int) error {
	mus, err := m.MusicDal.DelById(id)
	if err != nil {
		return err
	}

	// del file
	_, err = m.FileDal.DelByIds([]int{mus.CoverOss, mus.MpOss, mus.LyricOss, mus.SheetOss})
	if err != nil {
		return err
	}
	return nil
}

func (m *musicLogic) GetByTitle(title string) ([]*mysql.Music, error) {
	return m.MusicDal.GetByTitle(title)
}

func (m *musicLogic) GetByArtist(artist int) ([]*mysql.Music, error) {
	return m.MusicDal.GetByArtist(artist)
}

func (m *musicLogic) GetMusic(id int) (*resq.GetMusicResp, error) {
	mu, err := m.MusicDal.GetById(id)
	if err != nil {
		return nil, err
	}

	singerNames, composer, writer, err := getNameFormMusic(mu)
	if err != nil {
		return nil, err
	}

	urls, err := m.FileDal.GetUrlByIds(mu.CoverOss, mu.MpOss, mu.LyricOss, mu.SheetOss)
	if err != nil {
		return nil, err
	}

	return &resq.GetMusicResp{
		Id:          mu.Id,
		RootId:      mu.RootId,
		Title:       mu.Title,
		ArtistIds:   mu.ArtistIds,
		Composer:    mu.Composer,
		Writer:      mu.Writer,
		IssueYear:   mu.IssueYear,
		Language:    mu.Language,
		PerformType: mu.PerformType,
		Concert:     mu.Concert,
		ConcertYear: mu.ConcertYear,
		Sequence:    mu.Sequence,
		MvUrl:       mu.MvUrl,
		CoverOss:    mu.CoverOss,
		MpOss:       mu.MpOss,
		LyricOss:    mu.LyricOss,
		SheetOss:    mu.SheetOss,

		ArtistName:   singerNames,
		ComposerName: composer,
		WriterName:   writer,
		CoverUrl:     urls[0],
		MpUrl:        urls[1],
		LyricUrl:     urls[2],
		SheetUrl:     urls[3],
	}, nil
}

func (m *musicLogic) ListByCond(req *resq.ListMusicReq) ([]*resq.ListMusicResp, error) {
	mus, err := m.MusicDal.GetByTitleArtistPerType(req.Title, req.Artist, req.PerformType, req.Page, req.Size)
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
	per, err := m.PerformerDal.GetByIds(append(util.DbToList(mu.ArtistIds), mu.Composer, mu.Writer)...)
	if err != nil {
		return "", "", "", err
	}
	singerNames := collext.Select(per, func(p *mysql.Artist) (string, bool) {
		return p.Name, p.PerformerType == 1 //TODO
	})
	singer := strings.Join(singerNames, "&")
	composer := collext.SelectOne(per, func(p *mysql.Artist) (string, bool) {
		return p.Name, p.PerformerType == 2
	})
	writer := collext.SelectOne(per, func(p *mysql.Artist) (string, bool) {
		return p.Name, p.PerformerType == 3
	})
	return singer, composer, writer, nil
}
