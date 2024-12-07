# Bleve Mapping Types Guide

This document provides a comprehensive overview of the different mapping types available in Bleve and their configurations.

## Document Mapping

The document mapping is the top-level mapping that describes how a document should be indexed. It has the following properties:

```json
{
  "enabled": true,
  "dynamic": true,
  "properties": {},
  "fields": [],
  "default_analyzer": "",
  "struct_tag_key": "json"
}
```

- `enabled`: When false, the entire section is ignored
- `dynamic`: When false, only explicitly mapped fields are indexed
- `properties`: Map of sub-document mappings
- `fields`: Array of field mappings
- `default_analyzer`: Default analyzer to use for fields in this document
- `struct_tag_key`: Key to use when looking for field names in struct tags (defaults to "json")

## Field Types

### Text Field

The most common field type for indexing textual content.

```json
{
  "name": "description",
  "type": "text",
  "analyzer": "standard",
  "store": true,
  "index": true,
  "include_term_vectors": true,
  "include_in_all": true,
  "docvalues": false
}
```

- `analyzer`: The analyzer to use (e.g., "standard", "keyword", "en")
- `store`: Whether to store the field value
- `index`: Whether to index the field
- `include_term_vectors`: Whether to include term vectors
- `include_in_all`: Whether to include in _all field
- `docvalues`: Whether to store doc values for sorting/faceting

### Numeric Field

For indexing numeric values.

```json
{
  "name": "age",
  "type": "number",
  "store": true,
  "index": true,
  "docvalues": true
}
```

### Date Field

For indexing dates and times.

```json
{
  "name": "created_at",
  "type": "datetime",
  "store": true,
  "index": true,
  "docvalues": true
}
```

### Boolean Field

For indexing boolean values.

```json
{
  "name": "active",
  "type": "boolean",
  "store": true,
  "index": true
}
```

### Geopoint Field

For indexing geographical points.

```json
{
  "name": "location",
  "type": "geopoint",
  "store": true,
  "index": true
}
```

Supports multiple input formats:
- String: "lat,lon"
- Array: [lon, lat]
- Object: {"lon"/"lng": X, "lat": Y}

### Geoshape Field

For indexing geographical shapes.

```json
{
  "name": "area",
  "type": "geoshape",
  "store": true,
  "index": true
}
```

### Vector Field (v2.4.0+)

For indexing vector embeddings.

```json
{
  "name": "embedding",
  "type": "vector",
  "dims": 768,
  "similarity": "cosine",
  "optimization": "latency"
}
```

Key parameters:
- `dims`: Vector dimensionality (1-4096)
- `similarity`: Similarity metric ("cosine", "dot_product", "l2_norm")
- `optimization`: Index optimization strategy ("latency", "memory_efficient", "recall")

Alternative format using base64 encoding:

```json
{
  "name": "embedding",
  "type": "vector_base64",
  "dims": 768,
  "similarity": "cosine",
  "optimization": "latency"
}
```

### IP Field

For indexing IP addresses.

```json
{
  "name": "ip_address",
  "type": "IP",
  "store": true,
  "index": true
}
```

## Common Field Options

Most field types support these common options:

- `store`: Whether to store the original field value
- `index`: Whether to make the field searchable
- `docvalues`: Whether to enable sorting/faceting on the field
- `include_in_all`: Whether to include the field in the _all field
- `include_term_vectors`: Whether to store term vectors for the field

## Default Settings

Global defaults for dynamic fields:

```go
IndexDynamic = true      // Whether to index dynamic fields
StoreDynamic = true      // Whether to store dynamic fields
DocValuesDynamic = true  // Whether to enable docvalues for dynamic fields
```

## Analyzers

Common built-in analyzers:

- `standard`: Standard analyzer with reasonable defaults
- `keyword`: Treats the entire field value as a single token
- `en`: English analyzer with stop words and stemming
- Custom analyzers can be configured in the index mapping

## Example Complete Mapping

```json
{
  "types": {
    "document": {
      "properties": {
        "title": {
          "fields": [
            {
              "name": "title",
              "type": "text",
              "analyzer": "en",
              "store": true,
              "index": true
            }
          ]
        },
        "description": {
          "fields": [
            {
              "name": "description",
              "type": "text",
              "analyzer": "en",
              "store": true,
              "index": true
            }
          ]
        },
        "created_at": {
          "fields": [
            {
              "name": "created_at",
              "type": "datetime",
              "store": true,
              "index": true
            }
          ]
        },
        "location": {
          "fields": [
            {
              "name": "location",
              "type": "geopoint",
              "store": true,
              "index": true
            }
          ]
        },
        "embedding": {
          "fields": [
            {
              "name": "embedding",
              "type": "vector",
              "dims": 768,
              "similarity": "cosine"
            }
          ]
        }
      }
    }
  },
  "default_type": "document",
  "default_analyzer": "standard",
  "default_datetime_parser": "dateTimeOptional",
  "default_field": "_all",
  "byte_array_converter": "json"
}
``` 

