package logic

type MusicLogic interface {
}

type musicLogic struct {
}

func NewMusicLogic() MusicLogic {
	return &musicLogic{}
}
