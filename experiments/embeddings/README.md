# Bleve Vector Search Server with Geppetto Embeddings

This is an experimental server that demonstrates the vector search capabilities of Bleve, combined with Geppetto's embeddings package for generating embeddings with caching support. It provides a REST API for searching documents using text similarity, vector similarity, and hybrid approaches.

## Enhanced Features

- Full-text search with match and phrase queries
- Vector similarity search using embeddings
- Boolean queries combining multiple query types
- Document listing and inspection
- Web interface for exploring the index
- YAML-based query DSL
- **Automatic embedding generation using Ollama via Geppetto**
- **Caching support for embeddings (in-memory and disk-based)**
- **Advanced embeddings management through dedicated endpoints**

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

## Embeddings Architecture

The embeddings package now utilizes the Geppetto embeddings library for enhanced functionality:

### Key Components

1. **GeppettoClient**: 
   - Wrapper around the Geppetto embeddings provider
   - Supports multiple provider implementations
   - Compatible with the original client interface

2. **Embedding Providers**:
   - Direct Ollama provider
   - Memory-cached provider
   - Disk-cached provider

3. **Caching Mechanisms**:
   - In-memory LRU cache for short-lived applications
   - Persistent disk cache for long-running services
   - Configurable cache sizes and expiration

### Embeddings Endpoints

#### 1. Generate Embedding (`POST /embeddings`)

Generates an embedding vector for provided text.

```bash
curl -X POST http://localhost:8080/embeddings \
  -H "Content-Type: application/json" \
  -d '{"text": "Hello world"}'
```

Response:
```json
{
  "text": "Hello world",
  "dimensions": 384,
  "model": "all-minilm",
  "embedding": [0.123, 0.456, ...]
}
```

#### 2. Get Cache Stats (`GET /cache/stats`)

Returns the current state of the embeddings cache.

```bash
curl -X GET http://localhost:8080/cache/stats
```

Response:
```json
{
  "type": "memory",
  "size": 42,
  "max_size": 1000,
  "model": "all-minilm",
  "dimensions": 384
}
```

#### 3. Clear Cache (`POST /cache/clear`)

Clears the embeddings cache.

```bash
curl -X POST http://localhost:8080/cache/clear
```

Response:
```json
{
  "success": true,
  "message": "Cache cleared successfully"
}
```

## Web Interface

The web interface has been enhanced with:

1. **Embeddings Information Panel**:
   - Model details
   - Dimensions
   - Cache type and statistics

2. **Cache Management**:
   - Button to clear the cache
   - Button to refresh cache statistics

3. **Embedding Generation Tool**:
   - Text input for generating embeddings directly
   - Display of generated embeddings

4. **Improved YAML Query Interface**:
   - Syntax-highlighted textarea for YAML queries
   - Load example button for quick queries

## Examples

The `examples` directory contains scripts demonstrating the embeddings functionality:

- `similarity.go`: Demonstrates text similarity calculation using the cached embeddings client

### Running Examples

```bash
go run examples/similarity.go
```

## Vector Query Example

The vector query support has been enhanced to include more efficient caching:

```yaml
query:
  vector:
    field: vector
    text: "What is the meaning of life?"
    model: all-minilm
    k: 10
    boost: 1.0
options:
  size: 10
  highlight:
    fields: [content]
```

## Hybrid Search Example

```yaml
query:
  bool:
    must:
      - match:
          field: content
          value: fox
          boost: 1.0
    should:
      - vector:
          field: vector
          text: "animal jumping"
          model: all-minilm
          k: 10
          boost: 2.0
    minimum_should_match: 0
options:
  size: 10
  highlight:
    fields: [content]
  sort:
    - field: _score
      desc: true
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