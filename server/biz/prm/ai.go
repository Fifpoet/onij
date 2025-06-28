package prm

import "onij/model/api"


type AiChatParam struct {
	Prompt    string
	ImageURLs  []string
	VideoURLs  []string
	ThinkingType string
	SearchWeb bool
}

func NewAiChatParam(req *api.AiChatReq) *AiChatParam {
	tk := "disabled"
	switch req.Think {
	case api.ThinkingType_THT_Disabled:
		tk = "disabled"
	case api.ThinkingType_THT_Enabled:
		tk = "enabled"
	case api.ThinkingType_THT_Auto:
		tk = "auto"
	default:
		tk = "disabled"
	}
	return &AiChatParam{
		Prompt:    req.Prompt,
		ImageURLs:  req.ImageUrls,
		VideoURLs:  req.VideoUrls,
		ThinkingType: tk,
		SearchWeb: req.SearchWeb,
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
