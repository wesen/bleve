# Bleve Query DSL Specification

This document provides a comprehensive specification for a YAML-based Domain Specific Language (DSL) designed for constructing search queries against a Bleve index. It is intended for programmers who will implement the DSL parser and for users who will write queries using the DSL.

## 1. Introduction

The Bleve Query DSL aims to provide a human-readable, expressive, and composable way to define search queries for Bleve, a modern text indexing and search library written in Go. This DSL takes inspiration from Elasticsearch's query language while being tailored to Bleve's specific features and capabilities.

### 1.1 Design Goals

*   **Clarity:** The DSL should be easy to understand and write, even for users not intimately familiar with Bleve's internal workings.
*   **Expressiveness:** It should be capable of representing the full range of Bleve's query types and search options.
*   **Composability:** Queries should be composable, allowing complex searches to be built from simpler parts.
*   **Extensibility:** The DSL should be designed to accommodate future additions to Bleve's search capabilities.
*   **Validation:** The DSL should facilitate early error detection by allowing for schema validation and query verification.

## 2. Structure

A query in the DSL is represented as a YAML document with the following top-level structure:

````yaml
query:
  # Query clauses (see Section 3)
options:
  # Search options (optional, see Section 4)
facets:
  # Facet definitions (optional, see Section 5)
````

### 2.1 Components

*   **`query`:** (Required) This section defines the core search criteria using various query clauses.
*   **`options`:** (Optional) This section allows you to specify search parameters like result size, sorting, highlighting, etc.
*   **`facets`:** (Optional) This section defines aggregations to be performed on the search results, such as grouping by terms or ranges.

## 3. Query Clauses

The `query` section is the heart of the DSL. It defines what you are searching for using a combination of different query types. Each query type targets a specific kind of search operation.

### 3.1 Match Query

Performs a full-text search, analyzing the query text using the field's analyzer (if configured).

````yaml
query:
  match:
    field: content        # Name of the field to search
    value: golang programming # Text to search for
    boost: 1.0            # Optional: Boost factor for this query (default: 1.0)
    operator: or          # Optional: Operator to use between terms (or/and, default: or)
    fuzziness: 0          # Optional: Fuzziness (edit distance, default: 0)
    prefix_length: 0      # Optional: Number of leading characters that must match exactly (default: 0)
    analyzer: ""           # Optional: Specific analyzer to use (overrides field's default)
````

**Fields:**

*   **`field`:** (Required, string) The name of the field to search in.
*   **`value`:** (Required, string) The text to search for.
*   **`boost`:** (Optional, float) A positive floating-point number that increases the query's score, making it more relevant.
*   **`operator`:** (Optional, string) Either "or" or "and". Determines whether all terms must match ("and") or at least one term must match ("or").
*   **`fuzziness`:** (Optional, integer) The maximum edit distance (number of character changes) allowed for a term to be considered a match.
*   **`prefix_length`:** (Optional, integer) The number of initial characters that must match exactly, even when using fuzziness.
*   **`analyzer`:** (Optional, string) The name of a specific analyzer to use for this query, overriding the field's default analyzer.

### 3.2 Match Phrase Query

Searches for an exact phrase, where all terms must appear in the specified order and proximity.

````yaml
query:
  match_phrase:
    field: content
    value: "quick brown fox"
    boost: 1.0
    slop: 0        # Optional: Allowed distance between terms (default: 0)
    analyzer: ""
````

**Fields:**

*   **`field`:** (Required, string) The field to search.
*   **`value`:** (Required, string) The exact phrase to search for.
*   **`boost`:** (Optional, float) Boost factor.
*   **`slop`:** (Optional, integer) The maximum number of positions allowed between the terms in the phrase. A slop of 1 means terms can be one position apart.
*   **`analyzer`:** (Optional, string) Analyzer to use, overriding the field's default.

### 3.3 Term Query

Searches for an exact term without any text analysis. Useful for searching fields that are not analyzed (e.g., keywords, IDs).

