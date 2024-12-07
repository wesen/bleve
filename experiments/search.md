# Bleve Search Guide

This document provides a comprehensive guide to performing searches with Bleve, focusing on text, numeric, and date queries.

## Basic Text Search

Bleve provides several ways to perform text searches, from simple queries to more complex matching options.

### Simple Text Query

The most basic way to search is using a simple text query:

```go
package main

import (
    "fmt"
    "github.com/blevesearch/bleve/v2"
)

func main() {
    // Open index
    index, err := bleve.Open("example.bleve")
    if err != nil {
        panic(err)
    }
    defer index.Close()

    // Create a simple query
    query := bleve.NewQueryStringQuery("search text")
    
    // Perform search
    searchRequest := bleve.NewSearchRequest(query)
    searchResult, err := index.Search(searchRequest)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found %d matches\n", searchResult.Total)
    for _, hit := range searchResult.Hits {
        fmt.Printf("Score: %f, ID: %s\n", hit.Score, hit.ID)
    }
}
```

### Match Query

Match queries analyze the search text using the field's analyzer:

```go
// Match query on specific field
matchQuery := bleve.NewMatchQuery("golang programming")
matchQuery.SetField("content")

// Search with options
searchRequest := bleve.NewSearchRequest(matchQuery)
searchRequest.Size = 10
searchRequest.From = 0
```

### Match Phrase Query

Match phrase queries look for exact phrases:

```go
// Match exact phrase
phraseQuery := bleve.NewMatchPhraseQuery("quick brown fox")
phraseQuery.SetField("content")

// With slop (words can be this many positions apart)
phraseQuery.SetSlop(1)
```

### Prefix Query

Search for terms that begin with a prefix:

```go
// Prefix query
prefixQuery := bleve.NewPrefixQuery("prog")
prefixQuery.SetField("title")

// Search
searchRequest := bleve.NewSearchRequest(prefixQuery)
```

### Fuzzy Query

Fuzzy queries allow for character variations:

```go
// Fuzzy query with max edit distance
fuzzyQuery := bleve.NewFuzzyQuery("programmer")
fuzzyQuery.SetField("content")
fuzzyQuery.SetFuzziness(2) // Allow 2 character edits

// Search
searchRequest := bleve.NewSearchRequest(fuzzyQuery)
```

### Regular Expression Query

Search using regular expressions:

```go
// Regex query
regexpQuery := bleve.NewRegexpQuery("pro.*mer")
regexpQuery.SetField("content")

// Search
searchRequest := bleve.NewSearchRequest(regexpQuery)
```

### Wildcard Query

Search using wildcard patterns:

```go
// Wildcard query
wildcardQuery := bleve.NewWildcardQuery("pro*mer")
wildcardQuery.SetField("content")

// Search
searchRequest := bleve.NewSearchRequest(wildcardQuery)
```

## Numeric & Date Searches

Bleve provides powerful range queries for numeric and date fields.

### Numeric Range Query

Search for numeric values within a range:

```go
// Numeric range query
minValue := float64(100)
maxValue := float64(200)
rangeQuery := bleve.NewNumericRangeQuery(&minValue, &maxValue)
rangeQuery.SetField("price")

// Include/exclude boundaries
rangeQuery.InclusiveMin = true
rangeQuery.InclusiveMax = false

// Search
searchRequest := bleve.NewSearchRequest(rangeQuery)
```

### Date Range Query

Search for dates within a range:

```go
// Date range query
startDate := "2024-01-01T00:00:00Z"
endDate := "2024-12-31T23:59:59Z"
dateQuery := bleve.NewDateRangeQuery(&startDate, &endDate)
dateQuery.SetField("created_at")

// Include/exclude boundaries
dateQuery.InclusiveStart = true
dateQuery.InclusiveEnd = true

// Search
searchRequest := bleve.NewSearchRequest(dateQuery)
```

## Search Options

### Pagination

Control the number of results and starting position:

```go
searchRequest := bleve.NewSearchRequest(query)
searchRequest.Size = 10   // Number of results per page
searchRequest.From = 20   // Start from result #20
```

### Field Boost

Boost the importance of certain fields:

