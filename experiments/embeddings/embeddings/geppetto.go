package embeddings

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/go-go-golems/geppetto/pkg/embeddings"
)

// GeppettoClient represents an embeddings client using Geppetto's embeddings package
type GeppettoClient struct {
	provider embeddings.Provider
	model    string
	dims     int
}

// NewGeppettoClient creates a new Geppetto embeddings client
func NewGeppettoClient(provider embeddings.Provider) *GeppettoClient {
	model := provider.GetModel()
	return &GeppettoClient{
		provider: provider,
		model:    model.Name,
		dims:     model.Dimensions,
	}
}

// NewOllamaGeppettoClient creates a new Geppetto client with Ollama provider
func NewOllamaGeppettoClient(baseURL, model string, dims int) *GeppettoClient {
	provider := embeddings.NewOllamaProvider(baseURL, model, dims)
	return NewGeppettoClient(provider)
}

// NewCachedOllamaGeppettoClient creates a new Geppetto client with cached Ollama provider
func NewCachedOllamaGeppettoClient(baseURL, model string, dims int, cacheSize int) *GeppettoClient {
	baseProvider := embeddings.NewOllamaProvider(baseURL, model, dims)
	cachedProvider := embeddings.NewCachedProvider(baseProvider, cacheSize)
	return NewGeppettoClient(cachedProvider)
}

// NewDiskCachedOllamaGeppettoClient creates a new Geppetto client with disk-cached Ollama provider
func NewDiskCachedOllamaGeppettoClient(baseURL, model string, dims int, cacheDir string, maxSize int64, maxEntries int) (*GeppettoClient, error) {
	baseProvider := embeddings.NewOllamaProvider(baseURL, model, dims)

	// Create options for disk cache provider
	opts := []embeddings.DiskCacheProviderOption{
		embeddings.WithDirectory(cacheDir),
	}

	if maxSize > 0 {
		opts = append(opts, embeddings.WithMaxSize(maxSize))
	}

	if maxEntries > 0 {
		opts = append(opts, embeddings.WithMaxEntries(maxEntries))
	}

	// Create disk cache provider
	diskProvider, err := embeddings.NewDiskCacheProvider(baseProvider, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create disk cache provider: %w", err)
	}

	return NewGeppettoClient(diskProvider), nil
}

// DefaultGeppettoClient creates a new Geppetto client with default settings
func DefaultGeppettoClient() *GeppettoClient {
	return NewOllamaGeppettoClient("http://localhost:11434", "all-minilm", 384)
}

// DefaultCachedGeppettoClient creates a new Geppetto client with in-memory cache
func DefaultCachedGeppettoClient(cacheSize int) *GeppettoClient {
	baseProvider := embeddings.NewOllamaProvider("http://localhost:11434", "all-minilm", 384)
	cachedProvider := embeddings.NewCachedProvider(baseProvider, cacheSize)
	return NewGeppettoClient(cachedProvider)
}

// DefaultDiskCachedGeppettoClient creates a new Geppetto client with disk cache
func DefaultDiskCachedGeppettoClient(cacheDir string) (*GeppettoClient, error) {
	return NewDiskCachedOllamaGeppettoClient(
		"http://localhost:11434",
		"all-minilm",
		384,
		cacheDir,
		1<<30, // 1GB max size
		10000, // 10,000 entries max
	)
}

// GenerateEmbedding generates a vector embedding for the given text
func (c *GeppettoClient) GenerateEmbedding(text string) ([]float32, error) {
	startTime := time.Now()
	log.Printf("Generating embedding for text (length: %d characters): %q", len(text), truncateText(text, 50))

	// Create context
	ctx := context.Background()

	// Generate embedding using Geppetto provider
	embedding, err := c.provider.GenerateEmbedding(ctx, text)
	if err != nil {
		log.Printf("Error generating embedding: %v", err)
		return nil, err
	}

	duration := time.Since(startTime)
	vectorLen := len(embedding)

	// Log the first 10 numbers of the embedding
	log.Printf("Generated embedding: %d dimensions in %v", vectorLen, duration)

	if vectorLen > 0 {
		preview := make([]string, 0, min(10, vectorLen))
		for i := 0; i < min(10, vectorLen); i++ {
			preview = append(preview, fmt.Sprintf("%.4f", embedding[i]))
		}
		log.Printf("First %d values: [%s]", len(preview), fmt.Sprintf("%v", preview))
	}

	return embedding, nil
}

// GetDimensions returns the dimensions of the embeddings
func (c *GeppettoClient) GetDimensions() int {
	return c.dims
}

// GetModel returns the model name
func (c *GeppettoClient) GetModel() string {
	return c.model
}

// ClearCache clears the cache if the provider supports it
func (c *GeppettoClient) ClearCache() error {
	// Check if the provider is a CachedProvider
	if cachedProvider, ok := c.provider.(*embeddings.CachedProvider); ok {
		cachedProvider.ClearCache()
		return nil
	}

	// Check if the provider is a DiskCacheProvider
	if diskProvider, ok := c.provider.(*embeddings.DiskCacheProvider); ok {
		return diskProvider.ClearCache()
	}

	return fmt.Errorf("provider does not support cache clearing")
}

// GetCacheStats returns cache statistics if available
func (c *GeppettoClient) GetCacheStats() map[string]interface{} {
	stats := make(map[string]interface{})

	// Check if the provider is a CachedProvider
	if cachedProvider, ok := c.provider.(*embeddings.CachedProvider); ok {
		stats["type"] = "memory"
		stats["size"] = cachedProvider.Size()
		stats["max_size"] = cachedProvider.MaxSize()
		return stats
	}

	// Check if the provider is a DiskCacheProvider
	if diskProvider, ok := c.provider.(*embeddings.DiskCacheProvider); ok {
		stats["type"] = "disk"
		// Add disk-specific stats if needed
		return stats
	}

	stats["type"] = "none"
	return stats
}

// ComputeCosineSimilarity calculates cosine similarity between two vectors
func ComputeCosineSimilarity(a, b []float32) float64 {
	var dotProduct float64
	var normA float64
	var normB float64

	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	return dotProduct / (calculateSqrt(normA) * calculateSqrt(normB))
}

// calculateSqrt is a helper function to calculate square root
func calculateSqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return math.Sqrt(x)
}
