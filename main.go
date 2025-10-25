package main

import (
	"fmt"
	"log"
	"net/http"

	"gemini2openai/config"
	"gemini2openai/handler"
	"gemini2openai/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.InitLogger(cfg.Log)

	chatHandler := handler.NewChatHandler(cfg)

	http.Handle("/v1/chat/completions", corsMiddleware(http.HandlerFunc(chatHandler.Handle)))

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on port %s...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}