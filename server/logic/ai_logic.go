package logic

import (
	"context"
	"fmt"
	"onij/biz/prm"
	"onij/infra"
	"onij/util"
	"os"

	ark "github.com/sashabaranov/go-openai"
)

type AiLogic interface {
	Chat(ctx context.Context, param *prm.AiChatParam) (*prm.AiChatResult, error)
}

type aiLogic struct {
	*infra.AllInfra
}

func NewAiLogic(i *infra.AllInfra) AiLogic {
	return &aiLogic{
		AllInfra: i,
	}
}

func (l *aiLogic) Chat(ctx context.Context, param *prm.AiChatParam) (*prm.AiChatResult, error) {
	config := ark.DefaultConfig(os.Getenv("ark"))
	config.BaseURL = "https://ark.cn-beijing.volces.com/api/v3"
	client := ark.NewClientWithConfig(config)

	fmt.Println("----- image input -----")
	req := ark.ChatCompletionRequest{
		Model: "doubao-seed-1-6-250615",
		Messages: []ark.ChatCompletionMessage{
			{
				Role: ark.ChatMessageRoleUser,
				MultiContent: []ark.ChatMessagePart{
					{
						Type: ark.ChatMessagePartTypeImageURL,
						ImageURL: &ark.ChatMessageImageURL{
							URL: param.ImageURL,
						},
					},
					{

						Type: ark.ChatMessagePartTypeText,
						Text: param.Prompt,
					},
				},
			},
		},
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v", err)
		return nil, err
	}
	fmt.Println(resp.Choices[0].Message.Content)

	return &prm.AiChatResult{
		Code:    util.BaseCodeOK,
		Message: util.BaseMsgOK,
		Answer:  resp.Choices[0].Message.Content,
	}, nil
}
