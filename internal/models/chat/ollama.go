package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/models/utils/ollama"
	"github.com/Tencent/WeKnora/internal/types"
	ollamaapi "github.com/ollama/ollama/api"
)

// OllamaChat implements Ollama-based chat
type OllamaChat struct {
	modelName     string
	modelID       string
	ollamaService *ollama.OllamaService
}

// NewOllamaChat creates Ollama chat instance
func NewOllamaChat(config *ChatConfig, ollamaService *ollama.OllamaService) (*OllamaChat, error) {
	return &OllamaChat{
		modelName:     config.ModelName,
		modelID:       config.ModelID,
		ollamaService: ollamaService,
	}, nil
}

// convertMessages converts message format to Ollama API format
func (c *OllamaChat) convertMessages(messages []Message) []ollamaapi.Message {
	ollamaMessages := make([]ollamaapi.Message, 0, len(messages))
	for _, msg := range messages {
		msgOllama := ollamaapi.Message{
			Role:      msg.Role,
			Content:   msg.Content,
			ToolCalls: c.toolCallFrom(msg.ToolCalls),
		}
		if msg.Role == "tool" {
			msgOllama.ToolName = msg.Name
		}
		ollamaMessages = append(ollamaMessages, msgOllama)
	}
	return ollamaMessages
}

// buildChatRequest builds chat request parameters
func (c *OllamaChat) buildChatRequest(messages []Message, opts *ChatOptions, isStream bool) *ollamaapi.ChatRequest {
	// Set streaming flag
	streamFlag := isStream

	// Build request parameters
	chatReq := &ollamaapi.ChatRequest{
		Model:    c.modelName,
		Messages: c.convertMessages(messages),
		Stream:   &streamFlag,
		Options:  make(map[string]interface{}),
	}

	// Add optional parameters
	if opts != nil {
		if opts.Temperature > 0 {
			chatReq.Options["temperature"] = opts.Temperature
		}
		if opts.TopP > 0 {
			chatReq.Options["top_p"] = opts.TopP
		}
		if opts.MaxTokens > 0 {
			chatReq.Options["num_predict"] = opts.MaxTokens
		}
		if opts.Thinking != nil {
			chatReq.Think = &ollamaapi.ThinkValue{
				Value: *opts.Thinking,
			}
		}
		if len(opts.Format) > 0 {
			chatReq.Format = opts.Format
		}
		if len(opts.Tools) > 0 {
			chatReq.Tools = c.toolFrom(opts.Tools)
		}
	}

	return chatReq
}