## Creating Mappings in Go

This section demonstrates how to create different types of mappings programmatically using Bleve's Go API.

### Basic Index Mapping

```go
import (
    "github.com/blevesearch/bleve/v2"
    "github.com/blevesearch/bleve/v2/mapping"
)

// Create a new index mapping
indexMapping := bleve.NewIndexMapping()

// Configure default settings
indexMapping.DefaultAnalyzer = "en"
indexMapping.DefaultDateTimeParser = "dateTimeOptional"
```

### Document Mapping

```go
// Create a document mapping
documentMapping := bleve.NewDocumentMapping()

// Configure document mapping settings
documentMapping.Enabled = true
documentMapping.Dynamic = true

// Add it to the index mapping
indexMapping.AddDocumentMapping("article", documentMapping)
```

### Field Mappings

```go
// Text field
titleFieldMapping := bleve.NewTextFieldMapping()
titleFieldMapping.Analyzer = "en"
titleFieldMapping.Store = true
titleFieldMapping.Index = true
titleFieldMapping.IncludeTermVectors = true
titleFieldMapping.IncludeInAll = true

// Numeric field
ageFieldMapping := bleve.NewNumericFieldMapping()
ageFieldMapping.Store = true
ageFieldMapping.Index = true
ageFieldMapping.DocValues = true

// Date field
createdFieldMapping := bleve.NewDateTimeFieldMapping()
createdFieldMapping.Store = true
createdFieldMapping.Index = true
createdFieldMapping.DocValues = true

// Boolean field
activeFieldMapping := bleve.NewBooleanFieldMapping()
activeFieldMapping.Store = true
activeFieldMapping.Index = true

// Geopoint field
locationFieldMapping := bleve.NewGeoPointFieldMapping()
locationFieldMapping.Store = true
locationFieldMapping.Index = true

// Vector field
vectorFieldMapping := bleve.NewVectorFieldMapping()
vectorFieldMapping.Dims = 768
vectorFieldMapping.Similarity = "cosine"
vectorFieldMapping.Optimization = "latency"

// Add fields to document mapping
documentMapping.AddFieldMappingsAt("title", titleFieldMapping)
documentMapping.AddFieldMappingsAt("age", ageFieldMapping)
documentMapping.AddFieldMappingsAt("created", createdFieldMapping)
documentMapping.AddFieldMappingsAt("active", activeFieldMapping)
documentMapping.AddFieldMappingsAt("location", locationFieldMapping)
documentMapping.AddFieldMappingsAt("embedding", vectorFieldMapping)
```

### Complete Example

Here's a complete example that puts it all together:

```go
package main

import (
    "github.com/blevesearch/bleve/v2"
    "github.com/blevesearch/bleve/v2/mapping"
)

func buildIndexMapping() mapping.IndexMapping {
    // Create the index mapping
    indexMapping := bleve.NewIndexMapping()
    
    // Create document mapping for articles
    articleMapping := bleve.NewDocumentMapping()
    
    // Text fields
    titleField := bleve.NewTextFieldMapping()
    titleField.Analyzer = "en"
    titleField.Store = true
    titleField.Index = true
    
    contentField := bleve.NewTextFieldMapping()
    contentField.Analyzer = "en"
    contentField.Store = true
    contentField.Index = true
    contentField.IncludeTermVectors = true
    
    // Date field
    publishedField := bleve.NewDateTimeFieldMapping()
    publishedField.Store = true
    publishedField.Index = true
    publishedField.DocValues = true
    
    // Vector field for embeddings
    embeddingField := bleve.NewVectorFieldMapping()
    embeddingField.Dims = 768
    embeddingField.Similarity = "cosine"
    
    // Add fields to article mapping
    articleMapping.AddFieldMappingsAt("title", titleField)
    articleMapping.AddFieldMappingsAt("content", contentField)
    articleMapping.AddFieldMappingsAt("published", publishedField)
    articleMapping.AddFieldMappingsAt("embedding", embeddingField)
    
    // Add article mapping to index
    indexMapping.AddDocumentMapping("article", articleMapping)
    
    // Set default type
    indexMapping.DefaultType = "article"
    
    return indexMapping
}

func main() {
    // Create index with mapping
    mapping := buildIndexMapping()
    index, err := bleve.New("example.bleve", mapping)
    if err != nil {
        panic(err)
    }
    
    // Index a document
    article := struct {
        Type      string     `json:"type"`
        Title     string     `json:"title"`
        Content   string     `json:"content"`
        Published string     `json:"published"`
        Embedding []float32  `json:"embedding"`
    }{
        Type:      "article",
        Title:     "Example Article",
        Content:   "This is the article content.",
        Published: "2023-01-01T12:00:00Z",
        Embedding: make([]float32, 768),
    }
    
    if err := index.Index("article1", article); err != nil {
        panic(err)
    }
}
```

