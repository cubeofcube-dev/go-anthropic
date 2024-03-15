# Go Anthropic

**English** | [简体中文](./README.zh-CN.md)

Anthropic SDK implemented in Go, supporting models such as Claude 2.1, Claude 3 (supports sending images), etc.

## Installation

```
go get github.com/wengchaoxi/go-anthropic
```

## Usage

By default, it will fetch `ANTHROPIC_API_KEY` and `ANTHROPIC_BASE_URL` from the environment variables.

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
		// `MODEL_CLAUDE_3_SONNET`、`MODEL_CLAUDE_3_OPUS`、`MODEL_CLAUDE_2_1`
		Model: anthropic.MODEL_CLAUDE_3_HAIKU
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

## Streaming Responses

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
		// `MODEL_CLAUDE_3_SONNET`、`MODEL_CLAUDE_3_OPUS`、`MODEL_CLAUDE_2_1`
		Model: anthropic.MODEL_CLAUDE_3_HAIKU
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

## References

- https://docs.anthropic.com/claude/reference/messages_post
- https://docs.anthropic.com/claude/reference/messages-streaming
