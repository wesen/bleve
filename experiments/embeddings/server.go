package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/blevesearch/bleve/v2"
	"gopkg.in/yaml.v2"

	"github.com/blevesearch/bleve/v2/experiments/embeddings/query"
	"github.com/blevesearch/bleve/v2/experiments/embeddings/templates"
)

// Server represents the HTTP server and its dependencies
type Server struct {
	index bleve.Index
}

// NewServer creates a new server instance
func NewServer(index bleve.Index) *Server {
	return &Server{
		index: index,
	}
}

// handleIndex handles the index page
func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get index mapping
		mapping := s.index.Mapping()
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
		if err := templates.Templates.Base.Execute(w, data); err != nil {
			log.Printf("Template execution error: %v", err)
		}
	}
}

// handleSearch handles the search endpoint
func (s *Server) handleSearch() http.HandlerFunc {
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
		var searchReq query.SearchRequest
		if err := yaml.Unmarshal(body, &searchReq); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse YAML: %v", err), http.StatusBadRequest)
			return
		}

		// Build the Bleve query
		bleveQuery, err := query.BuildBleveQuery(searchReq.Query)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to build query: %v", err), http.StatusBadRequest)
			return
		}

		// Create search request
		searchRequest := bleve.NewSearchRequest(bleveQuery)

		// Apply options
		query.ApplySearchOptions(searchRequest, searchReq.Options)

		// Execute search
		searchResult, err := s.index.Search(searchRequest)
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

// handleListDocuments handles the list documents endpoint
func (s *Server) handleListDocuments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Create a match all query
		matchAllQuery := bleve.NewMatchAllQuery()
		searchRequest := bleve.NewSearchRequest(matchAllQuery)

		// Request all fields
		searchRequest.Fields = []string{"*"}
		searchRequest.Size = 1000 // Limit to 1000 documents for safety

		// Execute search
		searchResult, err := s.index.Search(searchRequest)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list documents: %v", err), http.StatusInternalServerError)
			return
		}

		// Format response
		type DocumentInfo struct {
			ID     string                 `json:"id"`
			Fields map[string]interface{} `json:"fields"`
		}

		docs := make([]DocumentInfo, len(searchResult.Hits))
		for i, hit := range searchResult.Hits {
			docs[i] = DocumentInfo{
				ID:     hit.ID,
				Fields: hit.Fields,
			}
		}

		response := struct {
			Total     uint64         `json:"total"`
			Documents []DocumentInfo `json:"documents"`
		}{
			Total:     searchResult.Total,
			Documents: docs,
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	// Set up routes
	http.HandleFunc("/", s.handleIndex())
	http.HandleFunc("/search", s.handleSearch())
	http.HandleFunc("/documents", s.handleListDocuments())

	// Start server
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, nil)
}
