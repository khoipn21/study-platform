package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"chatbot-service/internal/config"
	"chatbot-service/internal/model"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

type AIService struct {
	llm    llms.Model
	config *config.GeminiConfig
}

func NewAIService(cfg *config.GeminiConfig) (*AIService, error) {
	llm, err := googleai.New(
		context.Background(),
		googleai.WithAPIKey(cfg.APIKey),
		googleai.WithDefaultModel(cfg.Model),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &AIService{
		llm:    llm,
		config: cfg,
	}, nil
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

	// Build prompt from chat history
	prompt := s.buildPrompt(messages, chatContext)

	// Generate response using Gemini
	content, err := llms.GenerateFromSinglePrompt(
		ctx,
		s.llm,
		prompt,
		llms.WithTemperature(s.config.Temperature),
		llms.WithMaxTokens(s.config.MaxTokens),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %w", err)
	}

	_ = int(time.Since(startTime).Milliseconds())

	// Estimate token usage (Gemini doesn't always provide exact counts)
	tokensUsed := len(strings.Fields(prompt))/2 + len(strings.Fields(content))/2

	return &model.ChatResponse{
		Role:       model.RoleAssistant,
		Content:    content,
		TokensUsed: tokensUsed,
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

		// Build prompt from chat history
		prompt := s.buildPrompt(messages, chatContext)

		var fullContent strings.Builder

		// Stream response using Gemini
		_, err := s.llm.GenerateContent(
			ctx,
			[]llms.MessageContent{
				llms.TextParts(llms.ChatMessageTypeHuman, prompt),
			},
			llms.WithTemperature(s.config.Temperature),
			llms.WithMaxTokens(s.config.MaxTokens),
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				content := string(chunk)
				if content != "" {
					fullContent.WriteString(content)
					
					// Send partial response
					partialResponse := &model.ChatResponse{
						Role:       model.RoleAssistant,
						Content:    content,
						TokensUsed: 0,
						CreatedAt:  time.Now(),
					}
					responseChan <- partialResponse
				}
				return nil
			}),
		)

		if err != nil {
			errorChan <- fmt.Errorf("stream error: %w", err)
			return
		}

		// Send final response with complete content
		_ = int(time.Since(startTime).Milliseconds())
		tokensUsed := len(strings.Fields(prompt))/2 + len(strings.Fields(fullContent.String()))/2

		finalResponse := &model.ChatResponse{
			Role:       model.RoleAssistant,
			Content:    fullContent.String(),
			TokensUsed: tokensUsed,
			CreatedAt:  time.Now(),
		}
		responseChan <- finalResponse
	}()

	return responseChan, errorChan
}

func (s *AIService) buildPrompt(messages []*model.ChatMessage, chatContext *ChatContext) string {
	var prompt strings.Builder

	// Add system message with context
	systemPrompt := s.buildSystemPrompt(chatContext)
	prompt.WriteString(systemPrompt)
	prompt.WriteString("\n\n")

	// Convert chat messages to prompt format
	for _, msg := range messages {
		switch msg.Role {
		case model.RoleUser:
			prompt.WriteString(fmt.Sprintf("Student: %s\n", msg.Content))
		case model.RoleAssistant:
			prompt.WriteString(fmt.Sprintf("Assistant: %s\n", msg.Content))
		case model.RoleSystem:
			prompt.WriteString(fmt.Sprintf("System: %s\n", msg.Content))
		}
	}

	prompt.WriteString("Assistant: ")

	return prompt.String()
}

func (s *AIService) buildSystemPrompt(context *ChatContext) string {
	var prompt strings.Builder

	prompt.WriteString("You are an AI teaching assistant specialized in helping students solve academic problems. ")
	prompt.WriteString("Your primary role is to assist with homework, assignments, understanding complex concepts, ")
	prompt.WriteString("and providing step-by-step solutions to academic questions across various subjects. ")

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

	prompt.WriteString("\n\nCore Guidelines for Academic Assistance:\n")
	prompt.WriteString("1. PROBLEM SOLVING: Break down complex problems into manageable steps\n")
	prompt.WriteString("2. EXPLANATIONS: Provide clear, detailed explanations with examples and analogies\n")
	prompt.WriteString("3. STEP-BY-STEP: Show your work and reasoning process for mathematical and logical problems\n")
	prompt.WriteString("4. VERIFICATION: Help students verify their solutions and understand common mistakes\n")
	prompt.WriteString("5. CONCEPTUAL UNDERSTANDING: Focus on helping students understand WHY, not just HOW\n")
	prompt.WriteString("6. SUBJECT AREAS: Assist with mathematics, science, programming, languages, humanities, and more\n")
	prompt.WriteString("7. CLARIFICATION: Ask follow-up questions if the problem statement is unclear\n")
	prompt.WriteString("8. RESOURCES: Suggest relevant learning materials, practice problems, or topics to review\n")
	prompt.WriteString("9. ENCOURAGEMENT: Be patient, supportive, and encouraging - learning is a journey\n")
	prompt.WriteString("10. ACADEMIC INTEGRITY: Guide students toward understanding rather than just giving answers\n")
	prompt.WriteString("\nFormatting:\n")
	prompt.WriteString("- Use markdown for better readability (headers, lists, code blocks, math notation)\n")
	prompt.WriteString("- For math: Use LaTeX notation when appropriate (e.g., $x^2 + y^2 = r^2$)\n")
	prompt.WriteString("- For code: Use proper code blocks with syntax highlighting\n")
	prompt.WriteString("- Structure responses with clear sections for complex topics\n")

	return prompt.String()
}

func (s *AIService) GenerateSessionTitle(ctx context.Context, firstMessage string) (string, error) {
	prompt := fmt.Sprintf("Generate a short, descriptive title (max 50 characters) for a chat session based on this message: '%s'. Return only the title, nothing else.", firstMessage)

	title, err := llms.GenerateFromSinglePrompt(
		ctx,
		s.llm,
		prompt,
		llms.WithTemperature(0.3),
		llms.WithMaxTokens(20),
	)
	
	if err != nil {
		// Fallback to a simple title if AI call fails
		words := strings.Fields(firstMessage)
		if len(words) > 5 {
			return strings.Join(words[:5], " ") + "...", nil
		}
		return firstMessage, nil
	}

	title = strings.TrimSpace(title)
	title = strings.Trim(title, "\"'")
	
	if len(title) > 50 {
		title = title[:47] + "..."
	}

	if title == "" {
		return "New Chat", nil
	}

	return title, nil
}