package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"gemini2openai/config"
	"gemini2openai/converter"
	"gemini2openai/logger"
	"gemini2openai/model"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// ChatHandler handles the chat completions request.
type ChatHandler struct {
	cfg config.Config
}

// NewChatHandler creates a new ChatHandler.
func NewChatHandler(cfg config.Config) *ChatHandler {
	return &ChatHandler{cfg: cfg}
}

// Handle is the actual http.HandlerFunc for the chat handler.
func (h *ChatHandler) Handle(w http.ResponseWriter, r *http.Request) {
	logEntry := &logger.RequestLogEntry{
		Timestamp: time.Now(),
		ClientIP:  r.RemoteAddr,
	}
	// Use defer to ensure logging happens at the end of the request
	defer func() {
		logger.Log.Info().Object("request", logEntry).Send()
	}()

	if h.cfg.Auth.OpenAIAPIKey != "" {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			logEntry.StatusCode = http.StatusUnauthorized
			logEntry.Error = "Authorization header is required"
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != h.cfg.Auth.OpenAIAPIKey {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			logEntry.StatusCode = http.StatusUnauthorized
			logEntry.Error = "Invalid API key"
			return
		}
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		logEntry.StatusCode = http.StatusMethodNotAllowed
		logEntry.Error = "Only POST method is allowed"
		return
	}

	// Read the body to log it, then replace it for further processing
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		logEntry.StatusCode = http.StatusInternalServerError
		logEntry.Error = "Failed to read request body"
		return
	}
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var openAIReq model.OpenAIRequest
	if err := json.Unmarshal(bodyBytes, &openAIReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logEntry.StatusCode = http.StatusBadRequest
		logEntry.Error = "Invalid request body"
		return
	}
	logEntry.RequestToProxy = &openAIReq
	logEntry.ModelRequested = openAIReq.Model

	if openAIReq.Stream {
		h.stream(w, &openAIReq, logEntry)
	} else {
		h.nonStream(w, &openAIReq, logEntry)
	}
}

func (h *ChatHandler) nonStream(w http.ResponseWriter, openAIReq *model.OpenAIRequest, logEntry *logger.RequestLogEntry) {
	geminiReq := converter.OpenAIRequestToGemini(openAIReq)
	logEntry.RequestToUpstream = geminiReq

	geminiReqBytes, err := json.Marshal(geminiReq)
	if err != nil {
		http.Error(w, "Failed to marshal Gemini request", http.StatusInternalServerError)
		logEntry.StatusCode = http.StatusInternalServerError
		logEntry.Error = "Failed to marshal Gemini request"
		return
	}

	modelName := "gemini-pro" // Default model
	if openAIReq.Model != "" {
		modelName = openAIReq.Model
	}
	fullURL := h.cfg.Gemini.BaseURL + "/v1beta/models/" + modelName + ":generateContent?key=" + h.cfg.Gemini.APIKey
	logEntry.UpstreamURL = fullURL

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(geminiReqBytes))
	if err != nil {
		http.Error(w, "Failed to create request to Gemini API", http.StatusInternalServerError)
		logEntry.StatusCode = http.StatusInternalServerError
		logEntry.Error = "Failed to create request to Gemini API"
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to Gemini API", http.StatusInternalServerError)
		logEntry.StatusCode = http.StatusInternalServerError
		logEntry.Error = "Failed to send request to Gemini API"
		return
	}
	defer resp.Body.Close()

	upstreamBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from Gemini API", http.StatusInternalServerError)
		logEntry.StatusCode = http.StatusInternalServerError
		logEntry.Error = "Failed to read response from Gemini API"
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error from upstream API", resp.StatusCode)
		logEntry.StatusCode = resp.StatusCode
		logEntry.Error = "Error from upstream API: " + string(upstreamBody)
		return
	}

	var geminiResp model.GeminiResponse
	if err := json.Unmarshal(upstreamBody, &geminiResp); err != nil {
		http.Error(w, "Failed to unmarshal Gemini response", http.StatusInternalServerError)
		logEntry.StatusCode = http.StatusInternalServerError
		logEntry.Error = "Failed to unmarshal Gemini response"
		return
	}
	logEntry.ResponseFromUpstream = &geminiResp

	openAIResp := converter.GeminiResponseToOpenAI(&geminiResp)
	logEntry.FinalResponse = openAIResp
	logEntry.StatusCode = http.StatusOK

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(openAIResp); err != nil {
		logEntry.Error = "Failed to write final response"
		return
	}
}

func (h *ChatHandler) stream(w http.ResponseWriter, openAIReq *model.OpenAIRequest, logEntry *logger.RequestLogEntry) {
	geminiReq := converter.OpenAIRequestToGemini(openAIReq)
	logEntry.RequestToUpstream = geminiReq

	geminiReqBytes, err := json.Marshal(geminiReq)
	if err != nil {
		http.Error(w, "Failed to marshal Gemini request", http.StatusInternalServerError)
		return
	}

	modelName := "gemini-pro"
	if openAIReq.Model != "" {
		modelName = openAIReq.Model
	}
	fullURL := h.cfg.Gemini.BaseURL + "/v1beta/models/" + modelName + ":streamGenerateContent?alt=sse&key=" + h.cfg.Gemini.APIKey
	logEntry.UpstreamURL = fullURL

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(geminiReqBytes))
	if err != nil {
		http.Error(w, "Failed to create request to Gemini API", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to Gemini API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		http.Error(w, "Error from upstream API: "+string(bodyBytes), resp.StatusCode)
		logEntry.StatusCode = resp.StatusCode
		logEntry.Error = "Error from upstream API: " + string(bodyBytes)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			var geminiResp model.GeminiResponse
			if err := json.Unmarshal([]byte(data), &geminiResp); err != nil {
				continue
			}

			openAIStreamResp := converter.GeminiResponseToOpenAIStream(&geminiResp, modelName)
			respBytes, err := json.Marshal(openAIStreamResp)
			if err != nil {
				continue
			}

			fmt.Fprintf(w, "data: %s\n\n", respBytes)
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}

	if err := scanner.Err(); err != nil {
		// Handle scanner error
	}

	fmt.Fprintf(w, "data: [DONE]\n\n")
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}