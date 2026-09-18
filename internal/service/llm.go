package service

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type LLMService struct {
	client *genai.Client
	model string
}

func NewLLMService(ctx context.Context, projectID, location, model string) (*LLMService, error){
	client, err := genai.NewClient(ctx,&genai.ClientConfig{
		Backend: genai.BackendVertexAI,
		Project: projectID,
		Location: location,
	})
	
	if err != nil{
		return nil, fmt.Errorf("genai client olusturulamadi: %w",err)
	}
	return &LLMService{client: client, model: model}, nil
}

func (s *LLMService) GenerateDescription(ctx context.Context, diff string) (string, error){
	prompt := buildPrompt(diff)

	result, err := s.client.Models.GenerateContent(ctx,s.model,genai.Text(prompt),nil)
	if err != nil{
		return "", fmt.Errorf("LLM cagrisi basarisiz: %w",err)
	}
	return result.Text(), nil
}

func buildPrompt(diff string) string {
	return fmt.Sprintf(`Sen kıdemli bir yazılım mühendisisin. Aşağıdaki git diff'ini incele. Neyin, neden değiştiğini açıklayan, teknik ama sade bir dille profesyonel bir Pull Request açıklaması yaz. Maddeleme kullan. Sadece açıklama metnini döndür, başka bir şey ekleme.

Diff:
%s`, diff)
}