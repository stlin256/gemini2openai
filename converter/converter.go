package converter

import (
	"crypto/rand"
	"encoding/hex"
	"gemini2openai/model"
	"strings"
	"time"
)

// OpenAIRequestToGemini converts an OpenAI Chat API request to a Google Gemini API request.
func OpenAIRequestToGemini(req *model.OpenAIRequest) *model.GeminiRequest {
	geminiReq := &model.GeminiRequest{
		Contents: make([]model.GeminiContent, len(req.Messages)),
		GenerationConfig: model.GeminiGenerationConfig{
			Temperature:     req.Temperature,
			TopP:            req.TopP,
			MaxOutputTokens: req.MaxTokens,
		},
		Model: req.Model,
	}

	for i, msg := range req.Messages {
		role := "user"
		if strings.ToLower(msg.Role) == "assistant" {
			role = "model"
		}
		geminiReq.Contents[i] = model.GeminiContent{
			Role: role,
			Parts: []model.GeminiPart{
				{Text: msg.Content},
			},
		}
	}

	return geminiReq
}

// GeminiResponseToOpenAI converts a Google Gemini API response to an OpenAI Chat API response.
func GeminiResponseToOpenAI(geminiResp *model.GeminiResponse) *model.OpenAIResponse {
	openAIResp := &model.OpenAIResponse{
		ID:      "genai-" + generateRandomString(12), // Placeholder ID
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "gemini-converted", // Placeholder model name
		Choices: make([]model.OpenAIChoice, len(geminiResp.Candidates)),
		Usage:   model.OpenAIUsage{}, // Placeholder usage
	}

	var totalTokens int
	for i, candidate := range geminiResp.Candidates {
		totalTokens += candidate.TokenCount
		openAIResp.Choices[i] = model.OpenAIChoice{
			Index: candidate.Index,
			Message: model.OpenAIMessage{
				Role:    "assistant",
				Content: candidate.Content.Parts[0].Text,
			},
			FinishReason: candidate.FinishReason,
		}
	}

	openAIResp.Usage.TotalTokens = totalTokens
	// Note: Gemini API does not provide prompt_tokens and completion_tokens separately.
	// We are putting the total in TotalTokens as a best effort.
	openAIResp.Usage.CompletionTokens = totalTokens 

	return openAIResp
}

// generateRandomString creates a random hex string of a given length.
func generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
// GeminiResponseToOpenAIStream converts a single Gemini API response chunk to an OpenAI Chat API stream chunk.
func GeminiResponseToOpenAIStream(geminiResp *model.GeminiResponse, modelName string) *model.OpenAIStreamResponse {
	finishReason := geminiResp.Candidates[0].FinishReason
	var finishReasonPtr *string
	if finishReason != "" {
		finishReasonPtr = &finishReason
	}

	return &model.OpenAIStreamResponse{
		ID:      "genai-" + generateRandomString(12),
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   modelName,
		Choices: []model.OpenAIStreamChoice{
			{
				Index: 0,
				Delta: model.OpenAIStreamDelta{
					Content: geminiResp.Candidates[0].Content.Parts[0].Text,
				},
				FinishReason: finishReasonPtr,
			},
		},
	}
}