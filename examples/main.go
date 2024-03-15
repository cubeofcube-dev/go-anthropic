package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/wengchaoxi/go-anthropic"
)

func imageToBase64(filePath string) string {
	img, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	base64Str := base64.StdEncoding.EncodeToString(img)
	return base64Str
}

func main() {
	cli := anthropic.NewClient()

	req := anthropic.MessagesRequest{
		Model: anthropic.MODEL_CLAUDE_3_SONNET,
		Messages: []anthropic.Message{{
			Role: "user",
			Content: []anthropic.MessageContent{
				&anthropic.MessageContentFile{
					Type:   "image",
					Source: anthropic.MessageContentFileSource{Type: "base64", MediaType: "image/png", Data: imageToBase64("./image.png")},
				},
				&anthropic.MessageContentText{
					Type: "text",
					Text: "What is in this image?",
				},
			}},
		},
		MaxTokens: 1024,
		Stream:    true,
	}

	// resp, err := cli.CreateMessages(req)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println(resp.Content[0].Text)

	stream, err := cli.CreateMessagesStream(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stream.Close()
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err.Error())
			return
		}
		fmt.Print(resp.Delta.Text)
	}
	fmt.Println()
}
