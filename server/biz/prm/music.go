package prm

import (
	"onij/biz/getter"
	"onij/infra/mysql"
	"onij/model/api"
	"onij/util"
	"onij/util/boost/collection/collext"
)

type UploadMusicParam struct {
	Id           *int64
	Name         string
	ArtistIds    []int64
	Mp3FileId    int64
	LyricsFileId int64
	ComposerId   *int64
	WriterId     *int64
	AlbumId      *int64
	MvUrl        *string
	RootMusicId  *int64
	IssueTime    *int32
}

func NewUploadMusicParam(req *api.UploadMusicReq) *UploadMusicParam {
	return &UploadMusicParam{
		Id:           req.MusicId,
		Name:         req.Name,
		ArtistIds:    req.ArtistIds,
		Mp3FileId:    req.Mp3FileId,
		LyricsFileId: req.LyricsFileId,
		ComposerId:   req.ComposerId,
		WriterId:     req.WriterId,
		AlbumId:      req.AlbumId,
		MvUrl:        req.MvUrl,
		RootMusicId:  req.RootMusicId,
		IssueTime:    req.IssueTime,
	}
}

type UploadMusicResult struct {
	MusicId int64
}

func (p *UploadMusicResult) Resp() *api.UploadMusicResp {
	return &api.UploadMusicResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		MusicId: p.MusicId,
	}
}

type GetMusicDetailParam struct {
	Id int64
}

func NewGetMusicDetailParam(req *api.GetMusicDetailReq) *GetMusicDetailParam {
	return &GetMusicDetailParam{
		Id: req.MusicId,
	}
}

type GetMusicDetailResult struct {
	Music      *mysql.Music
	Tags []*mysql.Tag
	ArtistMap  map[int64]*mysql.Artist
	ArtistIds  []int64
	ComposerId int64
	WriterId   int64
	Albums     []*mysql.Album
	FileMap    map[int64]*mysql.File
}

func (p *GetMusicDetailResult) Resp() *api.GetMusicDetailResp {
	album := &mysql.Album{}
	if len(p.Albums) > 0 {
		album = p.Albums[0]
	}
	var mp3Url, lyrUrl, abmUrl string
	if p.FileMap[p.Music.Mp3FileId] != nil {
		mp3Url = p.FileMap[p.Music.Mp3FileId].StoreKey
	}
	if p.FileMap[p.Music.LyricFileId] != nil {
		lyrUrl = p.FileMap[p.Music.LyricFileId].StoreKey
	}
	if p.FileMap[album.CoverFileId] != nil {
		abmUrl = p.FileMap[album.CoverFileId].StoreKey
	}
	return &api.GetMusicDetailResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Detail: &api.MusicDetail{
			Id:            p.Music.Id,
			Name:          p.Music.Name,
			ArtistIds:     collext.Pick(p.ArtistIds, func(id int64) int64 { return p.ArtistMap[id].Id }),
			ArtistNames:   collext.Pick(p.ArtistIds, func(id int64) string { return p.ArtistMap[id].Name }),
			ComposerId:    p.ComposerId,
			ComposerName:  p.ArtistMap[p.ComposerId].Name,
			WriterId:      p.WriterId,
			WriterName:    p.ArtistMap[p.WriterId].Name,
			IssueTime:     p.Music.IssueTime,
			MvUrl:         p.Music.MvUrl,
			Mp3FileUrl:    util.DownloadFile(mp3Url),
			LyricsFileUrl: util.DownloadFile(lyrUrl),
			AlbumProfile: &api.AlbumProfile{
				Id:           album.Id,
				Name:         album.Name,
				CoverFileUrl: util.DownloadFile(abmUrl),
			},
			TagDetails: collext.Pick(p.Tags, func(tag *mysql.Tag)*api.TagDetail {
				return &api.TagDetail{
					ResourceId:   tag.ResourceId,
					ResourceType: api.ResourceType(tag.ResourceType),
					TagBiz:       api.TagBiz(tag.TagBiz),
					TagGroup:     api.TagGroup(tag.TagGroup),
					TagType:      api.TagType(tag.TagType),
					TargetId:     &tag.TargetId,
					TargetType:   &tag.TargetType,
					Extra:        tag.Extra,
					ListShow:     tag.ListShow,
				}
			}),
		},
	}
}

type GetMusicListParam struct {
	SortType api.MusicSortType
	Keywords  []string
	AlbumId  *int64
	ArtistIds []int64
	WriterIds []int64
	ComposerIds []int64
	TagTypes []int32
	Page     int32
	Limit    int32
}

func NewGetMusicListParam(req *api.GetMusicListReq) *GetMusicListParam {
	return &GetMusicListParam{
		SortType:    req.SortType,
		Keywords:     req.Keywords,
		AlbumId:     req.AlbumId,
		ArtistIds:   req.ArtistIds,
		WriterIds:   req.WriterIds,
		ComposerIds: req.ComposerIds,
		TagTypes:     collext.Pick(req.TagTypes, func(tt api.TagType) int32{return int32(tt)}),
		Page:        req.Page,
		Limit:       req.Limit,
	}
}

type GetMusicListResult struct {
	Musics []*mysql.Music
	MusicArtistsMap map[int64][]*mysql.Artist
	MusicTagsMap map[int64][]*mysql.Tag
}

func (p *GetMusicListResult) Resp() *api.GetMusicListResp {
	return &api.GetMusicListResp{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Musics: collext.Pick(p.Musics, func(music *mysql.Music) *api.MusicProfile {
			return &api.MusicProfile{
				Id:             music.Id,
				Name:           music.Name,
				SingerProfiles: collext.Pick(p.MusicArtistsMap[music.Id], getter.ArtistToProfile),
				Tags:           collext.Pick(p.MusicTagsMap[music.Id], getter.TagToDetail),
			}
		}),
	}
}
