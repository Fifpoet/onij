package model

/*********** 表字段常量定义 **********/

import (
	"gorm.io/gorm"
	"onij/util/cdb"
	"time"
)

const (
	TableName_Music = "music"

	Music_Id           = "id"
	Music_Name         = "name"
	Music_FullName     = "full_name"
	Music_ArtistIds    = "artist_ids"
	Music_ComposerIds  = "composer_ids"
	Music_WriterIds    = "writer_ids"
	Music_IssueTime    = "issue_time"
	Music_PerformType  = "perform_type"
	Music_TimeLength   = "time_length"
	Music_MvUrl        = "mv_url"
	Music_AudioFileId  = "audio_file_id"
	Music_AudioQuality = "audio_quality"
	Music_LyricFileId  = "lyric_file_id"
	Music_LyricContent = "lyric_content"
	Music_RootId       = "root_id"
	Music_Priority     = "priority"
	Music_CreatedAt    = "created_at"
	Music_UpdatedAt    = "updated_at"
	Music_DeletedAt    = "deleted_at"
)

/*********** 表结构定义 **********/

// Music 表 music 结构定义
type Music struct {
	Id           int64          `gorm:"column:id" json:"id"`                       // 主键id
	Name         string         `gorm:"column:name" json:"name"`                   // 名称
	FullName     string         `gorm:"column:full_name" json:"full_name"`         // 全名称
	ArtistIds    string         `gorm:"column:artist_ids" json:"artist_ids"`       // 艺术家id列表
	ComposerIds  string         `gorm:"column:composer_ids" json:"composer_ids"`   // 作曲家id列表
	WriterIds    string         `gorm:"column:writer_ids" json:"writer_ids"`       // 作词家id列表
	IssueTime    time.Time      `gorm:"column:issue_time" json:"issue_time"`       // 发行时间
	PerformType  int32          `gorm:"column:perform_type" json:"perform_type"`   // 表演类型
	TimeLength   int32          `gorm:"column:time_length" json:"time_length"`     // 时长(s)
	MvUrl        string         `gorm:"column:mv_url" json:"mv_url"`               // mv链接
	AudioFileId  int64          `gorm:"column:audio_file_id" json:"audio_file_id"` // 音频文件id
	AudioQuality int32          `gorm:"column:audio_quality" json:"audio_quality"` // 音频质量类型
	LyricFileId  int64          `gorm:"column:lyric_file_id" json:"lyric_file_id"` // 歌词文件id
	LyricContent string         `gorm:"column:lyric_content" json:"lyric_content"` // 歌词内容
	RootId       int64          `gorm:"column:root_id" json:"root_id"`             // 根id
	Priority     int32          `gorm:"column:priority" json:"priority"`           // 优先级
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`       // 创建时间
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`       // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`       // 删除时间
}

func NewMusic() *Music {
	return &Music{}
}

func (Music) TableName() string {
	return TableName_Music
}

func (Music) UpdatableColumns() []string {
	return []string{
		Music_Name,
		Music_FullName,
		Music_ArtistIds,
		Music_ComposerIds,
		Music_WriterIds,
		Music_IssueTime,
		Music_PerformType,
		Music_TimeLength,
		Music_MvUrl,
		Music_AudioFileId,
		Music_AudioQuality,
		Music_LyricFileId,
		Music_LyricContent,
		Music_RootId,
		Music_Priority,
		Music_CreatedAt,
		Music_UpdatedAt,
		Music_DeletedAt,
	}
}

/*********** 查询器与更新器 **********/

// MusicQuerier 表 music 查询器
type MusicQuerier cdb.BasicQuerier

func (q MusicQuerier) super() cdb.BasicQuerier {
	return (cdb.BasicQuerier)(q)
}

func (q MusicQuerier) ToOptions() []cdb.Option {
	return q.super().ToOptions()
}

func (q MusicQuerier) WithAsc(cols ...string) MusicQuerier {
	q.super().WithAsc(cols...)
	return q
}

func (q MusicQuerier) WithDesc(cols ...string) MusicQuerier {
	q.super().WithDesc(cols...)
	return q
}

func (q MusicQuerier) WithLimit(limit int) MusicQuerier {
	q.super().WithLimit(limit)
	return q
}

func (q MusicQuerier) WithUnscoped() MusicQuerier {
	q.super().WithUnscoped()
	return q
}

func (q MusicQuerier) Id(v any) MusicQuerier {
	q.super().Add(Music_Id, v)
	return q
}

func (q MusicQuerier) Name(v any) MusicQuerier {
	q.super().Add(Music_Name, v)
	return q
}

func (q MusicQuerier) FullName(v any) MusicQuerier {
	q.super().Add(Music_FullName, v)
	return q
}

func (q MusicQuerier) ArtistIds(v any) MusicQuerier {
	q.super().Add(Music_ArtistIds, v)
	return q
}

func (q MusicQuerier) ComposerIds(v any) MusicQuerier {
	q.super().Add(Music_ComposerIds, v)
	return q
}

