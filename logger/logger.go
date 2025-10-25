package logger

import (
	"gemini2openai/config"
	"gemini2openai/model"
	"github.com/rs/zerolog"
	"io"
	"os"
	"time"
)

var Log zerolog.Logger

// RequestLogEntry defines the structure for a single log entry.
type RequestLogEntry struct {
	Timestamp        time.Time             `json:"timestamp"`
	ClientIP         string                `json:"client_ip"`
	ModelRequested   string                `json:"model_requested"`
	UpstreamURL      string                `json:"upstream_url"`
	RequestToProxy   *model.OpenAIRequest  `json:"request_to_proxy"`
	RequestToUpstream *model.GeminiRequest  `json:"request_to_upstream"`
	ResponseFromUpstream *model.GeminiResponse `json:"response_from_upstream"`
	FinalResponse    *model.OpenAIResponse `json:"final_response"`
	StatusCode       int                   `json:"status_code"`
	Error            string                `json:"error,omitempty"`
}

// MarshalZerologObject makes RequestLogEntry implement zerolog.LogObjectMarshaler
func (e *RequestLogEntry) MarshalZerologObject(evt *zerolog.Event) {
	evt.Time("timestamp", e.Timestamp).
		Str("client_ip", e.ClientIP).
		Str("model_requested", e.ModelRequested).
		Str("upstream_url", e.UpstreamURL).
		Interface("request_to_proxy", e.RequestToProxy).
		Interface("request_to_upstream", e.RequestToUpstream).
		Interface("response_from_upstream", e.ResponseFromUpstream).
		Interface("final_response", e.FinalResponse).
		Int("status_code", e.StatusCode).
		Str("error", e.Error)
}


// InitLogger initializes the global logger based on the configuration.
func InitLogger(cfg config.LogConfig) {
	var output io.Writer
	if cfg.Enabled {
		file, err := os.OpenFile(cfg.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			// Fallback to console if file cannot be opened
			output = os.Stdout
		} else {
			output = file
		}
	} else {
		// If logging is disabled, discard all logs
		output = io.Discard
	}

	Log = zerolog.New(output).With().Timestamp().Logger()
}