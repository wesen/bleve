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
	"github.com/blevesearch/bleve/v2/search/query"
)

// Document represents the structure of our documents.
type Document struct {
	ID      string    `json:"id"`
	Content string    `json:"content"`
	Vector  []float32 `json:"vector"`
}

// SearchRequest represents the structure of the query DSL
type SearchRequest struct {
	Query   QueryDSL         `yaml:"query"`
	Options *SearchOptions   `yaml:"options,omitempty"`
	Facets  map[string]Facet `yaml:"facets,omitempty"`
}

// QueryDSL represents different types of queries
type QueryDSL struct {
	Match       *MatchQuery       `yaml:"match,omitempty"`
	MatchPhrase *MatchPhraseQuery `yaml:"match_phrase,omitempty"`
	Vector      *VectorQuery      `yaml:"vector,omitempty"`
	Bool        *BooleanQuery     `yaml:"bool,omitempty"`
	Term        *TermQuery        `yaml:"term,omitempty"`
	QueryString *QueryStringQuery `yaml:"query_string,omitempty"`
}

// MatchQuery represents a full-text search query
type MatchQuery struct {
	Field        string  `yaml:"field"`
	Value        string  `yaml:"value"`
	Boost        float64 `yaml:"boost,omitempty"`
	Operator     string  `yaml:"operator,omitempty"`
	Fuzziness    int     `yaml:"fuzziness,omitempty"`
	PrefixLength int     `yaml:"prefix_length,omitempty"`
	Analyzer     string  `yaml:"analyzer,omitempty"`
}

// MatchPhraseQuery represents a phrase search query
type MatchPhraseQuery struct {
	Field    string  `yaml:"field"`
	Value    string  `yaml:"value"`
	Boost    float64 `yaml:"boost,omitempty"`
	Slop     int     `yaml:"slop,omitempty"`
	Analyzer string  `yaml:"analyzer,omitempty"`
}

// VectorQuery represents a vector similarity search
type VectorQuery struct {
	Field  string    `yaml:"field"`
	Text   string    `yaml:"text,omitempty"`
	Vector []float32 `yaml:"vector,omitempty"`
	Model  string    `yaml:"model"`
	K      int       `yaml:"k"`
	Boost  float64   `yaml:"boost,omitempty"`
}

// BooleanQuery represents a boolean combination of queries
type BooleanQuery struct {
	Must               []QueryDSL `yaml:"must,omitempty"`
	Should             []QueryDSL `yaml:"should,omitempty"`
	MustNot            []QueryDSL `yaml:"must_not,omitempty"`
	MinimumShouldMatch int        `yaml:"minimum_should_match,omitempty"`
	Boost              float64    `yaml:"boost,omitempty"`
}

// TermQuery represents an exact term search
type TermQuery struct {
	Field string  `yaml:"field"`
	Value string  `yaml:"value"`
	Boost float64 `yaml:"boost,omitempty"`
}

// QueryStringQuery represents a query string search
type QueryStringQuery struct {
	Query string  `yaml:"query"`
	Boost float64 `yaml:"boost,omitempty"`
}

// SearchOptions represents search configuration options
type SearchOptions struct {
	Size      int          `yaml:"size,omitempty"`
	From      int          `yaml:"from,omitempty"`
	Explain   bool         `yaml:"explain,omitempty"`
	Fields    []string     `yaml:"fields,omitempty"`
	Sort      []SortOption `yaml:"sort,omitempty"`
	Highlight *Highlight   `yaml:"highlight,omitempty"`
}

// SortOption represents a sort configuration
type SortOption struct {
	Field string `yaml:"field"`
	Desc  bool   `yaml:"desc,omitempty"`
}

// Highlight represents highlighting configuration
type Highlight struct {
	Style  string   `yaml:"style,omitempty"`
	Fields []string `yaml:"fields,omitempty"`
}

// Facet represents a facet configuration
type Facet struct {
	Type   string       `yaml:"type"`
	Field  string       `yaml:"field"`
	Size   int          `yaml:"size,omitempty"`
	Ranges []FacetRange `yaml:"ranges,omitempty"`
}