### Sub-document Mappings

You can also create nested document structures:

```go
// Create parent document mapping
userMapping := bleve.NewDocumentMapping()

// Create child document mapping
addressMapping := bleve.NewDocumentMapping()

// Add fields to address mapping
streetField := bleve.NewTextFieldMapping()
cityField := bleve.NewTextFieldMapping()
addressMapping.AddFieldMappingsAt("street", streetField)
addressMapping.AddFieldMappingsAt("city", cityField)

// Add address as a sub-document to user
userMapping.AddSubDocumentMapping("address", addressMapping)

// Add user mapping to index
indexMapping.AddDocumentMapping("user", userMapping)
```

### Custom Analyzer Configuration

You can also configure custom analyzers:

```go
// Create custom analyzer
customAnalyzer := &analysis.Analyzer{
    Tokenizer: tokenizer.NewUnicodeTokenizer(),
    TokenFilters: []analysis.TokenFilter{
        token.NewLowerCaseFilter(),
        token.NewStopTokensFilter([]string{"the", "a", "an"}),
    },
}

// Add analyzer to index mapping
indexMapping.AddCustomAnalyzer("my_analyzer", customAnalyzer)

// Use custom analyzer in field mapping
fieldMapping := bleve.NewTextFieldMapping()
fieldMapping.Analyzer = "my_analyzer"
```

This section shows the most common ways to create and configure mappings using Bleve's Go API. The mappings created programmatically are equivalent to the JSON configurations shown in the previous sections. 

## YAML Mapping DSL

This section describes a YAML-based DSL (Domain Specific Language) for defining Bleve mappings, similar to Elasticsearch's mapping syntax.

### Basic Structure

```yaml
settings:
  default_analyzer: standard
  default_datetime_parser: dateTimeOptional
  default_field: _all
  default_type: _default
  byte_array_converter: json
  store_dynamic: true
  index_dynamic: true
  docvalues_dynamic: true

mappings:
  dynamic: true
  enabled: true
  
  # Document type definitions
  types:
    article:
      dynamic: true
      enabled: true
      default_analyzer: standard
      properties:
        title:
          type: text
          analyzer: en
          store: true
          index: true
          include_term_vectors: true
          include_in_all: true
          docvalues: false
```

### Field Type Definitions

#### Text Fields

```yaml
title:
  type: text
  analyzer: en  # standard, keyword, en, etc.
  store: true
  index: true
  include_term_vectors: true
  include_in_all: true
  docvalues: false
  
description:
  type: text
  analyzer: standard
  store: true
  index: true
  include_term_vectors: false
  include_in_all: true
  
keywords:
  type: text
  analyzer: keyword  # Treats entire field as single token
  store: true
  index: true
```

#### Numeric Fields

```yaml
age:
  type: number
  store: true
  index: true
  docvalues: true

price:
  type: number
  store: true
  index: true
  docvalues: true
  precision_step: 4  # Optional precision step for range queries
```

#### Date Fields

```yaml
created_at:
  type: datetime
  store: true
  index: true
  docvalues: true
  datetime_parser: dateTimeOptional  # default parser

updated_at:
  type: datetime
  store: true
  index: true
  docvalues: true
  datetime_parser: datetime_optional  # alternative format
```

#### Boolean Fields

```yaml
active:
  type: boolean
  store: true
  index: true

is_published:
  type: boolean
  store: true
  index: true
  docvalues: true  # Enable for faceting/aggregations
```

#### Geographic Fields

```yaml
location:
  type: geopoint
  store: true
  index: true
  precision_step: 6  # Optional precision for geospatial indexing

area:
  type: geoshape
  store: true
  index: true
```