````yaml
query:
  term:
    field: tags
    value: golang
    boost: 1.0
````

**Fields:**

*   **`field`:** (Required, string) The field to search.
*   **`value`:** (Required, string) The exact term to search for.
*   **`boost`:** (Optional, float) Boost factor.

### 3.4 Prefix Query

Searches for terms that start with a specific prefix.

````yaml
query:
  prefix:
    field: title
    value: prog
    boost: 1.0
````

**Fields:**

*   **`field`:** (Required, string) The field to search.
*   **`value`:** (Required, string) The prefix to search for.
*   **`boost`:** (Optional, float) Boost factor.

### 3.5 Fuzzy Query

Searches for terms that are similar to the given term, allowing for a certain number of character edits (insertions, deletions, substitutions).

````yaml
query:
  fuzzy:
    field: content
    value: programmer
    boost: 1.0
    fuzziness: 2  # Maximum edit distance
    prefix_length: 0
````

**Fields:**

*   **`field`:** (Required, string) The field to search.
*   **`value`:** (Required, string) The term to search for.
*   **`boost`:** (Optional, float) Boost factor.
*   **`fuzziness`:** (Optional, integer) The maximum edit distance allowed.
*   **`prefix_length`:** (Optional, integer) The number of leading characters that must match exactly.

### 3.6 Wildcard Query

Searches for terms that match a wildcard pattern, where `*` matches any sequence of characters and `?` matches any single character.

````yaml
query:
  wildcard:
    field: content
    value: pro*mer
    boost: 1.0
````

**Fields:**

*   **`field`:** (Required, string) The field to search.
*   **`value`:** (Required, string) The wildcard pattern to search for.
*   **`boost`:** (Optional, float) Boost factor.

### 3.7 Regexp Query

Searches for terms that match a regular expression.

````yaml
query:
  regexp:
    field: content
    value: pro.*mer
    boost: 1.0
````

**Fields:**

*   **`field`:** (Required, string) The field to search.
*   **`value`:** (Required, string) The regular expression to search for.
*   **`boost`:** (Optional, float) Boost factor.

### 3.8 Numeric Range Query

Searches for numeric values that fall within a specified range.

````yaml
query:
  numeric_range:
    field: price
    min: 100
    max: 200
    inclusive_min: true
    inclusive_max: false
    boost: 1.0
````

**Fields:**

*   **`field`:** (Required, string) The numeric field to search.
*   **`min`:** (Optional, float) The minimum value of the range.
*   **`max`:** (Optional, float) The maximum value of the range.
*   **`inclusive_min`:** (Optional, boolean) Whether the minimum value is included in the range (default: true).
*   **`inclusive_max`:** (Optional, boolean) Whether the maximum value is included in the range (default: true).
*   **`boost`:** (Optional, float) Boost factor.

### 3.9 Date Range Query

Searches for dates that fall within a specified range.

````yaml
query:
  date_range:
    field: created_at
    start: "2024-01-01T00:00:00Z"
    end: "2024-12-31T23:59:59Z"
    inclusive_start: true
    inclusive_end: true
    boost: 1.0
````

**Fields:**

*   **`field`:** (Required, string) The date field to search.
*   **`start`:** (Optional, string) The start date of the range (ISO 8601 format).
*   **`end`:** (Optional, string) The end date of the range (ISO 8601 format).
*   **`inclusive_start`:** (Optional, boolean) Whether the start date is included in the range (default: true).
*   **`inclusive_end`:** (Optional, boolean) Whether the end date is included in the range (default: true).
*   **`boost`:** (Optional, float) Boost factor.

### 3.10 Boolean Query

Combines multiple query clauses using boolean logic (AND, OR, NOT).

