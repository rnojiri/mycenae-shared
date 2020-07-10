package raw

import (
	"errors"
)

// DataMetadata - the raw data (metadata only)
type DataMetadata struct {
	Metric string            `json:"metric"`
	Tags   map[string]string `json:"tags"`
}

// DataQuery - the raw data query JSON
type DataQuery struct {
	DataMetadata
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

// DataNumberPoint - represents a raw number point result
type DataNumberPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// DataTextPoint - represents a raw text point result
type DataTextPoint struct {
	Timestamp int64  `json:"timestamp"`
	Text      string `json:"text"`
}

// DataQueryNumberPoints - the metadata and value results
type DataQueryNumberPoints struct {
	Metadata DataMetadata      `json:"metadata"`
	Values   []DataNumberPoint `json:"points"`
}

// DataQueryTextPoints - the metadata and text results
type DataQueryTextPoints struct {
	Metadata DataMetadata    `json:"metadata"`
	Texts    []DataTextPoint `json:"points"`
}

// DataQueryNumberResults - the final raw query number results
type DataQueryNumberResults struct {
	Results []DataQueryNumberPoints `json:"results"`
	Total   int                     `json:"total"`
}

// DataQueryTextResults - the final raw query text results
type DataQueryTextResults struct {
	Results []DataQueryTextPoints `json:"results"`
	Total   int                   `json:"total"`
}