#### Vector Fields

```yaml
embedding:
  type: vector
  dims: 768
  similarity: cosine  # cosine, dot_product, l2_norm
  optimization: latency  # latency, memory_efficient, recall
  store: true

embedding_base64:
  type: vector_base64
  dims: 768
  similarity: cosine
  optimization: memory_efficient
```

#### IP Fields

```yaml
ip_address:
  type: IP
  store: true
  index: true
  docvalues: true
```

### Nested Document Mappings

```yaml
user:
  properties:
    name:
      type: text
      analyzer: standard
      store: true
      index: true
    
    address:
      enabled: true
      dynamic: true
      properties:
        street:
          type: text
          analyzer: standard
          store: true
        city:
          type: text
          analyzer: keyword
          store: true
        location:
          type: geopoint
          store: true
          index: true
    
    contacts:
      enabled: true
      dynamic: false  # Strict mapping for contacts
      properties:
        email:
          type: text
          analyzer: keyword
          store: true
        phone:
          type: text
          analyzer: keyword
          store: true
```

### Custom Analyzer Definitions

```yaml
settings:
  analysis:
    analyzers:
      my_custom_analyzer:
        type: custom
        tokenizer: unicode
        token_filters:
          - lowercase
          - stop_en
          - stemmer_en
    
    token_filters:
      stop_en:
        type: stop
        tokens: [a, an, the, in, on, at, for]
      
      stemmer_en:
        type: stemmer
        language: english
    
    char_filters:
      html_strip:
        type: html

mappings:
  types:
    article:
      properties:
        content:
          type: text
          analyzer: my_custom_analyzer
```

### Complete Example

Here's a complete example showing various mapping features:

```yaml
settings:
  default_analyzer: standard
  default_datetime_parser: dateTimeOptional
  default_field: _all
  store_dynamic: true
  index_dynamic: true
  
  analysis:
    analyzers:
      my_text_analyzer:
        type: custom
        tokenizer: unicode
        token_filters:
          - lowercase
          - stop_en
          - stemmer_en

mappings:
  dynamic: true
  enabled: true
  
  types:
    article:
      dynamic: true
      enabled: true
      default_analyzer: standard
      
      properties:
        title:
          type: text
          analyzer: en
          store: true
          index: true
          include_term_vectors: true
        
        content:
          type: text
          analyzer: my_text_analyzer
          store: true
          index: true
          include_term_vectors: true
        
        summary:
          type: text
          analyzer: standard
          store: true
          index: true
        
        published_date:
          type: datetime
          store: true
          index: true
          docvalues: true
        
        author:
          properties:
            name:
              type: text
              analyzer: standard
              store: true
            
            email:
              type: text
              analyzer: keyword
              store: true
            
            bio:
              type: text
              analyzer: en
              store: true
        
        tags:
          type: text
          analyzer: keyword
          store: true
          index: true
          docvalues: true
        
        rating:
          type: number
          store: true
          index: true
          docvalues: true
        
        location:
          type: geopoint
          store: true
          index: true
        
        embedding:
          type: vector
          dims: 768
          similarity: cosine
          optimization: latency
        
        is_featured:
          type: boolean
          store: true
          index: true
```

### DSL to JSON Conversion

The YAML DSL can be converted to Bleve's JSON mapping format using standard YAML-to-JSON converters. The structure is designed to match Bleve's internal mapping representation while providing a more readable and maintainable format.

Key differences from Elasticsearch's mapping DSL:
1. Field options follow Bleve's naming conventions (e.g., `include_term_vectors` instead of `term_vector`)
2. Vector field configurations are specific to Bleve's implementation
3. Analyzer configurations follow Bleve's analysis chain structure
4. Geographic field types use Bleve's `geopoint` and `geoshape` terminology

The DSL supports all of Bleve's mapping features while providing a more concise and readable format for configuration. 

## Creating Index from JSON Mapping

You can create a Bleve index from a JSON mapping file in several ways. Here are the common approaches:

### Using JSON File Directly

```go
package main

import (
    "encoding/json"
    "os"

    "github.com/blevesearch/bleve/v2"
    "github.com/blevesearch/bleve/v2/mapping"
)

func main() {
    // Read the mapping file
    mappingBytes, err := os.ReadFile("mapping.json")
    if err != nil {
        panic(err)
    }

    // Parse the mapping
    var indexMapping mapping.IndexMapping
    if err := json.Unmarshal(mappingBytes, &indexMapping); err != nil {
        panic(err)
    }

    // Create a new index with the mapping
    index, err := bleve.New("example.bleve", &indexMapping)
    if err != nil {
        panic(err)
    }
    defer index.Close()
}
```

