package embeddings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Client represents an embeddings client
type Client struct {
	ollamaURL string
	model     string
}

// NewClient creates a new embeddings client
func NewClient(ollamaURL string, model string) *Client {
	return &Client{
		ollamaURL: ollamaURL,
		model:     model,
	}
}

// DefaultClient creates a new embeddings client with default settings
func DefaultClient() *Client {
	return NewClient("http://localhost:11434", "all-minilm")
}

// GenerateEmbedding generates a vector embedding for the given text using the Ollama API
func (c *Client) GenerateEmbedding(text string) ([]float32, error) {
	startTime := time.Now()
	log.Printf("Generating embedding for text (length: %d characters): %q", len(text), truncateText(text, 50))

	type EmbedRequest struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}

	type EmbedResponse struct {
		Embedding []float32 `json:"embedding"`
	}

	reqBody := EmbedRequest{
		Model:  c.model,
		Prompt: text,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return nil, err
	}

	resp, err := http.Post(c.ollamaURL+"/api/embeddings", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Printf("Error making HTTP request to Ollama API: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	var embedResponse EmbedResponse
	err = json.Unmarshal(respBody, &embedResponse)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	duration := time.Since(startTime)
	vectorLen := len(embedResponse.Embedding)

	// Log the first 10 numbers of the embedding
	var preview []string
	for i := 0; i < min(10, vectorLen); i++ {
		preview = append(preview, fmt.Sprintf("%.4f", embedResponse.Embedding[i]))
	}
	log.Printf("Generated embedding: %d dimensions in %v", vectorLen, duration)
	log.Printf("First %d values: [%s]", len(preview), strings.Join(preview, ", "))

	return embedResponse.Embedding, nil
}

// GetDimensions returns the dimensions of the embeddings for a given model
func (c *Client) GetDimensions() int {
	// For now, hardcode the dimensions for the all-minilm model
	// In a production environment, this should be retrieved from the model's metadata
	return 384
}

// truncateText truncates text to maxLen characters, adding "..." if truncated
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
