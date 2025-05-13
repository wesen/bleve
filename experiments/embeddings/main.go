package main

import (
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"

	"github.com/blevesearch/bleve/v2/experiments/embeddings/embeddings"
)

// Document represents the structure of our documents.
type Document struct {
	ID      string    `json:"id"`
	Content string    `json:"content"`
	Vector  []float32 `json:"vector"`
}

// indexDocument indexes a document into Bleve, generating and storing its embedding.
func indexDocument(index bleve.Index, doc Document, embeddingsClient *embeddings.GeppettoClient) error {
	embedding, err := embeddingsClient.GenerateEmbedding(doc.Content)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	doc.Vector = embedding

	err = index.Index(doc.ID, doc)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}

	return nil
}

// createIndex creates a new Bleve index with vector search capabilities.
func createIndex(indexPath string, embeddingsClient *embeddings.GeppettoClient) (bleve.Index, error) {
	// Create the index mapping
	indexMapping := bleve.NewIndexMapping()
	documentMapping := bleve.NewDocumentMapping()

	// Text field mappings
	textFieldMapping := bleve.NewTextFieldMapping()
	documentMapping.AddFieldMappingsAt("id", textFieldMapping)
	documentMapping.AddFieldMappingsAt("content", textFieldMapping)

	// Vector field mapping
	vectorFieldMapping := mapping.NewVectorFieldMapping()
	vectorFieldMapping.Dims = embeddingsClient.GetDimensions() // Get dimensions from the embeddings client
	vectorFieldMapping.Similarity = "cosine"
	documentMapping.AddFieldMappingsAt("vector", vectorFieldMapping)

	indexMapping.DefaultMapping = documentMapping

	// Create the index
	index, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}

	return index, nil
}

func main() {
	indexPath := "myindex.bleve"

	// Create embeddings client with disk caching
	cacheDir := "./embeddings_cache"

	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		log.Printf("Warning: Failed to create cache directory: %v. Falling back to memory cache.", err)
	}

	var embeddingsClient *embeddings.GeppettoClient
	var err error

	// Try to create disk-cached client
	embeddingsClient, err = embeddings.DefaultDiskCachedGeppettoClient(cacheDir)
	if err != nil {
		log.Printf("Warning: Failed to create disk-cached embeddings client: %v. Falling back to memory cache.", err)
		// Fall back to memory-cached client
		embeddingsClient = embeddings.DefaultCachedGeppettoClient(1000)
	}

	// Log client information
	log.Printf("Using embeddings client with model: %s, dimensions: %d",
		embeddingsClient.GetModel(), embeddingsClient.GetDimensions())

	// Log cache information
	cacheStats := embeddingsClient.GetCacheStats()
	log.Printf("Cache stats: %v", cacheStats)

	// Try to open existing index first
	index, err := bleve.Open(indexPath)
	if err != nil {
		// If index doesn't exist, create a new one
		index, err = createIndex(indexPath, embeddingsClient)
		if err != nil {
			log.Fatalf("Error creating index: %v", err)
		}

		// Index sample documents only for new index
		documents := []Document{
			{ID: "doc1", Content: "The quick brown fox jumps over the lazy dog"},
			{ID: "doc2", Content: "A journey of a thousand miles begins with a single step"},
			{ID: "doc3", Content: "To be or not to be, that is the question"},
		}

		for _, doc := range documents {
			err := indexDocument(index, doc, embeddingsClient)
			if err != nil {
				log.Printf("Error indexing document %s: %v", doc.ID, err)
			} else {
				fmt.Printf("Indexed document: %s\n", doc.ID)
			}
		}
	}
	defer index.Close()

	// Create and start server with the embeddings client
	server := NewServer(index, embeddingsClient)
	if err := server.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
