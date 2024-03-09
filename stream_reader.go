package anthropic

// See: https://docs.anthropic.com/claude/reference/messages-streaming

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

const (
	SSE_EVENT_PREFIX              = "event: "
	SSE_EVENT_MESSAGE_START       = "message_start"
	SSE_EVENT_CONTENT_BLOCK_START = "content_block_start"
	SSE_EVENT_PING                = "ping"
	SSE_EVENT_BLOCK_DELTA         = "content_block_delta"
	SSE_EVENT_CONTENT_BLOCK_STOP  = "content_block_stop"
	SSE_EVENT_MESSAGE_DELTA       = "message_delta"
	SSE_EVENT_MESSAGE_STOP        = "message_stop"
	SSE_EVENT_ERROR               = "error"

	SSE_DATA_PREFIX = "data: "
)

type MessagesStreamReader struct {
	reader   *bufio.Reader
	response *http.Response
}

func (s *MessagesStreamReader) Recv() (*MessagesStreamResponse, error) {
	var eventType string
	for {
		rawline, err := s.reader.ReadString('\n')
		if err != nil {
			return &MessagesStreamResponse{}, err
		}
		line := strings.TrimSpace(rawline)

		// event
		if strings.HasPrefix(line, SSE_EVENT_PREFIX) {
			eventType = strings.TrimPrefix(line, SSE_EVENT_PREFIX)
		}

		// handle data
		if strings.HasPrefix(line, SSE_DATA_PREFIX) {
			switch eventType {
			case SSE_EVENT_MESSAGE_START:
			case SSE_EVENT_CONTENT_BLOCK_START:
			case SSE_EVENT_PING:
			case SSE_EVENT_BLOCK_DELTA:
				var resp MessagesStreamResponse
				data := strings.TrimPrefix(line, SSE_DATA_PREFIX)
				if err = json.Unmarshal([]byte(data), &resp); err != nil {
					return &MessagesStreamResponse{}, errors.New("[Unmarshal MessagesStreamResponse] " + err.Error())
				}
				return &resp, nil
			case SSE_EVENT_CONTENT_BLOCK_STOP:
			case SSE_EVENT_MESSAGE_DELTA:
			case SSE_EVENT_MESSAGE_STOP:
				return &MessagesStreamResponse{}, io.EOF
			case SSE_EVENT_ERROR:
				var resp MessagesStreamResponse
				data := strings.TrimPrefix(line, SSE_DATA_PREFIX)
				if err = json.Unmarshal([]byte(data), &resp); err != nil {
					return &MessagesStreamResponse{}, errors.New("[Unmarshal MessagesStreamResponse] " + err.Error())
				}
				return &MessagesStreamResponse{}, errors.New("[Anthropic Error] " + resp.Error.Message)
			}
		}
	}
}

func (s *MessagesStreamReader) Close() {
	s.response.Body.Close()
}