### Using Command Line

You can also create an index using the `bleve` command-line tool:

```bash
# Create new index with mapping file
bleve create -mapping mapping.json myindex.bleve
```

### Example JSON Mapping File

Here's an example `mapping.json` file:

```json
{
    "types": {
        "article": {
            "properties": {
                "title": {
                    "fields": [
                        {
                            "name": "title",
                            "type": "text",
                            "analyzer": "en",
                            "store": true,
                            "index": true
                        }
                    ]
                },
                "content": {
                    "fields": [
                        {
                            "name": "content",
                            "type": "text",
                            "analyzer": "en",
                            "store": true,
                            "index": true,
                            "include_term_vectors": true
                        }
                    ]
                },
                "published": {
                    "fields": [
                        {
                            "name": "published",
                            "type": "datetime",
                            "store": true,
                            "index": true
                        }
                    ]
                }
            }
        }
    },
    "default_type": "article",
    "default_analyzer": "standard",
    "default_datetime_parser": "dateTimeOptional",
    "default_field": "_all"
}
```

### Converting YAML to JSON

If you're using the YAML DSL described earlier, you can convert it to JSON before creating the index:

```go
package main

import (
    "encoding/json"
    "os"

    "github.com/blevesearch/bleve/v2"
    "github.com/blevesearch/bleve/v2/mapping"
    "gopkg.in/yaml.v3"
)

func main() {
    // Read YAML mapping file
    yamlBytes, err := os.ReadFile("mapping.yaml")
    if err != nil {
        panic(err)
    }

    // Parse YAML to map
    var mappingMap map[string]interface{}
    if err := yaml.Unmarshal(yamlBytes, &mappingMap); err != nil {
        panic(err)
    }

    // Convert map to JSON
    jsonBytes, err := json.Marshal(mappingMap)
    if err != nil {
        panic(err)
    }

    // Parse JSON into Bleve mapping
    var indexMapping mapping.IndexMapping
    if err := json.Unmarshal(jsonBytes, &indexMapping); err != nil {
        panic(err)
    }

    // Create index with mapping
    index, err := bleve.New("example.bleve", &indexMapping)
    if err != nil {
        panic(err)
    }
    defer index.Close()
}
```

### Opening Existing Index

To open an existing index:

```go
// Open existing index
index, err := bleve.Open("example.bleve")
if err != nil {
    panic(err)
}
defer index.Close()

// Get the mapping
mapping := index.Mapping()
```

### Validating JSON Mapping

It's a good practice to validate your JSON mapping before using it:

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/blevesearch/bleve/v2/mapping"
)

func validateMapping(filename string) error {
    // Read mapping file
    mappingBytes, err := os.ReadFile(filename)
    if err != nil {
        return fmt.Errorf("error reading mapping file: %v", err)
    }

    // Try to parse into IndexMapping
    var indexMapping mapping.IndexMapping
    if err := json.Unmarshal(mappingBytes, &indexMapping); err != nil {
        return fmt.Errorf("invalid mapping JSON: %v", err)
    }

    // Validate the mapping
    if err := indexMapping.Validate(); err != nil {
        return fmt.Errorf("invalid mapping configuration: %v", err)
    }

    return nil
}

