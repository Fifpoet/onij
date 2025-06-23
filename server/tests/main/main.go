package main

import (
    "context"
    "fmt"
    "os"

    ark "github.com/sashabaranov/go-openai"
)

func main() {
    config := ark.DefaultConfig(os.Getenv("ARK_API_KEY"))
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
                      URL: "https://ark-project.tos-cn-beijing.ivolces.com/images/view.jpeg",
                   },
                },
                {

                   Type: ark.ChatMessagePartTypeText,
                   Text: "这是哪里？",
                },
             },
          },
       },
    }

    resp, err := client.CreateChatCompletion(context.Background(), req)
    if err != nil {
       fmt.Printf("ChatCompletion error: %v", err)
       return
    }
    fmt.Println(resp.Choices[0].Message.Content)
}