package opentsdb

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	validFieldRegexp   = regexp.MustCompile(`^[0-9A-Za-z-._%&#;\\/]+$`)
	validFieldWildcard = regexp.MustCompile(`^[0-9A-Za-z-._%&#;\\/*]+$`)
	validFieldLiteral  = regexp.MustCompile(`^[0-9A-Za-z-._%&#;\\/|]+$`)
)

const (
	stringsEmpty      string = ""
	stringsWhiteSpace string = " "
)

// Expression - opentsdb expression
type Expression struct {
	Aggregator  string            `json:"aggregator"`
	Downsample  string            `json:"downsample,omitempty"`
	Metric      string            `json:"metric"`
	Tags        map[string]string `json:"tags"`
	Rate        bool              `json:"rate,omitempty"`
	RateOptions Rate              `json:"rateOptions,omitempty"`
	Order       []string          `json:"order,omitempty"`
	FilterValue string            `json:"filterValue,omitempty"`
	Filters     []Filter          `json:"filters,omitempty"`
}

// Query - the main query container
type Query struct {
	Start        int64        `json:"start,omitempty"`
	End          int64        `json:"end,omitempty"`
	Relative     string       `json:"relative,omitempty"`
	Queries      []Expression `json:"queries"`
	ShowTSUIDs   bool         `json:"showTSUIDs"`
	MsResolution bool         `json:"msResolution"`
	EstimateSize bool         `json:"estimateSize"`
}

// Validate - validates the payload
func (query *Query) Validate() error {

	if query.Relative != stringsEmpty {
		if err := query.checkDuration(query.Relative); err != nil {
			return err
		}
	}

	if len(query.Queries) == 0 {
		return errors.New("at least one query should be present")
	}

	for i, q := range query.Queries {

		if err := query.checkField("metric", q.Metric); err != nil {
			return err
		}

		if err := query.checkAggregator(q.Aggregator); err != nil {
			return err
		}

		if q.Downsample != stringsEmpty {

			ds := strings.Split(q.Downsample, "-")

			if len(ds) < 2 {
				return errors.New("invalid downsample format")
			}

			if err := query.checkDuration(ds[0]); err != nil {
				return err
			}

			if err := query.checkDownsampler(ds[1]); err != nil {
				return err
			}

			if len(ds) > 2 {
				if err := query.checkFiller(ds[2]); err != nil {
					return err
				}
			}

		}

		if q.Rate {
			if err := query.checkRate(q.RateOptions); err != nil {
				return err
			}
		}

		if q.FilterValue != stringsEmpty {
			q.FilterValue = strings.Replace(q.FilterValue, stringsWhiteSpace, stringsEmpty, -1)
			query.Queries[i].FilterValue = q.FilterValue

			if len(q.FilterValue) < 2 {
				return fmt.Errorf("invalid filter value %s", q.FilterValue)
			}

			if q.FilterValue[:2] == ">=" || q.FilterValue[:2] == "<=" || q.FilterValue[:2] == "==" || q.FilterValue[:2] == "!=" {
				_, err := strconv.ParseFloat(q.FilterValue[2:], 64)
				if err != nil {
					return err
				}
			} else if q.FilterValue[:1] == ">" || q.FilterValue[:1] == "<" {
				_, err := strconv.ParseFloat(q.FilterValue[1:], 64)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("invalid filter value %s", q.FilterValue)
			}
		}

		if len(q.Order) == 0 {

			if q.FilterValue != stringsEmpty {
				query.Queries[i].Order = append(query.Queries[i].Order, "filterValue")
			}

			if q.Downsample != stringsEmpty {
				query.Queries[i].Order = append(query.Queries[i].Order, "downsample")
			}

			query.Queries[i].Order = append(query.Queries[i].Order, "aggregation")

			if q.Rate {
				query.Queries[i].Order = append(query.Queries[i].Order, "rate")
			}

		} else {

			orderCheck := make([]string, len(q.Order))

			copy(orderCheck, q.Order)

			k := 0
			occur := 0
			for j, order := range orderCheck {

				if order == "aggregation" {
					k = j
					occur++
				}

			}

			if occur == 0 {
				return errors.New("aggregation configured but no aggregation found in order array")
			}

			if occur > 1 {
				return errors.New("more than one aggregation found in order array")
			}

			if occur == 1 {
				orderCheck = append(orderCheck[:k], orderCheck[k+1:]...)
			}

			k = 0
			occur = 0
			for j, order := range orderCheck {

				if order == "filterValue" {
					k = j
					occur++
				}

			}

			if q.FilterValue != stringsEmpty && occur == 0 {
				return errors.New("filterValue configured but no filterValue found in order array")
			}

			if occur > 1 {
				return errors.New("more than one filterValue found in order array")
			}

			if occur == 1 {
				orderCheck = append(orderCheck[:k], orderCheck[k+1:]...)
			}

			k = 0
			occur = 0
			for j, order := range orderCheck {

				if order == "downsample" {
					k = j
					occur++
				}

			}

			if q.Downsample != stringsEmpty && occur == 0 {
				return errors.New("downsample configured but no downsample found in order array")
			}

			if occur > 1 {
				return errors.New("more than one downsample found in order array")
			}

			if occur == 1 {
				orderCheck = append(orderCheck[:k], orderCheck[k+1:]...)
			}

			k = 0
			occur = 0
			for j, order := range orderCheck {

				if order == "rate" {
					k = j
					occur++
				}

			}

			if q.Rate && occur == 0 {
				return errors.New("rate configured but no rate found in order array")
			}

			if occur > 1 {
				return errors.New("more than one rate found in order array")
			}

			if occur == 1 {
				orderCheck = append(orderCheck[:k], orderCheck[k+1:]...)
			}

			if len(orderCheck) != 0 {
				return fmt.Errorf("invalid operations in order array %v", orderCheck)
			}

		}

		if err := query.checkFilter(q.Filters); err != nil {
			return err
		}

	}

	return nil
}

