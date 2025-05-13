package main

import (
	"fmt"
	"os"
	"time"

	"github.com/blevesearch/bleve/v2/experiments/embeddings/embeddings"
)

// Example script that demonstrates how to use the Geppetto embeddings client
// to calculate similarity between pieces of text
func main() {
	// Create a cached embeddings client
	client := embeddings.DefaultCachedGeppettoClient(1000)

	// Print client info
	fmt.Printf("Embeddings Client Info:\n")
	fmt.Printf("  Model: %s\n", client.GetModel())
	fmt.Printf("  Dimensions: %d\n", client.GetDimensions())
	fmt.Printf("  Cache Type: %s\n", client.GetCacheStats()["type"])
	fmt.Println()

	// Example texts to compare
	texts := []string{
		"The quick brown fox jumps over the lazy dog",
		"A fast auburn fox leaps above the sleepy canine",
		"To be or not to be, that is the question",
		"Whether 'tis nobler in the mind to suffer the slings and arrows of outrageous fortune",
		"Golang is a programming language designed by Google",
		"Python is a high-level, interpreted programming language",
	}

	// Generate embeddings for each text
	embeddings := make([][]float32, len(texts))
	for i, text := range texts {
		fmt.Printf("Generating embedding for text %d: %s\n", i+1, text)
		embedding, err := client.GenerateEmbedding(text)
		if err != nil {
			fmt.Printf("Error generating embedding: %v\n", err)
			os.Exit(1)
		}
		embeddings[i] = embedding
	}

	fmt.Println("\nSimilarity Matrix:")
	fmt.Printf("%-5s", "")
	for i := range texts {
		fmt.Printf("| %-5d ", i+1)
	}
	fmt.Println("|")

	// Print separator line
	fmt.Printf("%-5s", "")
	for i := range texts {
		fmt.Printf("|%s", "-------")
	}
	fmt.Println("|")

	// Calculate and print similarity matrix
	for i := range texts {
		fmt.Printf("%-5d", i+1)
		for j := range texts {
			similarity := embeddings.ComputeCosineSimilarity(embeddings[i], embeddings[j])
			// Format: show as percentage with 2 decimal places
			fmt.Printf("| %5.2f ", similarity*100)
		}
		fmt.Println("|")
	}

	fmt.Println("\nTexts:")
	for i, text := range texts {
		fmt.Printf("%d: %s\n", i+1, text)
	}

	// Test cache efficiency
	fmt.Println("\nTesting cache efficiency:")
	startTime := time.Now()
	for i := 0; i < 3; i++ {
		for _, text := range texts {
			_, err := client.GenerateEmbedding(text)
			if err != nil {
				fmt.Printf("Error generating embedding: %v\n", err)
				os.Exit(1)
			}
		}
	}
	duration := time.Since(startTime)
	fmt.Printf("Time to generate 18 embeddings (with caching): %v\n", duration)
	fmt.Printf("Average time per embedding: %v\n", duration/18)

	// Print cache stats
	fmt.Println("\nCache Stats after tests:")
	stats := client.GetCacheStats()
	for key, value := range stats {
		fmt.Printf("  %s: %v\n", key, value)
	}
}
