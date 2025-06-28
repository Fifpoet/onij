package logic

import (
	"context"
	"errors"
	"fmt"
	"io"
	"onij/biz/prm"
	"onij/infra"
	"onij/util/boost/collection/collext"
	"os"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/utils"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
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
	client := arkruntime.NewClientWithApiKey(os.Getenv("ark"))

	msgs := []*model.ChatCompletionMessageContentPart{}
	if param.Prompt != "" {
		msgs = append(msgs, &model.ChatCompletionMessageContentPart{
			Type: model.ChatCompletionMessageContentPartTypeText,
			Text: param.Prompt,
		})
	}
	msgs = append(msgs, collext.Pick(param.ImageURLs, func(url string) *model.ChatCompletionMessageContentPart {
		return &model.ChatCompletionMessageContentPart{
			Type: model.ChatCompletionMessageContentPartTypeImageURL,
			ImageURL: &model.ChatMessageImageURL{
				URL: url,
				Detail: model.ImageURLDetailAuto, // 图片处理质量
			},
		}
	})...)
	msgs = append(msgs, collext.Pick(param.VideoURLs, func(url string) *model.ChatCompletionMessageContentPart {
		return &model.ChatCompletionMessageContentPart{
			Type: model.ChatCompletionMessageContentPartTypeVideoURL,
			VideoURL: &model.ChatMessageVideoURL{
				URL: url,
			},
		}
	})...)

	modelMsgs := []*model.ChatCompletionMessage{
		{
			Role: model.ChatMessageRoleSystem,
			Content: &model.ChatCompletionMessageContent{
				StringValue: volcengine.String("你是豆包，是由字节跳动开发的 AI 人工智能助手"),
			},
		},
		{
			Role: model.ChatMessageRoleUser,
			Content: &model.ChatCompletionMessageContent{
				ListValue: msgs,
			},
		},
	}

	var stream *utils.ChatCompletionStreamReader
	var err error
	if param.SearchWeb {
		client.CreateBotChatCompletionStream(ctx, model.BotChatCompletionRequest{
			Model: "doubao-seed-1-6-250615",
			Messages: modelMsgs,
			Stream: true,
			
		})
	} else {
		stream, err = client.CreateChatCompletionStream(ctx, model.CreateChatCompletionRequest{
			Model: "doubao-seed-1-6-250615",
			Messages: modelMsgs,
			Stream: volcengine.Bool(true),
			Thinking: &model.Thinking{	
				Type: model.ThinkingType(param.ThinkingType),
			},
		})
	}
	if err != nil {
		fmt.Printf("stream chat error: %v\n", err)
		return nil, err
	}
	
	defer stream.Close()

	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			return nil, errors.New("stream closed")
		}
		if err != nil {
			fmt.Printf("Stream chat error: %v\n", err)
			return nil, err
		}
		if len(recv.Choices) > 0 {
			fmt.Print(recv.Choices[0].Delta.Content)
		}
	}
}
