package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"chatbot-service/internal/config"
	"chatbot-service/internal/model"

	"github.com/sashabaranov/go-openai"
)

type AIService struct {
	client *openai.Client
	config *config.OpenAIConfig
}

func NewAIService(cfg *config.OpenAIConfig) *AIService {
	client := openai.NewClient(cfg.APIKey)
	return &AIService{
		client: client,
		config: cfg,
	}
}

type ChatContext struct {
	CourseID     *string `json:"course_id,omitempty"`
	CourseName   string  `json:"course_name,omitempty"`
	LectureName  string  `json:"lecture_name,omitempty"`
	UserRole     string  `json:"user_role,omitempty"`
	PreviousChat string  `json:"previous_chat,omitempty"`
}

func (s *AIService) GenerateResponse(ctx context.Context, messages []*model.ChatMessage, chatContext *ChatContext) (*model.ChatResponse, error) {
	startTime := time.Now()

	// Build OpenAI messages from chat history
	openAIMessages := s.buildOpenAIMessages(messages, chatContext)

	// Create chat completion request
	req := openai.ChatCompletionRequest{
		Model:       s.config.Model,
		Messages:    openAIMessages,
		MaxTokens:   s.config.MaxTokens,
		Temperature: s.config.Temperature,
		Stream:      false,
	}

	// Call OpenAI API
	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned from AI")
	}

	_ = int(time.Since(startTime).Milliseconds())
	content := resp.Choices[0].Message.Content

	return &model.ChatResponse{
		Role:       model.RoleAssistant,
		Content:    content,
		TokensUsed: resp.Usage.TotalTokens,
		CreatedAt:  time.Now(),
	}, nil
}

func (s *AIService) GenerateStreamResponse(ctx context.Context, messages []*model.ChatMessage, chatContext *ChatContext) (<-chan *model.ChatResponse, <-chan error) {
	responseChan := make(chan *model.ChatResponse, 10)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		startTime := time.Now()

		// Build OpenAI messages from chat history
		openAIMessages := s.buildOpenAIMessages(messages, chatContext)

		// Create streaming chat completion request
		req := openai.ChatCompletionRequest{
			Model:       s.config.Model,
			Messages:    openAIMessages,
			MaxTokens:   s.config.MaxTokens,
			Temperature: s.config.Temperature,
			Stream:      true,
		}

		stream, err := s.client.CreateChatCompletionStream(ctx, req)
		if err != nil {
			errorChan <- fmt.Errorf("failed to create AI stream: %w", err)
			return
		}
		defer stream.Close()

		var fullContent strings.Builder
		var totalTokens int

		for {
			response, err := stream.Recv()
			if err != nil {
				if err.Error() == "EOF" {
					// Send final response with complete content
					_ = int(time.Since(startTime).Milliseconds())
					finalResponse := &model.ChatResponse{
						Role:       model.RoleAssistant,
						Content:    fullContent.String(),
						TokensUsed: totalTokens,
						CreatedAt:  time.Now(),
					}
					responseChan <- finalResponse
					return
				}
				errorChan <- fmt.Errorf("stream error: %w", err)
				return
			}

			if len(response.Choices) > 0 {
				delta := response.Choices[0].Delta
				if delta.Content != "" {
					fullContent.WriteString(delta.Content)
					
					// Send partial response
					partialResponse := &model.ChatResponse{
						Role:       model.RoleAssistant,
						Content:    delta.Content,
						TokensUsed: 0, // We'll set this only in the final response
						CreatedAt:  time.Now(),
					}
					responseChan <- partialResponse
				}
			}

			// Note: ChatCompletionStreamResponse doesn't have Usage field in newer versions
			// We'll track tokens differently if needed
		}
	}()

	return responseChan, errorChan
}

func (s *AIService) buildOpenAIMessages(messages []*model.ChatMessage, chatContext *ChatContext) []openai.ChatCompletionMessage {
	var openAIMessages []openai.ChatCompletionMessage

	// Add system message with context
	systemPrompt := s.buildSystemPrompt(chatContext)
	openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemPrompt,
	})

	// Convert chat messages to OpenAI format
	for _, msg := range messages {
		var role string
		switch msg.Role {
		case model.RoleUser:
			role = openai.ChatMessageRoleUser
		case model.RoleAssistant:
			role = openai.ChatMessageRoleAssistant
		case model.RoleSystem:
			role = openai.ChatMessageRoleSystem
		default:
			role = openai.ChatMessageRoleUser
		}

		openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	return openAIMessages
}

func (s *AIService) buildSystemPrompt(context *ChatContext) string {
	var prompt strings.Builder

	prompt.WriteString("You are an AI teaching assistant for an online learning platform. ")
	prompt.WriteString("Your role is to help students understand course material, answer questions, and provide guidance. ")

	if context != nil {
		if context.CourseName != "" {
			prompt.WriteString(fmt.Sprintf("The current course is: %s. ", context.CourseName))
		}
		if context.LectureName != "" {
			prompt.WriteString(fmt.Sprintf("The current lecture is: %s. ", context.LectureName))
		}
		if context.UserRole != "" {
			prompt.WriteString(fmt.Sprintf("The user's role is: %s. ", context.UserRole))
		}
		if context.PreviousChat != "" {
			prompt.WriteString(fmt.Sprintf("Previous conversation context: %s. ", context.PreviousChat))
		}
	}

	prompt.WriteString("\nGuidelines:\n")
	prompt.WriteString("- Be helpful, educational, and encouraging\n")
	prompt.WriteString("- Provide clear explanations with examples when possible\n")
	prompt.WriteString("- Ask clarifying questions if the user's question is unclear\n")
	prompt.WriteString("- Suggest related topics or resources when appropriate\n")
	prompt.WriteString("- If you don't know something, be honest about it\n")
	prompt.WriteString("- Keep responses concise but thorough\n")
	prompt.WriteString("- Use markdown formatting for better readability when needed\n")

	return prompt.String()
}

func (s *AIService) GenerateSessionTitle(ctx context.Context, firstMessage string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Generate a short, descriptive title (max 50 characters) for a chat session based on the user's first message. The title should be clear and concise.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: firstMessage,
		},
	}

	req := openai.ChatCompletionRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    messages,
		MaxTokens:   20,
		Temperature: 0.3,
	}

	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		// Fallback to a simple title if AI call fails
		words := strings.Fields(firstMessage)
		if len(words) > 5 {
			return strings.Join(words[:5], " ") + "...", nil
		}
		return firstMessage, nil
	}

	if len(resp.Choices) == 0 {
		return "New Chat", nil
	}

	title := strings.TrimSpace(resp.Choices[0].Message.Content)
	title = strings.Trim(title, "\"'")
	
	if len(title) > 50 {
		title = title[:47] + "..."
	}

	return title, nil
}