// FacetRange represents a range for numeric or date facets
type FacetRange struct {
	Name  string      `yaml:"name"`
	Min   interface{} `yaml:"min,omitempty"`
	Max   interface{} `yaml:"max,omitempty"`
	Start string      `yaml:"start,omitempty"`
	End   string      `yaml:"end,omitempty"`
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

// buildBleveQuery converts a QueryDSL to a bleve.Query
func buildBleveQuery(q QueryDSL) (query.Query, error) {
	if q.Match != nil {
		query := bleve.NewMatchQuery(q.Match.Value)
		query.SetField(q.Match.Field)
		if q.Match.Boost != 0 {
			query.SetBoost(q.Match.Boost)
		}
		return query, nil
	}

	if q.MatchPhrase != nil {
		query := bleve.NewMatchPhraseQuery(q.MatchPhrase.Value)
		query.SetField(q.MatchPhrase.Field)
		if q.MatchPhrase.Boost != 0 {
			query.SetBoost(q.MatchPhrase.Boost)
		}
		return query, nil
	}

	if q.Vector != nil {
		var queryVector []float32
		var err error

		if q.Vector.Text != "" {
			queryVector, err = generateEmbedding(q.Vector.Text)
			if err != nil {
				return nil, fmt.Errorf("failed to generate vector embedding: %w", err)
			}
			log.Printf("Generated vector embedding: %v", queryVector)
		} else if q.Vector.Vector != nil {
			queryVector = q.Vector.Vector
		} else {
			return nil, fmt.Errorf("either text or vector must be provided for vector query")
		}

		searchRequest := bleve.NewSearchRequest(bleve.NewMatchNoneQuery())
		searchRequest.Size = q.Vector.K
		searchRequest.AddKNN(q.Vector.Field, queryVector, int64(q.Vector.K), q.Vector.Boost)
		return searchRequest.Query, nil
	}

	if q.Bool != nil {
		boolQuery := bleve.NewBooleanQuery()

		for _, must := range q.Bool.Must {
			q, err := buildBleveQuery(must)
			if err != nil {
				return nil, err
			}
			boolQuery.AddMust(q)
		}

		for _, should := range q.Bool.Should {
			q, err := buildBleveQuery(should)
			if err != nil {
				return nil, err
			}
			boolQuery.AddShould(q)
		}

		for _, mustNot := range q.Bool.MustNot {
			q, err := buildBleveQuery(mustNot)
			if err != nil {
				return nil, err
			}
			boolQuery.AddMustNot(q)
		}

		if q.Bool.MinimumShouldMatch > 0 {
			boolQuery.SetMinShould(float64(q.Bool.MinimumShouldMatch))
		}

		if q.Bool.Boost != 0 {
			boolQuery.SetBoost(q.Bool.Boost)
		}

		return boolQuery, nil
	}

	if q.Term != nil {
		query := bleve.NewTermQuery(q.Term.Value)
		query.SetField(q.Term.Field)
		if q.Term.Boost != 0 {
			query.SetBoost(q.Term.Boost)
		}
		return query, nil
	}

	if q.QueryString != nil {
		q_ := bleve.NewQueryStringQuery(q.QueryString.Query)
		if q.QueryString.Boost != 0 {
			q_.SetBoost(q.QueryString.Boost)
		}
		return q_, nil
	}

	return nil, fmt.Errorf("no valid query type found")
}

// handleSearch processes search requests using the query DSL
func handleSearch(index bleve.Index) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// Parse the YAML request
		var searchReq SearchRequest
		if err := yaml.Unmarshal(body, &searchReq); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse YAML: %v", err), http.StatusBadRequest)
			return
		}

		// Build the Bleve query
		bleveQuery, err := buildBleveQuery(searchReq.Query)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to build query: %v", err), http.StatusBadRequest)
			return
		}

		// Create search request
		searchRequest := bleve.NewSearchRequest(bleveQuery)

		// Apply options if provided
		if searchReq.Options != nil {
			if searchReq.Options.Size > 0 {
				searchRequest.Size = searchReq.Options.Size
			}
			if searchReq.Options.From > 0 {
				searchRequest.From = searchReq.Options.From
			}
			if len(searchReq.Options.Fields) > 0 {
				searchRequest.Fields = searchReq.Options.Fields
			}
			if searchReq.Options.Explain {
				searchRequest.Explain = true
			}
			if searchReq.Options.Highlight != nil {
				searchRequest.Highlight = bleve.NewHighlight()
				searchRequest.Highlight.Fields = searchReq.Options.Highlight.Fields
			}
			// Apply sorting
			for _, sort := range searchReq.Options.Sort {
				if sort.Field == "_score" {
					searchRequest.SortBy([]string{"-_score"})
				} else {
					if sort.Desc {
						searchRequest.SortBy([]string{"-" + sort.Field})
					} else {
						searchRequest.SortBy([]string{sort.Field})
					}
				}
			}
		}

		// Execute search
		searchResult, err := index.Search(searchRequest)
		if err != nil {
			http.Error(w, fmt.Sprintf("Search failed: %v", err), http.StatusInternalServerError)
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(searchResult); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	}
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

func main() {
	indexPath := "myindex.bleve"

	// Try to open existing index first
	index, err := bleve.Open(indexPath)
	if err != nil {
		// If index doesn't exist, create a new one
		index, err = createIndex(indexPath)
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
			err := indexDocument(index, doc)
			if err != nil {
				log.Printf("Error indexing document %s: %v", doc.ID, err)
			} else {
				fmt.Printf("Indexed document: %s\n", doc.ID)
			}
		}
	}
	defer index.Close()

	// Set up HTTP server
	http.HandleFunc("/", handleIndex(index))
	http.HandleFunc("/search", handleSearch(index))

	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
