package logic

import (
	"mime/multipart"
	"onij/infra/mysql"
)

type MusicLogic interface {
	Save(music *mysql.Music, file *multipart.FileHeader) (int, error)
}

type musicLogic struct {
}

func (m *musicLogic) Save(music *mysql.Music, file *multipart.FileHeader) (int, error) {
	//TODO implement me
	panic("implement me")
}

func NewMusicLogic() MusicLogic {
	return &musicLogic{}
}
