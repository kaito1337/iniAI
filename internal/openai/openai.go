package openai

import (
	"context"
	"fmt"
	"inivoice/internal/constants"
	"inivoice/libs"
	"strings"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	client  *openai.Client
	enabled bool
	mu      sync.Mutex
	logger  libs.Logger
}

func New(apiKey string, logger libs.Logger) *OpenAI {
	s := &OpenAI{
		client: openai.NewClient(apiKey),
		logger: logger,
	}

	if apiKey != "" {
		s.enabled = true
	}

	return s
}

func (s *OpenAI) Summarize(text string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.enabled {
		s.logger.Info("Summary disabled")
		return "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	response, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf("%s%s", constants.PROMPT, text),
			},
		},
		MaxTokens:   256,
		Temperature: 0.7,
		TopP:        1,
	})
	if err != nil {
		s.logger.Errorf("Summary error: %s", err.Error())
		return "", err
	}
	return strings.TrimSpace(response.Choices[0].Message.Content), nil
}

func (s *OpenAI) SpeechToText(filepath string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.enabled {
		s.logger.Info("Transcription disabled")
		return "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	response, err := s.client.CreateTranscription(ctx, openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: filepath,
	})
	if err != nil {
		s.logger.Errorf("Transcription error: %s", err.Error())
		return "", err
	}
	return response.Text, nil
}