````yaml
query:
  bool:
    must:          # All of these queries must match (logical AND)
      - match:
          field: title
          value: golang
      - numeric_range:
          field: price
          min: 10
          max: 100
    should:        # At least one of these queries should match (logical OR)
      - match:
          field: description
          value: golang
      - match:
          field: tags
          value: programming
    must_not:     # None of these queries must match (logical NOT)
      - term:
          field: status
          value: draft
    minimum_should_match: 1 # Optional: Minimum number of "should" clauses that must match (default: 0)
    boost: 1.0
````

**Fields:**

*   **`must`:** (Optional, array of queries) A list of queries that **must** all match.
*   **`should`:** (Optional, array of queries) A list of queries where at least one **should** match.
*   **`must_not`:** (Optional, array of queries) A list of queries that **must not** match.
*   **`minimum_should_match`:** (Optional, integer) The minimum number of `should` queries that must be satisfied.
*   **`boost`:** (Optional, float) Boost factor for the entire boolean query.

### 3.11 Query String Query

Allows you to use Bleve's query string syntax, which is similar to the Lucene query language.

````yaml
query:
  query_string:
    query: "title:golang AND (tags:programming OR tags:database) NOT status:draft"
    default_field: content # Optional: The default field to search if no field is specified in the query string
    boost: 1.0
````

**Fields:**

*   **`query`:** (Required, string) The query string.
*   **`default_field`:** (Optional, string) The default field to use when no field is explicitly specified in the query string.
*   **`boost`:** (Optional, float) Boost factor.

### 3.12 Doc ID Query

Searches for documents with specific IDs.

````yaml
query:
  doc_id:
    ids:
      - doc1
      - doc2
      - doc3
````

**Fields:**

*   **`ids`:** (Required, array of strings) A list of document IDs to search for.

### 3.13 Exists Query

Checks if a field exists in a document (i.e., the field has a value).

````yaml
query:
  exists:
    field: optional_field
````

**Fields:**

*   **`field`:** (Required, string) The name of the field to check for existence.

### 3.14 Vector Search

Performs a vector similarity search, finding documents whose vectors are closest to a given query vector.

````yaml
query:
  vector:
    field: embedding  # The field containing the vector data
    # Option 1: Text input
    text: "query string"
    model: all-minilm  # Model to use for generating embedding (required with text)
    # Option 2: Precomputed vector
    vector: [0.1, 0.2, 0.3, 0.4] # The vector embedding itself
    model: all-minilm # Model used to generate the vector (required with vector)
    k: 10             # Number of nearest neighbors to retrieve
    boost: 1.0        # Boost for this query
````

**Fields:**

*   **`field`:** (Required, string) The name of the field that contains the vector embeddings.
*   **`text`:** (Optional, string) The query text. If provided, the `model` will be used to generate the embedding vector for this text. **Mutually exclusive with `vector`.**
*   **`vector`:** (Optional, array of floats) The precomputed vector embedding to search for. **Mutually exclusive with `text`.**
*   **`model`:** (Required, string)
    *   If `text` is used: The name of the embedding model to use for generating the vector from the `text`.
    *   If `vector` is used: The name of the model that was used to generate the provided `vector`. This is used for validation to ensure compatibility with the indexed vectors.
*   **`k`:** (Required, integer) The number of nearest neighbors (most similar documents) to retrieve.
*   **`boost`:** (Optional, float) Boost factor for this query.

**Important Notes on Vector Search:**

*   **Model Compatibility:** When using a precomputed `vector`, the parser must verify that the specified `model` matches the model used to create the vector field in the index mapping. If there's a mismatch, the parser should return an error.
*   **Embedding Generation:** If `text` is provided, the parser must use the specified `model` to generate the vector embedding before constructing the Bleve `search.VectorQuery`.

## 4. Search Options

The `options` section allows you to control various aspects of the search execution and results.

````yaml
options:
  size: 10       # Number of results to return per page (default: 10)
  from: 0        # Offset of the first result (for pagination, default: 0)
  explain: false # Include detailed explanation of score calculation (default: false)
  fields:        # List of fields to return in the results (default: all fields)
    - title
    - created_at
  sort:          # Sorting order (default: by score, descending)
    - field: created_at
      desc: true
    - field: _score  # Special field representing the document's score
      desc: true
  highlight:     # Highlight matching terms in the results
    style: html  # Highlighting style (html or ansi, default: none)
    fields:      # Fields to highlight
      - content