```go
// Boost title field
matchQuery := bleve.NewMatchQuery("golang")
matchQuery.SetField("title")
matchQuery.SetBoost(2.0)

// Combine with content field
contentQuery := bleve.NewMatchQuery("golang")
contentQuery.SetField("content")

// Combine queries
shouldQuery := bleve.NewDisjunctionQuery(matchQuery, contentQuery)
```

### Result Highlighting

Enable highlighting of matching terms:

```go
searchRequest := bleve.NewSearchRequest(query)

// Configure highlighting
searchRequest.Highlight = bleve.NewHighlight()
searchRequest.Highlight.Fields = []string{"content"}
searchRequest.Highlight.Style = "html"  // or "ansi"

// Process results
for _, hit := range searchResult.Hits {
    if fragments, ok := hit.Fragments["content"]; ok {
        for _, fragment := range fragments {
            fmt.Printf("Fragment: %s\n", fragment)
        }
    }
}
```

### Field Selection

Control which fields are returned in results:

```go
searchRequest := bleve.NewSearchRequest(query)

// Only return specific fields
searchRequest.Fields = []string{"title", "created_at"}

// Access fields in results
for _, hit := range searchResult.Hits {
    title := hit.Fields["title"]
    createdAt := hit.Fields["created_at"]
    fmt.Printf("Title: %v, Created: %v\n", title, createdAt)
}
```

### Sorting

Sort results by field values:

```go
searchRequest := bleve.NewSearchRequest(query)

// Sort by field
searchRequest.SortBy([]string{"created_at"})

// Sort by multiple fields
searchRequest.SortBy([]string{"_score", "-created_at"}) // '-' for descending

// Custom sort
searchRequest.SortBy([]string{
    "_score",
    "-created_at",
    "title",
})
```

## Boolean Queries

Bleve provides powerful boolean query combinations using conjunction (AND), disjunction (OR), and negation (NOT).

### Boolean Combinations

```go
// Create individual queries
titleQuery := bleve.NewMatchQuery("golang")
titleQuery.SetField("title")

priceQuery := bleve.NewNumericRangeQuery(float64(10), float64(100))
priceQuery.SetField("price")

statusQuery := bleve.NewTermQuery("published")
statusQuery.SetField("status")

// AND query (all must match)
mustQuery := bleve.NewConjunctionQuery(titleQuery, priceQuery, statusQuery)

// OR query (any can match)
shouldQuery := bleve.NewDisjunctionQuery(titleQuery, priceQuery)

// NOT query (must not match)
notQuery := bleve.NewBooleanQuery()
notQuery.AddMustNot(statusQuery)
```

### Must, Should, Must Not Clauses

```go
// Create a boolean query with multiple clauses
boolQuery := bleve.NewBooleanQuery()

// Must match these (AND)
boolQuery.AddMust(titleQuery)
boolQuery.AddMust(priceQuery)

// Should match these (OR, boosts score)
boolQuery.AddShould(descriptionQuery)
boolQuery.AddShould(tagsQuery)

// Must not match these (NOT)
boolQuery.AddMustNot(statusQuery)
```

### Minimum Should Match

```go
// Create disjunction with minimum matches
shouldQuery := bleve.NewDisjunctionQuery(
    titleQuery,
    contentQuery,
    tagsQuery,
)
shouldQuery.SetMin(2) // At least 2 should match
```

## Advanced Query Types

### Term Query

Exact term matching without analysis:

```go
// Term query (exact match)
termQuery := bleve.NewTermQuery("golang")
termQuery.SetField("tags")
termQuery.SetBoost(2.0)
```

### Phrase Query

Match exact phrases with term positions:

```go
// Phrase query
phraseQuery := bleve.NewPhraseQuery([]string{"quick", "brown", "fox"})
phraseQuery.SetField("content")
```

### Query String Query

Lucene-style query syntax:

```go
// Query string with Lucene syntax
queryString := `title:golang AND (tags:programming OR tags:database) NOT status:draft`
query := bleve.NewQueryStringQuery(queryString)

// With default field
query := bleve.NewQueryStringQuery("golang programming")
query.SetDefaultField("content")
```

### Document ID Query

Search for specific document IDs:

```go
// Document ID query
docQuery := bleve.NewDocIDQuery([]string{"doc1", "doc2", "doc3"})
```

### Field Existence Query

Check if a field exists:

```go
// Field exists query
existsQuery := bleve.NewDocIDQuery([]string{"*"})
existsQuery.SetField("optional_field")
```

## Aggregations & Facets

### Terms Facet

Group results by field terms:

```go
func termsFacetSearch(index bleve.Index, query bleve.Query) (*bleve.SearchResult, error) {
    searchRequest := bleve.NewSearchRequest(query)
    
    // Add terms facet
    termsFacet := bleve.NewFacetRequest("tags", 10)
    searchRequest.AddFacet("tags_facet", termsFacet)
    
    result, err := index.Search(searchRequest)
    if err != nil {
        return nil, err
    }
    
    // Process facets
    if facets, found := result.Facets["tags_facet"]; found {
        for _, term := range facets.Terms {
            fmt.Printf("Term: %s, Count: %d\n", term.Term, term.Count)
        }
    }
    
    return result, nil
}
```

### Numeric Range Facet

Group results by numeric ranges:

```go
func numericFacetSearch(index bleve.Index, query bleve.Query) (*bleve.SearchResult, error) {
    searchRequest := bleve.NewSearchRequest(query)
    
    // Create numeric ranges
    numericRanges := []*bleve.NumericRange{
        {Name: "low", Min: 0, Max: 100},
        {Name: "medium", Min: 100, Max: 500},
        {Name: "high", Min: 500, Max: nil},
    }
    
    // Add numeric facet
    priceFacet := bleve.NewFacetRequest("price", 3)
    priceFacet.AddNumericRanges(numericRanges...)
    searchRequest.AddFacet("price_ranges", priceFacet)
    
    return index.Search(searchRequest)
}
```

### Date Range Facet

Group results by date ranges:

```go
func dateFacetSearch(index bleve.Index, query bleve.Query) (*bleve.SearchResult, error) {
    searchRequest := bleve.NewSearchRequest(query)
    
    // Create date ranges
    now := time.Now()
    dateRanges := []*bleve.DateTimeRange{
        {
            Name: "last_24h",
            Start: now.Add(-24 * time.Hour).Format(time.RFC3339),
            End:   now.Format(time.RFC3339),
        },
        {
            Name: "last_week",
            Start: now.Add(-7 * 24 * time.Hour).Format(time.RFC3339),
            End:   now.Add(-24 * time.Hour).Format(time.RFC3339),
        },
        {
            Name: "last_month",
            Start: now.Add(-30 * 24 * time.Hour).Format(time.RFC3339),
            End:   now.Add(-7 * 24 * time.Hour).Format(time.RFC3339),
        },
    }
    
    // Add date facet
    dateFacet := bleve.NewFacetRequest("created_at", 3)
    dateFacet.AddDateTimeRanges(dateRanges...)
    searchRequest.AddFacet("date_ranges", dateFacet)
    
    return index.Search(searchRequest)
}
```

### Multiple Facets Example

```go
func multiFacetSearch(index bleve.Index, query bleve.Query) (*bleve.SearchResult, error) {
    searchRequest := bleve.NewSearchRequest(query)
    
    // Terms facet
    tagsFacet := bleve.NewFacetRequest("tags", 10)
    searchRequest.AddFacet("popular_tags", tagsFacet)
    
    // Numeric facet
    priceRanges := []*bleve.NumericRange{
        {Name: "budget", Min: 0, Max: 50},
        {Name: "mid_range", Min: 50, Max: 200},
        {Name: "premium", Min: 200, Max: nil},
    }
    priceFacet := bleve.NewFacetRequest("price", 3)
    priceFacet.AddNumericRanges(priceRanges...)
    searchRequest.AddFacet("price_ranges", priceFacet)
    
    // Date facet
    now := time.Now()
    dateRanges := []*bleve.DateTimeRange{
        {
            Name: "recent",
            Start: now.Add(-7 * 24 * time.Hour).Format(time.RFC3339),
            End:   now.Format(time.RFC3339),
        },
        {
            Name: "older",
            End:   now.Add(-7 * 24 * time.Hour).Format(time.RFC3339),
        },
    }
    dateFacet := bleve.NewFacetRequest("created_at", 2)
    dateFacet.AddDateTimeRanges(dateRanges...)
    searchRequest.AddFacet("date_ranges", dateFacet)
    
    // Execute search
    result, err := index.Search(searchRequest)
    if err != nil {
        return nil, err
    }
    
    // Process all facets
    for facetName, facet := range result.Facets {
        fmt.Printf("\nFacet: %s\n", facetName)
        
        // Terms
        for _, term := range facet.Terms {
            fmt.Printf("  Term: %s, Count: %d\n", term.Term, term.Count)
        }
        
        // Numeric ranges
        for _, numRange := range facet.NumericRanges {
            fmt.Printf("  Range: %s, Count: %d\n", numRange.Name, numRange.Count)
        }
        
        // Date ranges
        for _, dateRange := range facet.DateRanges {
            fmt.Printf("  Range: %s, Count: %d\n", dateRange.Name, dateRange.Count)
        }
    }
    
    return result, nil
}
```