// Chat performs non-streaming chat
func (c *OllamaChat) Chat(ctx context.Context, messages []Message, opts *ChatOptions) (*types.ChatResponse, error) {
	// Ensure model is available
	if err := c.ensureModelAvailable(ctx); err != nil {
		return nil, err
	}

	// Build request parameters
	chatReq := c.buildChatRequest(messages, opts, false)

	// Log request
	logger.GetLogger(ctx).Infof("Sending chat request to model %s", c.modelName)

	var responseContent string
	var toolCalls []types.LLMToolCall
	var promptTokens, completionTokens int

	// Send request using Ollama client
	err := c.ollamaService.Chat(ctx, chatReq, func(resp ollamaapi.ChatResponse) error {
		responseContent = resp.Message.Content
		toolCalls = c.toolCallTo(resp.Message.ToolCalls)

		// Get token count
		if resp.EvalCount > 0 {
			promptTokens = resp.PromptEvalCount
			completionTokens = resp.EvalCount - promptTokens
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("chat request failed: %w", err)
	}

	// 构建响应
	return &types.ChatResponse{
		Content:   responseContent,
		ToolCalls: toolCalls,
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		},
	}, nil
}

// ChatStream performs streaming chat
func (c *OllamaChat) ChatStream(
	ctx context.Context,
	messages []Message,
	opts *ChatOptions,
) (<-chan types.StreamResponse, error) {
	// Ensure model is available
	if err := c.ensureModelAvailable(ctx); err != nil {
		return nil, err
	}

	// Build request parameters
	chatReq := c.buildChatRequest(messages, opts, true)

	// Log request
	logger.GetLogger(ctx).Infof("Sending streaming chat request to model %s", c.modelName)

	// Create streaming response channel
	streamChan := make(chan types.StreamResponse)

	// Start goroutine to handle streaming response
	go func() {
		defer close(streamChan)

		err := c.ollamaService.Chat(ctx, chatReq, func(resp ollamaapi.ChatResponse) error {
			if resp.Message.Content != "" {
				streamChan <- types.StreamResponse{
					ResponseType: types.ResponseTypeAnswer,
					Content:      resp.Message.Content,
					Done:         false,
				}
			}

			if len(resp.Message.ToolCalls) > 0 {
				streamChan <- types.StreamResponse{
					ResponseType: types.ResponseTypeToolCall,
					ToolCalls:    c.toolCallTo(resp.Message.ToolCalls),
					Done:         false,
				}
			}

			if resp.Done {
				streamChan <- types.StreamResponse{
					ResponseType: types.ResponseTypeAnswer,
					Done:         true,
				}
			}

			return nil
		})
		if err != nil {
			logger.GetLogger(ctx).Errorf("Streaming chat request failed: %v", err)
			// Send error response
			streamChan <- types.StreamResponse{
				ResponseType: types.ResponseTypeError,
				Content:      err.Error(),
				Done:         true,
			}
		}
	}()

	return streamChan, nil
}

// ensureModelAvailable ensures model is available
func (c *OllamaChat) ensureModelAvailable(ctx context.Context) error {
	logger.GetLogger(ctx).Infof("Ensuring model %s is available", c.modelName)
	return c.ollamaService.EnsureModelAvailable(ctx, c.modelName)
}

// GetModelName gets model name
func (c *OllamaChat) GetModelName() string {
	return c.modelName
}

// GetModelID gets model ID
func (c *OllamaChat) GetModelID() string {
	return c.modelID
}

// toolFrom converts this module's Tool to Ollama's Tool
func (c *OllamaChat) toolFrom(tools []Tool) ollamaapi.Tools {
	if len(tools) == 0 {
		return nil
	}
	ollamaTools := make(ollamaapi.Tools, 0, len(tools))
	for _, tool := range tools {
		function := ollamaapi.ToolFunction{
			Name:        tool.Function.Name,
			Description: tool.Function.Description,
		}
		if len(tool.Function.Parameters) > 0 {
			_ = json.Unmarshal(tool.Function.Parameters, &function.Parameters)
		}

		ollamaTools = append(ollamaTools, ollamaapi.Tool{
			Type:     tool.Type,
			Function: function,
		})
	}
	return ollamaTools
}

// toolTo converts Ollama's Tool to this module's Tool
func (c *OllamaChat) toolTo(ollamaTools ollamaapi.Tools) []Tool {
	if len(ollamaTools) == 0 {
		return nil
	}
	tools := make([]Tool, 0, len(ollamaTools))
	for _, tool := range ollamaTools {
		paramsBytes, _ := json.Marshal(tool.Function.Parameters)
		tools = append(tools, Tool{
			Type: tool.Type,
			Function: FunctionDef{
				Name:        tool.Function.Name,
				Description: tool.Function.Description,
				Parameters:  paramsBytes,
			},
		})
	}
	return tools
}

// toolCallFrom converts this module's ToolCall to Ollama's ToolCall
func (c *OllamaChat) toolCallFrom(toolCalls []ToolCall) []ollamaapi.ToolCall {
	if len(toolCalls) == 0 {
		return nil
	}
	ollamaToolCalls := make([]ollamaapi.ToolCall, 0, len(toolCalls))
	for _, tc := range toolCalls {
		var args map[string]interface{}
		if tc.Function.Arguments != "" {
			_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
		}
		ollamaToolCalls = append(ollamaToolCalls, ollamaapi.ToolCall{
			Function: ollamaapi.ToolCallFunction{
				Index:     tools2i(tc.ID),
				Name:      tc.Function.Name,
				Arguments: args,
			},
		})
	}
	return ollamaToolCalls
}

// toolCallTo converts Ollama's ToolCall to this module's ToolCall
func (c *OllamaChat) toolCallTo(ollamaToolCalls []ollamaapi.ToolCall) []types.LLMToolCall {
	if len(ollamaToolCalls) == 0 {
		return nil
	}
	toolCalls := make([]types.LLMToolCall, 0, len(ollamaToolCalls))
	for _, tc := range ollamaToolCalls {
		argsBytes, _ := json.Marshal(tc.Function.Arguments)
		toolCalls = append(toolCalls, types.LLMToolCall{
			ID:   tooli2s(tc.Function.Index),
			Type: "function",
			Function: types.FunctionCall{
				Name:      tc.Function.Name,
				Arguments: string(argsBytes),
			},
		})
	}
	return toolCalls
}

func tooli2s(i int) string {
	return strconv.Itoa(i)
}

func tools2i(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