````

**Fields:**

*   **`size`:** (Optional, integer) The number of search results to return.
*   **`from`:** (Optional, integer) The starting offset of the results (used for pagination).
*   **`explain`:** (Optional, boolean) Whether to include a detailed explanation of the score calculation for each result.
*   **`fields`:** (Optional, array of strings) The list of fields to retrieve and include in the search results. If not specified, all fields are returned. You can use `"*"` to explicitly include all fields or `"_id"` to include only the document ID.
*   **`sort`:** (Optional, array of sort objects) Specifies the sorting order of the results. Each sort object has:
    *   **`field`:** (string) The field to sort by. Use `"_score"` to sort by relevance score and `"_id"` to sort by document ID.
    *   **`desc`:** (boolean) Whether to sort in descending order (default: false).
*   **`highlight`:** (Optional, highlight object) Specifies how to highlight matching terms in the results.
    *   **`style`:** (string) The highlighting style: "html" (wraps matches in HTML tags) or "ansi" (uses ANSI escape codes).
    *   **`fields`:** (array of strings) The fields to highlight.

## 5. Facets

Facets provide aggregated data based on the search results, allowing users to drill down and refine their searches.

````yaml
facets:
  tags:             # Facet name (can be anything)
    type: terms     # Facet type (terms, numeric, or date)
    field: tags     # Field to facet on
    size: 10        # Maximum number of facet terms to return
  price_ranges:
    type: numeric
    field: price
    ranges:
      - name: low     # Range name
        min: 0
        max: 100
      - name: medium
        min: 100
        max: 500
      - name: high
        min: 500
  date_ranges:
    type: date
    field: created_at
    ranges:
      - name: last_24h
        start: now-24h # Relative time expressions are allowed
        end: now
      - name: last_week
        start: now-7d
        end: now-24h
      - name: last_month
        start: now-30d
        end: now-7d
````

**Fields:**

*   **`facet_name`:** (Required, string) The name you give to the facet (e.g., "tags", "price\_ranges"). This name will be used in the search results to identify the facet's data.
*   **`type`:** (Required, string) The type of facet:
    *   **`terms`:** Groups results by unique terms in a field.
    *   **`numeric`:** Groups results into numeric ranges.
    *   **`date`:** Groups results into date ranges.
*   **`field`:** (Required, string) The field to perform faceting on.
*   **`size`:** (Required for `terms` facet, integer) The maximum number of terms to return for a terms facet.
*   **`ranges`:** (Required for `numeric` and `date` facets, array of range objects) Defines the ranges for numeric and date facets. Each range object has:
    *   **`name`:** (string) The name of the range (e.g., "low", "last\_24h").
    *   **`min`:** (float, for `numeric` facet) The minimum value of the range.
    *   **`max`:** (float, for `numeric` facet) The maximum value of the range.
    *   **`start`:** (string, for `date` facet) The start date/time of the range (ISO 8601 format or relative expressions like "now-1d", "now-7d").
    *   **`end`:** (string, for `date` facet) The end date/time of the range.

**Relative Date Expressions:**

*   **`now`:** Represents the current time.
*   **`now-{N}d`:** {N} days before now (e.g., `now-1d`, `now-7d`, `now-30d`).
*   **`now-{N}h`:** {N} hours before now (e.g., `now-1h`, `now-24h`).
*   **`now-{N}m`:** {N} minutes before now.
*   **`now-{N}s`:** {N} seconds before now.

## 6. Examples

### 6.1. Basic Text Search with Highlighting

````yaml
query:
  match:
    field: content
    value: golang programming
options:
  highlight:
    style: html
    fields:
      - content
````

### 6.2. Boolean Query with Range and Facets