func main() {
    if err := validateMapping("mapping.json"); err != nil {
        fmt.Printf("Mapping validation failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Mapping is valid")
}
```

### Best Practices

1. Always validate your mapping file before creating an index
2. Use constants for field names to avoid typos
3. Consider using the YAML DSL for better readability and maintenance
4. Keep mapping files under version control
5. Document any custom analyzers or special field configurations
6. Test your mapping with sample documents before production use

## Working with Documents

This section demonstrates how to load documents from JSON, validate them against mappings, and index them in Bleve.

### Loading and Indexing JSON Documents

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/blevesearch/bleve/v2"
    "github.com/blevesearch/bleve/v2/mapping"
)

// Document struct matching the mapping
type Article struct {
    Type      string     `json:"type"`
    Title     string     `json:"title"`
    Content   string     `json:"content"`
    Published string     `json:"published"`
    Tags      []string   `json:"tags"`
    Rating    float64    `json:"rating"`
    Embedding []float32  `json:"embedding"`
}

func main() {
    // Open or create index
    index, err := bleve.Open("example.bleve")
    if err == bleve.ErrorIndexPathDoesNotExist {
        // Create new index with mapping
        mapping := buildIndexMapping()
        index, err = bleve.New("example.bleve", mapping)
        if err != nil {
            panic(err)
        }
    } else if err != nil {
        panic(err)
    }
    defer index.Close()

    // Read JSON document
    docBytes, err := os.ReadFile("document.json")
    if err != nil {
        panic(err)
    }

    // Parse JSON into struct
    var article Article
    if err := json.Unmarshal(docBytes, &article); err != nil {
        panic(err)
    }

    // Index the document
    if err := index.Index("article1", article); err != nil {
        panic(err)
    }
}
```

### Example JSON Document

```json
{
    "type": "article",
    "title": "Introduction to Bleve",
    "content": "Bleve is a modern text indexing library for Go...",
    "published": "2024-01-15T10:30:00Z",
    "tags": ["search", "golang", "indexing"],
    "rating": 4.5,
    "embedding": [0.1, 0.2, 0.3, 0.4]
}
```

### Batch Processing Multiple Documents

```go
func indexBatchDocuments(index bleve.Index, documentsDir string) error {
    // Create a new batch
    batch := index.NewBatch()
    batchSize := 100
    count := 0

    // Walk through JSON files
    entries, err := os.ReadDir(documentsDir)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
            // Read document
            docBytes, err := os.ReadFile(filepath.Join(documentsDir, entry.Name()))
            if err != nil {
                return err
            }

            // Parse JSON
            var doc map[string]interface{}
            if err := json.Unmarshal(docBytes, &doc); err != nil {
                return err
            }

            // Add to batch
            if err := batch.Index(entry.Name(), doc); err != nil {
                return err
            }
            count++

            // Execute batch when size is reached
            if count >= batchSize {
                if err := index.Batch(batch); err != nil {
                    return err
                }
                batch = index.NewBatch()
                count = 0
            }
        }
    }

    // Index remaining documents
    if count > 0 {
        if err := index.Batch(batch); err != nil {
            return err
        }
    }

    return nil
}
```

### Document Validation

Here's a utility to validate documents against a mapping:

```go
func validateDocument(doc interface{}, indexMapping mapping.IndexMapping) error {
    // Convert document to map if it's not already
    var docMap map[string]interface{}
    
    switch v := doc.(type) {
    case map[string]interface{}:
        docMap = v
    default:
        // Convert struct to map using JSON marshaling
        jsonBytes, err := json.Marshal(doc)
        if err != nil {
            return fmt.Errorf("error marshaling document: %v", err)
        }
        
        if err := json.Unmarshal(jsonBytes, &docMap); err != nil {
            return fmt.Errorf("error unmarshaling document to map: %v", err)
        }
    }

    // Get document type
    docType := indexMapping.DefaultType
    if typeVal, ok := docMap["type"].(string); ok {
        docType = typeVal
    }

    // Get mapping for this document type
    docMapping := indexMapping.DocumentMapping(docType)
    if docMapping == nil {
        return fmt.Errorf("no mapping found for document type: %s", docType)
    }

    // Validate required fields
    for fieldName, fieldMapping := range docMapping.Properties {
        if !fieldMapping.Dynamic && fieldMapping.Enabled {
            if _, exists := docMap[fieldName]; !exists {
                return fmt.Errorf("required field missing: %s", fieldName)
            }
        }
    }

    // Validate field types
    var validateField func(string, interface{}, *mapping.DocumentMapping) error
    validateField = func(path string, value interface{}, mapping *mapping.DocumentMapping) error {
        if mapping == nil {
            return nil
        }

        switch v := value.(type) {
        case map[string]interface{}:
            for key, val := range v {
                fieldPath := path
                if fieldPath != "" {
                    fieldPath += "."
                }
                fieldPath += key

                if subMapping, exists := mapping.Properties[key]; exists {
                    if err := validateField(fieldPath, val, subMapping); err != nil {
                        return err
                    }
                }
            }
        case []interface{}:
            // Validate array elements if needed
            return nil
        default:
            // Validate field type
            if len(mapping.Fields) > 0 {
                field := mapping.Fields[0]
                switch field.Type {
                case "text", "keyword":
                    if _, ok := value.(string); !ok {
                        return fmt.Errorf("field %s should be string, got %T", path, value)
                    }
                case "number":
                    switch value.(type) {
                    case float64, int, int64, float32:
                        // Valid numeric types
                    default:
                        return fmt.Errorf("field %s should be numeric, got %T", path, value)
                    }
                case "datetime":
                    if _, ok := value.(string); !ok {
                        return fmt.Errorf("field %s should be datetime string, got %T", path, value)
                    }
                    // Could add datetime parsing validation here
                case "boolean":
                    if _, ok := value.(bool); !ok {
                        return fmt.Errorf("field %s should be boolean, got %T", path, value)
                    }
                case "vector":
                    if arr, ok := value.([]float32); !ok {
                        return fmt.Errorf("field %s should be []float32, got %T", path, value)
                    } else if len(arr) != field.Dims {
                        return fmt.Errorf("field %s vector dimension mismatch: expected %d, got %d", 
                            path, field.Dims, len(arr))
                    }
                }
            }
        }
        return nil
    }

    return validateField("", docMap, docMapping)
}

// Example usage
func main() {
    // Open index
    index, err := bleve.Open("example.bleve")
    if err != nil {
        panic(err)
    }
    defer index.Close()

    // Read document
    docBytes, err := os.ReadFile("document.json")
    if err != nil {
        panic(err)
    }

    var doc map[string]interface{}
    if err := json.Unmarshal(docBytes, &doc); err != nil {
        panic(err)
    }

    // Validate document
    if err := validateDocument(doc, index.Mapping()); err != nil {
        fmt.Printf("Document validation failed: %v\n", err)
        return
    }

    // Index document
    if err := index.Index("doc1", doc); err != nil {
        panic(err)
    }
}
```

### Best Practices for Document Processing

1. **Batch Processing**
   - Use batch operations for indexing multiple documents
   - Choose appropriate batch sizes (100-1000 documents)
   - Monitor memory usage during batch operations

2. **Validation**
   - Validate documents before indexing
   - Check required fields
   - Verify field types match mapping
   - Validate vector dimensions for vector fields

3. **Error Handling**
   - Handle JSON parsing errors gracefully
   - Log validation failures with details
   - Implement retry logic for failed documents
   - Keep track of failed documents for later processing

4. **Performance**
   - Use goroutines for parallel processing
   - Monitor index size and performance
   - Consider using bulk indexing for large datasets
   - Implement rate limiting if needed

5. **Data Quality**
   - Clean and normalize data before indexing
   - Handle missing or null values appropriately
   - Convert field values to correct types
   - Validate date formats

6. **Monitoring**
   - Log indexing statistics
   - Track processing time
   - Monitor error rates
   - Implement health checks

## Retrieving Index Mappings

You can retrieve mappings from an existing Bleve index in several ways:

### Basic Mapping Retrieval

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/blevesearch/bleve/v2"
)

