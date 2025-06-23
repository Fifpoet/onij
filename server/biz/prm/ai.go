package prm

import "onij/model/api"


type AiChatParam struct {
	Prompt    string
	ImageURL  string
}

func NewAiChatParam(req *api.AiChatReq) *AiChatParam {
	return &AiChatParam{
		Prompt:    req.Prompt,
		ImageURL:  req.ImageUrl,
	}
}

type AiChatResult struct {
	Code    int32
	Message string
	Answer  string
}

func (r *AiChatResult) Resp() *api.AiChatResp {
	return &api.AiChatResp{
		Code:    r.Code,
		Message: r.Message,
		Answer:  r.Answer,
	}
}
