package errdef

import "errors"

var (
	ErrAlbumExisted = errors.New("专辑已存在")
)

var (
	ErrUrlFileExtUnknown = errors.New("未知链接文件后缀")
	ErrFolderIdNotExist  = errors.New("文件夹不存在")
)
