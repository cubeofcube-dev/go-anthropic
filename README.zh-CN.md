# Go Anthropic

**简体中文** | [English](./README.md)

用 Go 实现的 Anthropic SDK，支持：Claude 2.1、Claude 3（支持发送图片）等模型。

## 安装

```
go get github.com/wengchaoxi/go-anthropic
```

## 简单使用

> 默认会从环境变量中获取 `ANTHROPIC_API_KEY`、`ANTHROPIC_BASE_URL`。

```go
package main

import (
	"fmt"

	"github.com/wengchaoxi/go-anthropic"
)

func main() {
	cli := anthropic.NewClient()
	// cli := anthropic.NewClient(anthropic.ClientOptions{
	// 	ApiKey:  os.Getenv("ANTHROPIC_API_KEY"),
	// 	BaseUrl: anthropic.DEFAULT_ANTHROPIC_BASE_URL,
	// })

	resp, err := cli.CreateMessages(anthropic.MessagesRequest{
		Model: anthropic.MODEL_CLAUDE_3_SONNET, // 或者 `MODEL_CLAUDE_3_OPUS`、`MODEL_CLAUDE_2_1`
		Messages: []anthropic.Message{{
			Role: "user",
			Content: []anthropic.MessageContent{
				&anthropic.MessageContentText{
					Type: "text",
					Text: "Hello Claude!",
				},
			}},
		},
		MaxTokens: 1024,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp.Content[0].Text)
}
```

## 流式响应

```go
package main

import (
	"fmt"

	"github.com/wengchaoxi/go-anthropic"
)

func main() {
	cli := anthropic.NewClient()
	// cli := anthropic.NewClient(anthropic.ClientOptions{
	// 	ApiKey:  os.Getenv("ANTHROPIC_API_KEY"),
	// 	BaseUrl: anthropic.DEFAULT_ANTHROPIC_BASE_URL,
	// })

	stream, _ := cli.CreateMessagesStream(anthropic.MessagesRequest{
		Model: anthropic.MODEL_CLAUDE_3_SONNET, // 或者 `MODEL_CLAUDE_3_OPUS`、`MODEL_CLAUDE_2_1`
		Messages: []anthropic.Message{{
			Role: "user",
			Content: []anthropic.MessageContent{
				&anthropic.MessageContentText{
					Type: "text",
					Text: "Hello Claude!",
				},
			}},
		},
		MaxTokens: 1024,
		Stream:    true,
	})
	defer stream.Close()
	for {
		resp, err := stream.Recv()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(resp.Delta.Text)
	}
}
```

## 相关文档

- https://docs.anthropic.com/claude/reference/messages_post
- https://docs.anthropic.com/claude/reference/messages-streaming
