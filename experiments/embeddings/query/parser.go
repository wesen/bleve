package query

import (
	"fmt"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	bleve_query "github.com/blevesearch/bleve/v2/search/query"

	"github.com/blevesearch/bleve/v2/experiments/embeddings/embeddings"
)

var embeddingsClient *embeddings.Client

func init() {
	embeddingsClient = embeddings.DefaultClient()
}

// BuildBleveQuery converts a QueryDSL to a bleve.Query
func BuildBleveQuery(q QueryDSL) (bleve_query.Query, error) {
	if q.Match != nil {
		// Create match query
		query := bleve.NewMatchQuery(q.Match.Value)
		query.SetField(q.Match.Field)

		// Set boost if provided
		if q.Match.Boost != 0 {
			query.SetBoost(q.Match.Boost)
		}

		// Set operator if provided
		if q.Match.Operator != "" {
			op := strings.ToLower(q.Match.Operator)
			switch op {
			case "or":
				query.SetOperator(bleve_query.MatchQueryOperatorOr)
			case "and":
				query.SetOperator(bleve_query.MatchQueryOperatorAnd)
			default:
				return nil, fmt.Errorf("invalid operator %q, must be 'and' or 'or'", q.Match.Operator)
			}
		}

		// Set fuzziness if provided
		if q.Match.Fuzziness != 0 {
			if q.Match.Fuzziness < 0 {
				return nil, fmt.Errorf("fuzziness must be non-negative, got %d", q.Match.Fuzziness)
			}
			query.SetFuzziness(q.Match.Fuzziness)
		}

		// Set prefix length if provided
		if q.Match.PrefixLength != 0 {
			if q.Match.PrefixLength < 0 {
				return nil, fmt.Errorf("prefix_length must be non-negative, got %d", q.Match.PrefixLength)
			}
			query.SetPrefix(q.Match.PrefixLength)
		}

		return query, nil
	}

	if q.MatchPhrase != nil {
		// Create match phrase query
		query := bleve.NewMatchPhraseQuery(q.MatchPhrase.Value)
		query.SetField(q.MatchPhrase.Field)

		// Set boost if provided
		if q.MatchPhrase.Boost != 0 {
			query.SetBoost(q.MatchPhrase.Boost)
		}

		return query, nil
	}

	if q.Vector != nil {
		var queryVector []float32
		var err error

		if q.Vector.Text != "" {
			queryVector, err = embeddingsClient.GenerateEmbedding(q.Vector.Text)
			if err != nil {
				return nil, fmt.Errorf("failed to generate vector embedding: %w", err)
			}
		} else if q.Vector.Vector != nil {
			queryVector = q.Vector.Vector
		} else {
			return nil, fmt.Errorf("either text or vector must be provided for vector query")
		}

		searchRequest := bleve.NewSearchRequest(bleve.NewMatchAllQuery())
		searchRequest.Size = q.Vector.K
		searchRequest.AddKNN(q.Vector.Field, queryVector, int64(q.Vector.K), q.Vector.Boost)
		return searchRequest.Query, nil
	}

	if q.Bool != nil {
		boolQuery := bleve.NewBooleanQuery()

		for _, must := range q.Bool.Must {
			q, err := BuildBleveQuery(must)
			if err != nil {
				return nil, err
			}
			boolQuery.AddMust(q)
		}

		for _, should := range q.Bool.Should {
			q, err := BuildBleveQuery(should)
			if err != nil {
				return nil, err
			}
			boolQuery.AddShould(q)
		}

		for _, mustNot := range q.Bool.MustNot {
			q, err := BuildBleveQuery(mustNot)
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
		query := bleve.NewQueryStringQuery(q.QueryString.Query)
		if q.QueryString.Boost != 0 {
			query.SetBoost(q.QueryString.Boost)
		}
		// XXX remove default field
		return query, nil
	}

	if q.Prefix != nil {
		query := bleve.NewPrefixQuery(q.Prefix.Value)
		query.SetField(q.Prefix.Field)
		if q.Prefix.Boost != 0 {
			query.SetBoost(q.Prefix.Boost)
		}
		return query, nil
	}

	if q.Wildcard != nil {
		query := bleve.NewWildcardQuery(q.Wildcard.Value)
		query.SetField(q.Wildcard.Field)
		if q.Wildcard.Boost != 0 {
			query.SetBoost(q.Wildcard.Boost)
		}
		return query, nil
	}

	if q.NumericRange != nil {
		// // NewNumericRangeInclusiveQuery creates a new Query for ranges
		// of numeric values.
		// Either, but not both endpoints can be nil.
		// Control endpoint inclusion with inclusiveMin, inclusiveMax.
		// func NewNumericRangeInclusiveQuery(min, max *float64, minInclusive, maxInclusive *bool) *query.NumericRangeQuery {
		// 	return query.NewNumericRangeInclusiveQuery(min, max, minInclusive, maxInclusive)
		// }

		// XXX Handle inclusive/exclusive
		query := bleve.NewNumericRangeQuery(q.NumericRange.Min, q.NumericRange.Max)
		query.SetField(q.NumericRange.Field)
		if q.NumericRange.Boost != 0 {
			query.SetBoost(q.NumericRange.Boost)
		}
		return query, nil
	}

	if q.DateRange != nil {
		var startTime, endTime *time.Time

		if q.DateRange.Start != "" {
			t, err := time.Parse(time.RFC3339, q.DateRange.Start)
			if err != nil {
				return nil, fmt.Errorf("invalid start date %q: %w", q.DateRange.Start, err)
			}
			startTime = &t
		}

		if q.DateRange.End != "" {
			t, err := time.Parse(time.RFC3339, q.DateRange.End)
			if err != nil {
				return nil, fmt.Errorf("invalid end date %q: %w", q.DateRange.End, err)
			}
			endTime = &t
		}

		query := bleve.NewDateRangeQuery(*startTime, *endTime)
		query.SetField(q.DateRange.Field)
		// XXX handle inclusive date range query
		if q.DateRange.Boost != 0 {
			query.SetBoost(q.DateRange.Boost)
		}
		return query, nil
	}

	return nil, fmt.Errorf("no valid query type found")
}

// ApplySearchOptions applies the search options to a search request
func ApplySearchOptions(searchRequest *bleve.SearchRequest, options *SearchOptions) {
	if options == nil {
		return
	}

	if options.Size > 0 {
		searchRequest.Size = options.Size
	}
	if options.From > 0 {
		searchRequest.From = options.From
	}
	if len(options.Fields) > 0 {
		searchRequest.Fields = options.Fields
	}
	if options.Explain {
		searchRequest.Explain = true
	}
	if options.Highlight != nil {
		searchRequest.Highlight = bleve.NewHighlight()
		searchRequest.Highlight.Fields = options.Highlight.Fields
	}
	// Apply sorting
	for _, sort := range options.Sort {
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
