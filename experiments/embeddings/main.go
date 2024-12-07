package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search"
)

// Document represents the structure of our documents.
type Document struct {
	ID      string    `json:"id"`
	Content string    `json:"content"`
	Vector  []float32 `json:"vector"`
}

// generateEmbedding generates a vector embedding for the given text using the Ollama API.
func generateEmbedding(text string) ([]float32, error) {
	startTime := time.Now()
	log.Printf("Generating embedding for text (length: %d characters)", len(text))

	type EmbedRequest struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}

	type EmbedResponse struct {
		Embedding []float32 `json:"embedding"`
	}

	reqBody := EmbedRequest{
		Model:  "all-minilm",
		Prompt: text,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return nil, err
	}

	resp, err := http.Post("http://localhost:11434/api/embeddings", "application/json", bytes.NewBuffer(reqBytes))
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
	log.Printf("Successfully generated embedding: %d dimensions in %v", vectorLen, duration)

	return embedResponse.Embedding, nil
}

// indexDocument indexes a document into Bleve, generating and storing its embedding.
func indexDocument(index bleve.Index, doc Document) error {
	embedding, err := generateEmbedding(doc.Content)
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
func createIndex(indexPath string) (bleve.Index, error) {
	// Create the index mapping
	indexMapping := bleve.NewIndexMapping()
	documentMapping := bleve.NewDocumentMapping()

	// Text field mappings
	textFieldMapping := bleve.NewTextFieldMapping()
	documentMapping.AddFieldMappingsAt("id", textFieldMapping)
	documentMapping.AddFieldMappingsAt("content", textFieldMapping)

	// Vector field mapping
	vectorFieldMapping := mapping.NewVectorFieldMapping()
	vectorFieldMapping.Dims = 384 // Set to match nomic-embed-text dimensions
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

// searchIndex performs a hybrid search on the Bleve index.
func searchIndex(index bleve.Index, queryString string, k int) (*bleve.SearchResult, error) {
	// Generate embedding for the query
	queryVector, err := generateEmbedding(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Create text query
	textQuery := bleve.NewMatchQuery(queryString)
	textQuery.SetField("content")

	// Create search request with text query
	searchRequest := bleve.NewSearchRequest(textQuery)
	searchRequest.Size = k
	searchRequest.Fields = []string{"id", "content"}
	searchRequest.Highlight = bleve.NewHighlight()
	searchRequest.Highlight.Fields = []string{"content"}

	// Add vector search
	searchRequest.AddKNN("vector", queryVector, int64(k), 1.0)

	// Execute search
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Log the search result fields for debugging
	log.Printf("Search result: %+v", searchResult)

	return searchResult, nil
}

// Templates
const baseHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Bleve Search</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 p-8">
    <div class="max-w-4xl mx-auto">
        <h1 class="text-3xl font-bold mb-8">Bleve Search Demo</h1>
        
        <div class="bg-white rounded-lg shadow p-6 mb-8">
            <h2 class="text-xl font-semibold mb-4">Index Mappings</h2>
            <pre class="bg-gray-50 p-4 rounded overflow-auto max-h-96"><code>{{.Mappings}}</code></pre>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
            <h2 class="text-xl font-semibold mb-4">Search</h2>
            <form hx-post="/search" hx-target="#results" class="mb-6">
                <div class="flex gap-4">
                    <input type="text" name="query" placeholder="Enter your search query" 
                           class="flex-1 px-4 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500">
                    <button type="submit" 
                            class="px-6 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500">
                        Search
                    </button>
                </div>
            </form>
            <div id="results"></div>
        </div>
    </div>
</body>
</html>`

const searchResultsHTML = `
{{if .Error}}
    <div class="text-red-500 mb-4">{{.Error}}</div>
{{else}}
    <div class="space-y-4">
        {{range .Hits}}
            <div class="border rounded p-4 hover:bg-gray-50">
                <div class="font-semibold mb-2">Document ID: {{.ID}}</div>
                <div class="text-gray-600">Score: {{printf "%.4f" .Score}}</div>
                {{if .Fragments}}
                    <div class="mt-2 text-sm">
                        {{range $field, $fragments := .Fragments}}
                            {{range $fragments}}
                                <div class="mt-1">... {{.}} ...</div>
                            {{end}}
                        {{end}}
                    </div>
                {{end}}
            </div>
        {{else}}
            <div class="text-gray-500">No results found</div>
        {{end}}
    </div>
{{end}}`

var (
	baseTemplate          = template.Must(template.New("base").Parse(baseHTML))
	searchResultsTemplate = template.Must(template.New("results").Parse(searchResultsHTML))
)

// HTTP Handlers
func handleIndex(index bleve.Index) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get index mapping
		mapping := index.Mapping()
		mappingBytes, err := json.MarshalIndent(mapping, "", "  ")
		if err != nil {
			http.Error(w, "Failed to get index mapping", http.StatusInternalServerError)
			return
		}

		// Convert JSON to YAML for better readability
		var mapData interface{}
		if err := json.Unmarshal(mappingBytes, &mapData); err != nil {
			http.Error(w, "Failed to parse mapping", http.StatusInternalServerError)
			return
		}

		mappingYAML, err := yaml.Marshal(mapData)
		if err != nil {
			http.Error(w, "Failed to convert mapping to YAML", http.StatusInternalServerError)
			return
		}

		data := struct {
			Mappings string
		}{
			Mappings: string(mappingYAML),
		}

		w.Header().Set("Content-Type", "text/html")
		if err := baseTemplate.Execute(w, data); err != nil {
			log.Printf("Template execution error: %v", err)
		}
	}
}

func handleSearch(index bleve.Index) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := r.FormValue("query")
		if query == "" {
			http.Error(w, "Query is required", http.StatusBadRequest)
			return
		}

		results, err := searchIndex(index, query, 10)
		data := struct {
			Error string
			Hits  search.DocumentMatchCollection
		}{}

		if err != nil {
			data.Error = fmt.Sprintf("Search failed: %v", err)
		} else {
			data.Hits = results.Hits
		}

		w.Header().Set("Content-Type", "text/html")
		if err := searchResultsTemplate.Execute(w, data); err != nil {
			log.Printf("Template execution error: %v", err)
		}
	}
}

func main() {
	indexPath := "myindex.bleve"

	// Create a new index
	index, err := createIndex(indexPath)
	if err != nil {
		log.Fatalf("Error creating index: %v", err)
	}
	defer index.Close()

	// Index some sample documents
	documents := []Document{
		{ID: "doc1", Content: "The quick brown fox jumps over the lazy dog"},
		{ID: "doc2", Content: "A journey of a thousand miles begins with a single step"},
		{ID: "doc3", Content: "To be or not to be, that is the question"},
	}

	for _, doc := range documents {
		err := indexDocument(index, doc)
		if err != nil {
			log.Printf("Error indexing document %s: %v", doc.ID, err)
		} else {
			fmt.Printf("Indexed document: %s\n", doc.ID)
		}
	}

	// Set up HTTP server
	http.HandleFunc("/", handleIndex(index))
	http.HandleFunc("/search", handleSearch(index))

	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
