package service

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type AIService struct {
	client *openai.Client
}

func NewAIService(apiKey string) *AIService {
	return &AIService{
		client: openai.NewClient(apiKey),
	}
}

func (ai *AIService) ProcessScheduleImage(ctx context.Context, imageURL string) (string, error) {
	// Здесь можно добавить обработку изображения с расписанием
	// Пока просто заглушка
	return "Расписание успешно обработано", nil
}

func (ai *AIService) HandleQuestion(ctx context.Context, question string) (string, error) {
	resp, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Ты помощник для студентов. Отвечай кратко и по делу.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: question,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