````yaml
query:
  bool:
    must:
      - match:
          field: title
          value: "Bleve Search"
      - date_range:
          field: created_at
          start: "2023-01-01T00:00:00Z"
          end: "2024-01-01T00:00:00Z"
    should:
      - match:
          field: tags
          value: tutorial
          boost: 2.0
    minimum_should_match: 0
options:
  size: 20
  sort:
    - field: created_at
      desc: true
facets:
  tags:
    type: terms
    field: tags
    size: 5
  price_ranges:
    type: numeric
    field: price
    ranges:
      - name: budget
        min: 0
        max: 50
      - name: mid-range
        min: 50
        max: 200
````

### 6.3. Vector Search with Text Input

````yaml
query:
  vector:
    field: embedding
    text: "What is the best programming language?"
    model: all-minilm
    k: 5
````

### 6.4. Vector Search with Precomputed Vector

````yaml
query:
  vector:
    field: embedding
    vector: [0.85, 0.22, 0.45, 0.67, 0.11]
    model: all-minilm
    k: 10
options:
  fields:
    - title
    - category
````

### 6.5. Combined Vector and Keyword Search

````yaml
query:
  bool:
    must:
      - vector:
          field: embedding
          text: "modern web development"
          model: all-minilm
          k: 50
          boost: 2.0
      - match:
          field: content
          value: "golang"
    should:
      - term:
          field: category
          value: "technology"
          boost: 0.5
    minimum_should_match: 0
options:
  size: 20
  sort:
    - field: _score
      desc: true
````

## 7. Implementation Notes

### 7.1. Parser

*   The parser should be implemented to translate the YAML DSL into Bleve's Go API calls, constructing the appropriate `search.Query` objects and setting the search parameters.
*   The parser should handle type checking and validation, ensuring that the provided values match the expected types and constraints defined in this specification.
*   The parser should provide informative error messages when encountering invalid DSL syntax, unsupported query combinations, or type mismatches.

### 7.2. Error Handling

*   **Invalid YAML:** Handle YAML parsing errors gracefully.
*   **Unknown Fields:** Implement options for strictness (e.g., reject unknown fields or query types) or allow unknown fields to be ignored.
*   **Type Mismatches:** Validate field types and values against the expected types.
*   **Missing Required Fields:** Ensure that all required fields for each query type are present.
*   **Mutually Exclusive Fields:** Enforce mutual exclusivity where applicable (e.g., `text` and `vector` in the `vector` query).
*   **Model Validation (Vector Search):** Verify that the `model` specified in a `vector` query with a precomputed vector matches the model used to create the vector field in the index mapping.
*   **Range Errors:** Check for invalid numeric or date ranges (e.g., `min` greater than `max`).
*   **Unsupported Combinations:** Handle cases where certain query types or options cannot be combined.

### 7.3. Extensibility

*   The DSL should be designed in a way that allows for easy addition of new query types, search options, and facet types in the future.
*   Consider using a modular design for the parser to facilitate adding new features without disrupting existing functionality.

### 7.4. Testing

*   Thoroughly test the parser and the generated Bleve queries with a wide range of valid and invalid DSL inputs.
*   Use unit tests to verify individual components of the parser and integration tests to ensure that the entire query generation process works correctly.
*   Test with various Bleve index configurations and data sets.

## 8. Future Considerations

*   **Geo Queries:** Add support for geospatial queries (e.g., searching within a certain distance of a point or within a bounding box).
*   **More Like This (MLT) Queries:** Implement queries that find documents similar to a given document.
*   **Scripting:** Allow for more complex logic using scripting within queries (e.g., custom scoring functions).
*   **Nested Documents:** Support querying nested document structures.
*   **Advanced Aggregations:** Expand the faceting capabilities with more advanced aggregation types (e.g., histograms, percentiles).

This specification provides a comprehensive guide for implementing and using the Bleve Query DSL. It should serve as a living document that evolves alongside Bleve's capabilities.