func main() {
    // Open existing index
    index, err := bleve.Open("example.bleve")
    if err != nil {
        panic(err)
    }
    defer index.Close()

    // Get the mapping
    mapping := index.Mapping()

    // Print mapping details
    fmt.Printf("Default Type: %s\n", mapping.DefaultType)
    fmt.Printf("Default Analyzer: %s\n", mapping.DefaultAnalyzer)
    fmt.Printf("Default Field: %s\n", mapping.DefaultField)
}
```

### Export Mapping to JSON

```go
func exportMapping(indexPath string, outputPath string) error {
    // Open index
    index, err := bleve.Open(indexPath)
    if err != nil {
        return fmt.Errorf("error opening index: %v", err)
    }
    defer index.Close()

    // Get mapping
    mapping := index.Mapping()

    // Marshal to JSON with indentation
    jsonBytes, err := json.MarshalIndent(mapping, "", "    ")
    if err != nil {
        return fmt.Errorf("error marshaling mapping: %v", err)
    }

    // Write to file
    if err := os.WriteFile(outputPath, jsonBytes, 0644); err != nil {
        return fmt.Errorf("error writing mapping file: %v", err)
    }

    return nil
}
```

### Inspect Specific Document Type Mapping

```go
func inspectDocumentMapping(index bleve.Index, docType string) error {
    mapping := index.Mapping()
    
    // Get document mapping for specific type
    docMapping := mapping.DocumentMapping(docType)
    if docMapping == nil {
        return fmt.Errorf("no mapping found for document type: %s", docType)
    }

    // Print document mapping details
    fmt.Printf("Document Type: %s\n", docType)
    fmt.Printf("Enabled: %v\n", docMapping.Enabled)
    fmt.Printf("Dynamic: %v\n", docMapping.Dynamic)
    
    // Print field mappings
    for fieldName, fieldMapping := range docMapping.Properties {
        fmt.Printf("\nField: %s\n", fieldName)
        for _, field := range fieldMapping.Fields {
            fmt.Printf("  Type: %s\n", field.Type)
            fmt.Printf("  Analyzer: %s\n", field.Analyzer)
            fmt.Printf("  Store: %v\n", field.Store)
            fmt.Printf("  Index: %v\n", field.Index)
        }
    }

    return nil
}
```

### Analyze Mapping Configuration

```go
func analyzeMappingConfig(index bleve.Index) {
    mapping := index.Mapping()

    // Check global settings
    fmt.Println("Global Settings:")
    fmt.Printf("Default Type: %s\n", mapping.DefaultType)
    fmt.Printf("Default Analyzer: %s\n", mapping.DefaultAnalyzer)
    fmt.Printf("Default Field: %s\n", mapping.DefaultField)
    fmt.Printf("Default DateTime Parser: %s\n", mapping.DefaultDateTimeParser)
    fmt.Printf("Byte Array Converter: %s\n", mapping.ByteArrayConverter)

    // List all document types
    fmt.Println("\nDocument Types:")
    for typeName := range mapping.TypeMapping {
        docMapping := mapping.DocumentMapping(typeName)
        fmt.Printf("\nType: %s\n", typeName)
        fmt.Printf("  Enabled: %v\n", docMapping.Enabled)
        fmt.Printf("  Dynamic: %v\n", docMapping.Dynamic)
        
        // List fields
        fmt.Println("  Fields:")
        for fieldName, fieldMapping := range docMapping.Properties {
            fmt.Printf("    %s:\n", fieldName)
            for _, field := range fieldMapping.Fields {
                fmt.Printf("      Type: %s\n", field.Type)
                if field.Type == "vector" {
                    fmt.Printf("      Dims: %d\n", field.Dims)
                    fmt.Printf("      Similarity: %s\n", field.Similarity)
                }
            }
        }
    }

    // List custom analyzers
    fmt.Println("\nCustom Analyzers:")
    for name := range mapping.CustomAnalysis.Analyzers {
        fmt.Printf("  - %s\n", name)
    }
}
```

### Complete Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/blevesearch/bleve/v2"
)

func main() {
    // Open index
    index, err := bleve.Open("example.bleve")
    if err != nil {
        panic(err)
    }
    defer index.Close()

    // Get and analyze mapping
    mapping := index.Mapping()

    // 1. Print basic info
    fmt.Println("Basic Index Information:")
    fmt.Printf("Default Type: %s\n", mapping.DefaultType)
    fmt.Printf("Default Analyzer: %s\n", mapping.DefaultAnalyzer)

    // 2. Export mapping to JSON
    jsonBytes, err := json.MarshalIndent(mapping, "", "    ")
    if err != nil {
        panic(err)
    }
    fmt.Println("\nFull Mapping JSON:")
    fmt.Println(string(jsonBytes))

    // 3. Analyze specific document type
    docType := "article"  // replace with your document type
    if err := inspectDocumentMapping(index, docType); err != nil {
        fmt.Printf("Error inspecting document mapping: %v\n", err)
    }

    // 4. Analyze full configuration
    fmt.Println("\nFull Configuration Analysis:")
    analyzeMappingConfig(index)
}

// Save mapping to file
func saveMappingToFile(mapping interface{}, filename string) error {
    jsonBytes, err := json.MarshalIndent(mapping, "", "    ")
    if err != nil {
        return err
    }
    return os.WriteFile(filename, jsonBytes, 0644)
}
```

### Using Command Line

You can also use the `bleve` command-line tool to inspect mappings:

```bash
# Show mapping for an index
bleve mapping example.bleve

# Export mapping to JSON file
bleve mapping example.bleve > mapping.json
```

### Best Practices

1. **Version Control**
   - Save exported mappings in version control
   - Document mapping changes
   - Keep a history of mapping evolution

2. **Documentation**
   - Document custom analyzers
   - Document field usage and constraints
   - Keep mapping documentation up to date

3. **Validation**
   - Validate mapping changes before applying
   - Test with sample documents
   - Check for backward compatibility

4. **Maintenance**
   - Regularly review mappings
   - Clean up unused field mappings
   - Monitor mapping size and complexity

5. **Security**
   - Secure mapping files
   - Control access to mapping information
   - Sanitize mapping output in logs