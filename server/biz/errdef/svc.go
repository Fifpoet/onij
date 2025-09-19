package errdef

import "errors"

// album

var (
	ErrAlbumExisted = errors.New("专辑已存在")
)

// file

var (
	ErrUrlFileExtUnknown = errors.New("未知链接文件后缀")
	ErrFolderIdNotExist  = errors.New("文件夹不存在")
)

// music

var (
	ErrResponseNoLyric      = errors.New("响应中无歌词")
	ErrThirdMusicIdNotFound = errors.New("第三方音乐不存在")
)
