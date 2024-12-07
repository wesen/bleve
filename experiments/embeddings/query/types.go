package query

// SearchRequest represents the structure of the query DSL
type SearchRequest struct {
	Query   QueryDSL         `yaml:"query"`
	Options *SearchOptions   `yaml:"options,omitempty"`
	Facets  map[string]Facet `yaml:"facets,omitempty"`
}

// QueryDSL represents different types of queries
type QueryDSL struct {
	Match        *MatchQuery        `yaml:"match,omitempty"`
	MatchPhrase  *MatchPhraseQuery  `yaml:"match_phrase,omitempty"`
	Vector       *VectorQuery       `yaml:"vector,omitempty"`
	Bool         *BooleanQuery      `yaml:"bool,omitempty"`
	Term         *TermQuery         `yaml:"term,omitempty"`
	QueryString  *QueryStringQuery  `yaml:"query_string,omitempty"`
	Prefix       *PrefixQuery       `yaml:"prefix,omitempty"`
	Wildcard     *WildcardQuery     `yaml:"wildcard,omitempty"`
	NumericRange *NumericRangeQuery `yaml:"numeric_range,omitempty"`
	DateRange    *DateRangeQuery    `yaml:"date_range,omitempty"`
}

// MatchQuery represents a full-text search query
type MatchQuery struct {
	Field        string  `yaml:"field"`
	Value        string  `yaml:"value"`
	Boost        float64 `yaml:"boost,omitempty"`
	Operator     string  `yaml:"operator,omitempty"`
	Fuzziness    int     `yaml:"fuzziness,omitempty"`
	PrefixLength int     `yaml:"prefix_length,omitempty"`
}

// MatchPhraseQuery represents a phrase search query
type MatchPhraseQuery struct {
	Field string  `yaml:"field"`
	Value string  `yaml:"value"`
	Boost float64 `yaml:"boost,omitempty"`
	Slop  int     `yaml:"slop,omitempty"`
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
	Query        string  `yaml:"query"`
	DefaultField string  `yaml:"default_field,omitempty"`
	Boost        float64 `yaml:"boost,omitempty"`
}

// PrefixQuery represents a prefix-based search
type PrefixQuery struct {
	Field string  `yaml:"field"`
	Value string  `yaml:"value"`
	Boost float64 `yaml:"boost,omitempty"`
}

// WildcardQuery represents a wildcard pattern search
type WildcardQuery struct {
	Field string  `yaml:"field"`
	Value string  `yaml:"value"`
	Boost float64 `yaml:"boost,omitempty"`
}

// NumericRangeQuery represents a numeric range search
type NumericRangeQuery struct {
	Field        string   `yaml:"field"`
	Min          *float64 `yaml:"min,omitempty"`
	Max          *float64 `yaml:"max,omitempty"`
	InclusiveMin bool     `yaml:"inclusive_min,omitempty"`
	InclusiveMax bool     `yaml:"inclusive_max,omitempty"`
	Boost        float64  `yaml:"boost,omitempty"`
}

// DateRangeQuery represents a date range search
type DateRangeQuery struct {
	Field          string  `yaml:"field"`
	Start          string  `yaml:"start,omitempty"`
	End            string  `yaml:"end,omitempty"`
	InclusiveStart bool    `yaml:"inclusive_start,omitempty"`
	InclusiveEnd   bool    `yaml:"inclusive_end,omitempty"`
	Boost          float64 `yaml:"boost,omitempty"`
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
