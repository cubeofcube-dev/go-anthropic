package anthropic

// See: https://docs.anthropic.com/claude/reference/messages_post

const (
	// See: https://docs.anthropic.com/claude/docs/models-overview#model-recommendations
	MODEL_CLAUDE_2_1      = "claude-2.1"
	MODEL_CLAUDE_3_SONNET = "claude-3-sonnet-20240229"
	MODEL_CLAUDE_3_OPUS   = "claude-3-opus-20240229"
	// MODEL_CLAUDE_3_HAIKU  = "claude-3-haiku-20240229"
)

// MessageContentText | MessageContentFile
type MessageContent interface {
	GetType() string
}

type MessageContentText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (m *MessageContentText) GetType() string {
	return m.Type
}

type MessageContentFileSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type MessageContentFile struct {
	Type   string                   `json:"type"`
	Source MessageContentFileSource `json:"source,omitempty"`
}

func (m *MessageContentFile) GetType() string {
	return m.Type
}

type Message struct {
	Role    string           `json:"role"`
	Content []MessageContent `json:"content"`
	// Content string `json:"content"`
}

type MessagesRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	System    string    `json:"system,omitempty"`
	MaxTokens int       `json:"max_tokens"`
	Metadata  struct {
		UserID string `json:"user_id"`
	} `json:"metadata,omitempty"`
	StopSequences []string `json:"stop_sequences,omitempty"`
	Stream        bool     `json:"stream,omitempty"`
	Temperature   float64  `json:"temperature,omitempty"`
	TopP          float64  `json:"top_p,omitempty"`
	TopK          int      `json:"top_k,omitempty"`
}

type MessagesResponse struct {
	ID           string               `json:"id"`
	Type         string               `json:"type"`
	Role         string               `json:"role"`
	Content      []MessageContentText `json:"content"`
	Model        string               `json:"model"`
	StopReaon    string               `json:"stop_reason"`
	StopSequence string               `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

type MessagesStreamResponseContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type MessagesStreamResponseDelta struct {
	Type         string `json:"type"`
	Text         string `json:"text"`
	StopReaon    string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

type MessagesStream struct {
	*MessagesStreamReader
}

type MessagesStreamResponse struct {
	Type         string                             `json:"type"`
	Message      MessagesResponse                   `json:"message,omitempty"`
	Index        int                                `json:"index,omitempty"`
	ContentBlock MessagesStreamResponseContentBlock `json:"content_block,omitempty"`
	Delta        MessagesStreamResponseDelta        `json:"delta,omitempty"`
	Error        struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type MessagesResponseError struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}
