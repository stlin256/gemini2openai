package model

// GeminiRequest represents the request body for the Google Gemini API.
type GeminiRequest struct {
	Contents         []GeminiContent      `json:"contents"`
	GenerationConfig GeminiGenerationConfig `json:"generationConfig,omitempty"`
	Model            string               `json:"model,omitempty"` // Added model field
}

// GeminiContent represents the content of a single turn in the conversation.
type GeminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a single part of the content.
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiGenerationConfig allows you to configure the model's generation parameters.
type GeminiGenerationConfig struct {
	Temperature     float64 `json:"temperature,omitempty"`
	TopP            float64 `json:"topP,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

// GeminiResponse represents the response from the Google Gemini API.
type GeminiResponse struct {
	Candidates     []GeminiCandidate `json:"candidates"`
	PromptFeedback struct {
		BlockReason   string `json:"blockReason"`
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"promptFeedback"`
}

// GeminiCandidate represents a single candidate in the response.
type GeminiCandidate struct {
	Content       GeminiContent `json:"content"`
	FinishReason  string        `json:"finishReason"`
	Index         int           `json:"index"`
	SafetyRatings []struct {
		Category    string `json:"category"`
		Probability string `json:"probability"`
	} `json:"safetyRatings"`
	TokenCount int `json:"tokenCount"`
}