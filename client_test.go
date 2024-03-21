package anthropic_test

import (
	"io"
	"testing"

	"github.com/cubeofcube-dev/go-anthropic"
)

func TestCreateMessages(t *testing.T) {
	cli := anthropic.NewClient()

	model := anthropic.MODEL_CLAUDE_3_SONNET
	content := anthropic.MessageContentText{Type: "text", Text: "Ping"}
	max_tokens := 9

	t.Run("Messages", func(t *testing.T) {
		messages := anthropic.MessagesRequest{
			Model: model,
			Messages: []anthropic.Message{
				{Role: "user", Content: []anthropic.MessageContent{&content}},
			},
			MaxTokens: max_tokens,
		}
		_, err := cli.CreateMessages(messages)
		if err != nil {
			t.Errorf("CreateMessages() error = %v", err)
			return
		}
	})
	t.Run("Messages Stream", func(t *testing.T) {
		messages := anthropic.MessagesRequest{
			Model: model,
			Messages: []anthropic.Message{
				{Role: "user", Content: []anthropic.MessageContent{&content}},
			},
			MaxTokens: max_tokens,
			Stream:    true,
		}
		stream, err := cli.CreateMessagesStream(messages)
		if err != nil {
			t.Errorf("CreateMessagesStream() error = %v", err)
			return
		}
		defer stream.Close()
		for {
			_, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Errorf("Recv() error = %v", err)
				return
			}
		}
	})
}
