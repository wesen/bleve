# Bleve Vector Search Server

This is an experimental server that demonstrates the vector search capabilities of Bleve, combined with Ollama for generating embeddings. It provides a REST API for searching documents using text similarity, vector similarity, and hybrid approaches.

## Features

- Full-text search with match and phrase queries
- Vector similarity search using embeddings
- Boolean queries combining multiple query types
- Document listing and inspection
- Web interface for exploring the index
- YAML-based query DSL
- Automatic embedding generation using Ollama

## Prerequisites

1. Go 1.19 or later
2. [Ollama](https://ollama.ai/) installed and running
3. The `all-minilm` model pulled in Ollama

```bash
# Install Ollama from https://ollama.ai/
# Pull the all-minilm model
ollama pull all-minilm
```

## Installation

```bash
git clone https://github.com/blevesearch/bleve.git
cd bleve/experiments/embeddings
go build
```

## Running the Server

```bash
# Start Ollama (in a separate terminal)
ollama serve

# Start the server
go run .
```

The server will start on http://localhost:8080.

## API Endpoints

### 1. List Documents (`GET /documents`)

Lists all documents in the index with their fields.

```bash
curl -X GET http://localhost:8080/documents
```

Response:
```json
{
  "total": 3,
  "documents": [
    {
      "id": "doc1",
      "fields": {
        "content": "The quick brown fox jumps over the lazy dog",
        "vector": [0.123, 0.456, ...]
      }
    }
  ]
}
```

### 2. Search (`POST /search`)

Performs a search using the query DSL. Accepts YAML-formatted requests.

#### Match Query Example:
```bash
curl -X POST http://localhost:8080/search \
  -H "Content-Type: application/yaml" \
  -d '
query:
  match:
    field: content
    value: quick fox
    boost: 1.0
options:
  size: 10
  highlight:
    fields: [content]'
```

#### Match Phrase Query Example:
```bash
curl -X POST http://localhost:8080/search \
  -H "Content-Type: application/yaml" \
  -d '
query:
  match_phrase:
    field: content
    value: "quick brown fox"
    boost: 1.0'
```

#### Vector Query Example:
```bash
curl -X POST http://localhost:8080/search \
  -H "Content-Type: application/yaml" \
  -d '
query:
  vector:
    field: vector
    text: "What is the meaning of life?"
    model: all-minilm
    k: 10
    boost: 1.0'
```

#### Boolean Query Example (Hybrid Search):
```bash
curl -X POST http://localhost:8080/search \
  -H "Content-Type: application/yaml" \
  -d '
query:
  bool:
    must:
      - match:
          field: content
          value: fox
    should:
      - vector:
          field: vector
          text: "animal jumping"
          model: all-minilm
          k: 10
    minimum_should_match: 0
options:
  size: 10
  highlight:
    fields: [content]
  sort:
    - field: _score
      desc: true'
```

## Query DSL

The server uses a YAML-based query DSL inspired by Elasticsearch. The main components are:

1. `query`: The search criteria
2. `options`: Search configuration (size, from, fields, highlighting, sorting)
3. `facets`: Aggregation definitions (not implemented yet)

### Query Types

1. **Match Query**
   - Full-text search with analyzed text
   - Supports operators (and/or), fuzziness, and boost

2. **Match Phrase Query**
   - Exact phrase matching
   - Supports boost and slop

3. **Vector Query**
   - Semantic similarity search using embeddings
   - Can use text input (auto-generated embedding) or raw vector
   - Configurable k-nearest neighbors

4. **Boolean Query**
   - Combines multiple queries with boolean logic
   - Supports must, should, and must_not clauses
   - Allows hybrid search combining text and vector queries

### Search Options

- `size`: Number of results to return
- `from`: Pagination offset
- `fields`: Fields to include in results
- `highlight`: Configure result highlighting
- `sort`: Sort results by field or score
- `explain`: Include score explanation

## Architecture

The server is organized into several packages:

- `embeddings`: Handles vector embedding generation using Ollama
- `query`: Contains query DSL types and query building logic
- `templates`: HTML templates for the web interface
- `server.go`: HTTP server and request handling
- `main.go`: Application entry point and index management

## Development

### Adding New Query Types

1. Add the query type to `query/types.go`
2. Implement query building in `query/parser.go`
3. Update documentation

### Modifying Embeddings

The embeddings package (`embeddings/embeddings.go`) can be modified to:
- Use different embedding models
- Change embedding dimensions
- Implement caching
- Add different embedding providers

## Limitations

1. Maximum 1000 documents in listing
2. No authentication/authorization
3. No caching of embeddings
4. Limited error handling
5. No faceted search implementation yet
6. No document updates (only creation on first run)

## Contributing

Contributions are welcome! Please submit issues and pull requests to the main Bleve repository. 