## Vector Search

Bleve supports vector similarity search and hybrid search combining vectors with keywords.

### Basic Vector Search

Simple vector similarity search:

```go
// Create search request
searchRequest := bleve.NewSearchRequest(bleve.NewMatchNoneQuery())

// Add vector search
vector := []float32{0.1, 0.2, 0.3, 0.4} // Your vector embedding
searchRequest.AddKNN("embedding", vector, 10, 1.0) // field, vector, k, boost

// Execute search
searchResult, err := index.Search(searchRequest)
```

### Hybrid Search

Combine vector similarity with text search:

```go
// Create text query
textQuery := bleve.NewMatchQuery("golang programming")
textQuery.SetField("content")
textQuery.SetBoost(0.3)

// Create search request
searchRequest := bleve.NewSearchRequest(textQuery)

// Add vector search with higher weight
vector := []float32{0.1, 0.2, 0.3, 0.4}
searchRequest.AddKNN("embedding", vector, 100, 0.7) // 70% weight for vector similarity

// Execute search
searchResult, err := index.Search(searchRequest)
```

### Vector Search with Filtering

Filter vector search results using boolean queries:

```go
// Create base query for filtering
minRating := float64(4.0)
ratingFilter := bleve.NewNumericRangeQuery(&minRating, nil)
ratingFilter.SetField("rating")

statusFilter := bleve.NewTermQuery("published")
statusFilter.SetField("status")

// Combine filters
filterQuery := bleve.NewConjunctionQuery(ratingFilter, statusFilter)

// Create search request
searchRequest := bleve.NewSearchRequest(filterQuery)

// Add vector search
vector := []float32{0.1, 0.2, 0.3, 0.4}
searchRequest.AddKNN("embedding", vector, 50, 1.0)

// Execute search
searchResult, err := index.Search(searchRequest)
```

### Vector Search with Multiple Fields

Search across multiple vector fields:

```go
func multiVectorSearch(index bleve.Index, textVector, imageVector []float32) (*bleve.SearchResult, error) {
    // Create base query
    searchRequest := bleve.NewSearchRequest(bleve.NewMatchNoneQuery())
    
    // Add text embedding search with weight
    searchRequest.AddKNN("text_embedding", textVector, 20, 0.6)  // Text similarity weight
    
    // Add image embedding search with weight
    searchRequest.AddKNN("image_embedding", imageVector, 20, 0.4)  // Image similarity weight
    
    // Execute search
    return index.Search(searchRequest)
}
```

### Advanced Vector Search Example

Here's a complete example showing various vector search features:

```go
func advancedVectorSearch(
    index bleve.Index,
    vector []float32,
    searchText string,
    minRating float64,
    category string,
) (*bleve.SearchResult, error) {
    // Text match with boost
    textQuery := bleve.NewMatchQuery(searchText)
    textQuery.SetField("content")
    textQuery.SetBoost(0.3)
    
    // Create filters
    ratingFilter := bleve.NewNumericRangeQuery(&minRating, nil)
    ratingFilter.SetField("rating")
    
    categoryFilter := bleve.NewTermQuery(category)
    categoryFilter.SetField("category")
    
    // Combine text query with filters
    finalQuery := bleve.NewConjunctionQuery(
        textQuery,
        ratingFilter,
        categoryFilter,
    )
    
    // Create search request with combined query
    searchRequest := bleve.NewSearchRequest(finalQuery)
    
    // Configure request options
    searchRequest.Size = 20
    searchRequest.From = 0
    searchRequest.Fields = []string{
        "title",
        "category",
        "rating",
        "created_at",
    }
    searchRequest.Highlight = bleve.NewHighlight()
    searchRequest.Highlight.Fields = []string{"content"}
    
    // Add vector search with weight
    searchRequest.AddKNN("embedding", vector, 100, 0.7)
    
    // Add facets
    categoryFacet := bleve.NewFacetRequest("category", 10)
    searchRequest.AddFacet("categories", categoryFacet)
    
    ratingRanges := []*bleve.NumericRange{
        {Name: "good", Min: 4.0, Max: 5.0},
        {Name: "average", Min: 3.0, Max: 4.0},
        {Name: "poor", Min: 0.0, Max: 3.0},
    }
    ratingFacet := bleve.NewFacetRequest("rating", 3)
    ratingFacet.AddNumericRanges(ratingRanges...)
    searchRequest.AddFacet("ratings", ratingFacet)
    
    // Execute search
    result, err := index.Search(searchRequest)
    if err != nil {
        return nil, err
    }
    
    // Process results
    fmt.Printf("Found %d matches in %s\n", 
        result.Total, 
        result.Took)
    
    for _, hit := range result.Hits {
        fmt.Printf("\nScore: %f\n", hit.Score)
        fmt.Printf("Title: %v\n", hit.Fields["title"])
        fmt.Printf("Category: %v\n", hit.Fields["category"])
        fmt.Printf("Rating: %v\n", hit.Fields["rating"])
        
        if fragments, ok := hit.Fragments["content"]; ok {
            fmt.Println("Matching content:")
            for _, fragment := range fragments {
                fmt.Printf("  %s\n", fragment)
            }
        }
    }
    
    return result, nil
}
```

### Vector Search Best Practices

1. **Performance Optimization**
   - Use appropriate vector dimensions (128-1024 typically)
   - Set reasonable K values for nearest neighbor search
   - Consider using filters to reduce search space
   - Index vectors with appropriate similarity metric

2. **Hybrid Search**
   - Balance weights between vector and text search using boosts
   - Use query combinations for complex filtering
   - Consider multiple ranking stages
   - Test different combination strategies

3. **Vector Quality**
   - Use high-quality embeddings
   - Normalize vectors when using cosine similarity
   - Consider dimensionality reduction for large vectors
   - Validate vector quality before indexing

4. **Resource Management**
   - Monitor memory usage with large vector indices
   - Consider batch processing for bulk operations
   - Implement timeouts for vector searches
   - Cache frequently used vectors

5. **Error Handling**
   - Validate vector dimensions
   - Handle missing embeddings gracefully
   - Implement fallback strategies
   - Monitor search quality metrics

## Complete Example

Here's a complete example showing various search features:

```go
package main

import (
    "fmt"
    "time"

    "github.com/blevesearch/bleve/v2"
)

func main() {
    // Open index
    index, err := bleve.Open("example.bleve")
    if err != nil {
        panic(err)
    }
    defer index.Close()

    // Create compound query
    // Match title or content, within date range, with minimum rating
    
    // 1. Text match
    titleQuery := bleve.NewMatchQuery("golang")
    titleQuery.SetField("title")
    titleQuery.SetBoost(2.0)

    contentQuery := bleve.NewMatchQuery("golang")
    contentQuery.SetField("content")

    textQuery := bleve.NewDisjunctionQuery(titleQuery, contentQuery)

    // 2. Date range
    startDate := time.Now().AddDate(-1, 0, 0).Format(time.RFC3339)
    endDate := time.Now().Format(time.RFC3339)
    dateQuery := bleve.NewDateRangeQuery(&startDate, &endDate)
    dateQuery.SetField("created_at")

    // 3. Numeric range for rating
    minRating := float64(4.0)
    ratingQuery := bleve.NewNumericRangeQuery(&minRating, nil)
    ratingQuery.SetField("rating")

    // Combine all queries
    conjunctionQuery := bleve.NewConjunctionQuery(
        textQuery,
        dateQuery,
        ratingQuery,
    )

    // Create search request
    searchRequest := bleve.NewSearchRequest(conjunctionQuery)
    
    // Configure request options
    searchRequest.Size = 10
    searchRequest.From = 0
    searchRequest.Fields = []string{"title", "rating", "created_at"}
    searchRequest.Highlight = bleve.NewHighlight()
    searchRequest.Highlight.Fields = []string{"content"}
    searchRequest.SortBy([]string{"_score", "-created_at"})

    // Execute search
    searchResult, err := index.Search(searchRequest)
    if err != nil {
        panic(err)
    }

    // Process results
    fmt.Printf("Found %d matches in %s\n", 
        searchResult.Total, 
        searchResult.Took)

    for _, hit := range searchResult.Hits {
        fmt.Printf("\nScore: %f\n", hit.Score)
        fmt.Printf("Title: %v\n", hit.Fields["title"])
        fmt.Printf("Rating: %v\n", hit.Fields["rating"])
        fmt.Printf("Created: %v\n", hit.Fields["created_at"])
        
        if fragments, ok := hit.Fragments["content"]; ok {
            fmt.Println("Matching content:")
            for _, fragment := range fragments {
                fmt.Printf("  %s\n", fragment)
            }
        }
    }
}
```

## Best Practices

1. **Query Performance**
   - Use field-specific queries when possible
   - Avoid wildcard queries with leading wildcards
   - Use match queries instead of term queries for text fields
   - Consider pagination for large result sets

2. **Result Handling**
   - Always check error returns from Search()
   - Use appropriate field types in result processing
   - Handle highlighting appropriately for your UI
   - Consider implementing retry logic for timeouts

3. **Memory Management**
   - Limit result size appropriately
   - Use field selection to reduce memory usage
   - Close index when done
   - Monitor memory usage during searches

4. **Query Design**
   - Use boost to control relevance
   - Consider analyzer choice when designing queries
   - Use compound queries for complex searches
   - Test queries with representative data

5. **Error Handling**
   - Handle common errors gracefully
   - Provide meaningful error messages
   - Implement timeout handling
   - Log search errors appropriately

## Common Patterns

### Search with Fallback

```go
func searchWithFallback(index bleve.Index, searchText string) (*bleve.SearchResult, error) {
    // Try exact phrase first
    phraseQuery := bleve.NewMatchPhraseQuery(searchText)
    searchRequest := bleve.NewSearchRequest(phraseQuery)
    result, err := index.Search(searchRequest)
    
    // If no results, try regular match
    if err == nil && result.Total == 0 {
        matchQuery := bleve.NewMatchQuery(searchText)
        searchRequest = bleve.NewSearchRequest(matchQuery)
        result, err = index.Search(searchRequest)
    }
    
    // If still no results, try fuzzy match
    if err == nil && result.Total == 0 {
        fuzzyQuery := bleve.NewFuzzyQuery(searchText)
        fuzzyQuery.SetFuzziness(2)
        searchRequest = bleve.NewSearchRequest(fuzzyQuery)
        result, err = index.Search(searchRequest)
    }
    
    return result, err
}
```

### Faceted Search

```go
func facetedSearch(index bleve.Index, query bleve.Query) (*bleve.SearchResult, error) {
    searchRequest := bleve.NewSearchRequest(query)
    
    // Add facets
    searchRequest.AddFacet("type", 
        bleve.NewFacetRequest("type", 10))
    searchRequest.AddFacet("rating",
        bleve.NewNumericFacetRequest("rating", 5))
    searchRequest.AddFacet("created",
        bleve.NewDateTimeFacetRequest("created_at", 5))
    
    return index.Search(searchRequest)
}
```

### Search with Timeout

```go
func searchWithTimeout(index bleve.Index, query bleve.Query, timeout time.Duration) (*bleve.SearchResult, error) {
    searchRequest := bleve.NewSearchRequest(query)
    
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    return index.SearchInContext(ctx, searchRequest)
}
```

These patterns provide a foundation for building robust search functionality in your applications. Adapt them to your specific needs and requirements. 