func (query *Query) checkRate(opts Rate) error {

	if opts.CounterMax != nil && *opts.CounterMax < 0 {
		return errors.New("counter max needs to be a positive integer")
	}

	return nil
}

func (query *Query) checkAggregator(aggr string) error {

	ok := false

	for _, vAggr := range GetAggregators() {
		if vAggr == aggr {
			ok = true
			break
		}
	}

	if !ok {
		return errors.New("unknown aggregation value")
	}

	return nil
}

func (query *Query) checkDownsampler(DSr string) error {

	ok := false

	for _, vDSr := range GetDownsamplers() {
		if vDSr == DSr {
			ok = true
			break
		}
	}

	if !ok {
		return errors.New("invalid downsample")
	}

	return nil
}

func (query *Query) checkFiller(DSf string) error {

	ok := false

	for _, vDSf := range GetDownsampleFillers() {
		if vDSf == DSf {
			ok = true
			break
		}
	}

	if !ok {
		return errors.New("invalid fill value")
	}

	return nil
}

func (query *Query) checkFilter(filters []Filter) error {

	vFilters := GetFilters()

	for _, filter := range filters {

		ok := false

		ft := filter.Ftype

		if ft == "iliteral_or" {
			ft = "literal_or"
		} else if ft == "not_iliteral_or" {
			ft = "not_literal_or"
		} else if ft == "iwildcard" {
			ft = "wildcard"
		}

		for _, vFilter := range vFilters {
			if ft == vFilter {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("invalid filter type %s", filter.Ftype)
		}

		if err := query.checkField("tagk", filter.Tagk); err != nil {
			return err
		}

		if err := query.checkFilterField("filter", ft, filter.Filter); err != nil {
			return err
		}
	}

	return nil
}

func (query *Query) checkDuration(s string) error {

	if len(s) < 2 {
		return errors.New("invalid time interval")
	}

	var n int
	var err error

	if string(s[len(s)-2:]) == "ms" {
		n, err = strconv.Atoi(string(s[:len(s)-2]))
		if err != nil {
			return err
		}
		return nil
	}

	switch s[len(s)-1:] {
	case "s", "m", "h", "d", "w", "n", "y":
		n, err = strconv.Atoi(string(s[:len(s)-1]))
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid unit")
	}

	if n < 1 {
		return errors.New("interval needs to be bigger than 0")
	}

	return nil
}

func (query *Query) checkField(n, f string) error {

	if !validFieldRegexp.MatchString(f) {
		return fmt.Errorf("Invalid characters in field %s: %s", n, f)
	}

	return nil
}

func (query *Query) checkFilterField(n, tf, f string) error {

	match := false

	switch tf {
	case "wildcard":
		match = validFieldWildcard.MatchString(f)
	case "literal_or", "not_literal_or":
		match = validFieldLiteral.MatchString(f)
	case "regexp":
		match = true
	}

	if !match {
		return fmt.Errorf("Invalid characters in field %s: %s", n, f)
	}

	return nil
}

// Rate - rate options
type Rate struct {
	Counter    bool   `json:"counter"`
	CounterMax *int64 `json:"counterMax,omitempty"`
	ResetValue int64  `json:"resetValue,omitempty"`
}

// Filter - filter options
type Filter struct {
	Ftype   string `json:"type"`
	Tagk    string `json:"tagk"`
	Filter  string `json:"filter"`
	GroupBy bool   `json:"groupBy"`
}

// Points - an array of point
type Points []*Point

// Tag - a tag from the opentsdb point
type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Point - an opentsdb point
type Point struct {
	Metric    string   `json:"metric"`
	Timestamp int64    `json:"timestamp"`
	Value     *float64 `json:"value"`
	Text      string   `json:"text"`
	Tags      []Tag    `json:"tags"`
	TTL       int      `json:"ttl"`
	Keyset    string   `json:"keyset"`
}