func (q MusicQuerier) WriterIds(v any) MusicQuerier {
	q.super().Add(Music_WriterIds, v)
	return q
}

func (q MusicQuerier) IssueTime(v any) MusicQuerier {
	q.super().Add(Music_IssueTime, v)
	return q
}

func (q MusicQuerier) PerformType(v any) MusicQuerier {
	q.super().Add(Music_PerformType, v)
	return q
}

func (q MusicQuerier) TimeLength(v any) MusicQuerier {
	q.super().Add(Music_TimeLength, v)
	return q
}

func (q MusicQuerier) MvUrl(v any) MusicQuerier {
	q.super().Add(Music_MvUrl, v)
	return q
}

func (q MusicQuerier) AudioFileId(v any) MusicQuerier {
	q.super().Add(Music_AudioFileId, v)
	return q
}

func (q MusicQuerier) AudioQuality(v any) MusicQuerier {
	q.super().Add(Music_AudioQuality, v)
	return q
}

func (q MusicQuerier) LyricFileId(v any) MusicQuerier {
	q.super().Add(Music_LyricFileId, v)
	return q
}

func (q MusicQuerier) LyricContent(v any) MusicQuerier {
	q.super().Add(Music_LyricContent, v)
	return q
}

func (q MusicQuerier) RootId(v any) MusicQuerier {
	q.super().Add(Music_RootId, v)
	return q
}

func (q MusicQuerier) Priority(v any) MusicQuerier {
	q.super().Add(Music_Priority, v)
	return q
}

func (q MusicQuerier) CreatedAt(v any) MusicQuerier {
	q.super().Add(Music_CreatedAt, v)
	return q
}

func (q MusicQuerier) UpdatedAt(v any) MusicQuerier {
	q.super().Add(Music_UpdatedAt, v)
	return q
}

func (q MusicQuerier) DeletedAt(v any) MusicQuerier {
	q.super().Add(Music_DeletedAt, v)
	return q
}

// MusicUpdater 表 music 更新器
type MusicUpdater cdb.BasicUpdater

func (u MusicUpdater) super() cdb.BasicUpdater {
	return (cdb.BasicUpdater)(u)
}

func (u MusicUpdater) Id(v int64) MusicUpdater {
	u.super().Add(Music_Id, v)
	return u
}

func (u MusicUpdater) Name(v string) MusicUpdater {
	u.super().Add(Music_Name, v)
	return u
}

func (u MusicUpdater) FullName(v string) MusicUpdater {
	u.super().Add(Music_FullName, v)
	return u
}

func (u MusicUpdater) ArtistIds(v string) MusicUpdater {
	u.super().Add(Music_ArtistIds, v)
	return u
}

func (u MusicUpdater) ComposerIds(v string) MusicUpdater {
	u.super().Add(Music_ComposerIds, v)
	return u
}

func (u MusicUpdater) WriterIds(v string) MusicUpdater {
	u.super().Add(Music_WriterIds, v)
	return u
}

func (u MusicUpdater) IssueTime(v time.Time) MusicUpdater {
	u.super().Add(Music_IssueTime, v)
	return u
}

func (u MusicUpdater) PerformType(v int32) MusicUpdater {
	u.super().Add(Music_PerformType, v)
	return u
}

func (u MusicUpdater) TimeLength(v int32) MusicUpdater {
	u.super().Add(Music_TimeLength, v)
	return u
}

func (u MusicUpdater) MvUrl(v string) MusicUpdater {
	u.super().Add(Music_MvUrl, v)
	return u
}

func (u MusicUpdater) AudioFileId(v int64) MusicUpdater {
	u.super().Add(Music_AudioFileId, v)
	return u
}

func (u MusicUpdater) AudioQuality(v int32) MusicUpdater {
	u.super().Add(Music_AudioQuality, v)
	return u
}

func (u MusicUpdater) LyricFileId(v int64) MusicUpdater {
	u.super().Add(Music_LyricFileId, v)
	return u
}

func (u MusicUpdater) LyricContent(v string) MusicUpdater {
	u.super().Add(Music_LyricContent, v)
	return u
}

func (u MusicUpdater) RootId(v int64) MusicUpdater {
	u.super().Add(Music_RootId, v)
	return u
}

func (u MusicUpdater) Priority(v int32) MusicUpdater {
	u.super().Add(Music_Priority, v)
	return u
}

func (u MusicUpdater) CreatedAt(v time.Time) MusicUpdater {
	u.super().Add(Music_CreatedAt, v)
	return u
}

func (u MusicUpdater) UpdatedAt(v time.Time) MusicUpdater {
	u.super().Add(Music_UpdatedAt, v)
	return u
}

func (u MusicUpdater) DeletedAt(v time.Time) MusicUpdater {
	u.super().Add(Music_DeletedAt, v)
	return u
}

func (u MusicUpdater) ToMap() map[string]any {
	return u.super().ToMap()
}

func (u MusicUpdater) IsEmpty() bool {
	return u.super().IsEmpty()
}
