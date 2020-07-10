package raw

import (
	"errors"
)

//
// The raw query related structs.
// author: rnojiri
//

// Metadata - the raw data (metric and tags only)
type Metadata struct {
	Metric string            `json:"metric"`
	Tags   map[string]string `json:"tags"`
}

// Query - the raw data query JSON
type Query struct {
	Metadata
	Type         string `json:"type"`
	Since        string `json:"since"`
	Until        string `json:"until"`
	EstimateSize bool   `json:"estimateSize"`
}

const (
	rawDataQueryNumberType   string = "number"
	rawDataQueryTextType     string = "text"
	rawDataQueryMetricParam  string = "metric"
	rawDataQueryTagsParam    string = "tags"
	rawDataQuerySinceParam   string = "since"
	rawDataQueryUntilParam   string = "until"
	rawDataQueryEstimateSize string = "estimateSize"
	rawDataQueryTypeParam    string = "type"
	rawDataQueryFunc         string = "Parse"
	rawDataQueryKSID         string = "ksid"
	rawDataQueryTTL          string = "ttl"
)

var (
	// ErrUnmarshalling - unmarshalling error
	ErrUnmarshalling error = errors.New("error unmarshalling data")

	// ErrMissingMandatoryFields - mandatory fields are missing
	ErrMissingMandatoryFields error = errors.New("mandatory fields are missing")
)

// NumberPoint - represents a raw number point result
type NumberPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// TextPoint - represents a raw text point result
type TextPoint struct {
	Timestamp int64  `json:"timestamp"`
	Text      string `json:"text"`
}

// NumberPoints - the metadata and value results
type NumberPoints struct {
	Metadata Metadata      `json:"metadata"`
	Values   []NumberPoint `json:"points"`
}

// TextPoints - the metadata and text results
type TextPoints struct {
	Metadata Metadata    `json:"metadata"`
	Texts    []TextPoint `json:"points"`
}

// NumberQueryResults - the final raw query number results
type NumberQueryResults struct {
	Results []NumberPoints `json:"results"`
	Total   int            `json:"total"`
}

// TextQueryResults - the final raw query text results
type TextQueryResults struct {
	Results []TextPoints `json:"results"`
	Total   int          `json:"total"`
}
