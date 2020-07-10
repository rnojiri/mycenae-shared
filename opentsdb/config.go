package opentsdb

// FilterInfo - the description from a filter
type FilterInfo struct {
	Examples    string `json:"examples"`
	Description string `json:"description"`
}

// GetAggregators - returns the implemented aggregators
func GetAggregators() []string {
	return []string{
		"avg",
		"count",
		"min",
		"max",
		"sum",
	}
}

// GetFilters - returns the implemented filters
func GetFilters() []string {
	return []string{
		"literal_or",
		"not_literal_or",
		"wildcard",
		"regexp",
	}
}

// GetFiltersFull - returns a map of filters and their descriptions
func GetFiltersFull() map[string]FilterInfo {
	return map[string]FilterInfo{
		"literal_or": {
			Examples:    `host=iliteral_or(web01),  host=iliteral_or(web01|web02|web03)  {\"type\":\"iliteral_or\",\"tagk\":\"host\",\"filter\":\"web01|web02|web03\",\"groupBy\":false}`,
			Description: `Accepts one or more exact values and matches if the series contains any of them. Multiple values can be included and must be separated by the | (pipe) character. The filter is case insensitive and will not allow characters that TSDB does not allow at write time.`,
		},
		"not_literal_or": {
			Examples:    `host=not_literal_or(web01),  host=not_literal_or(web01|web02|web03)  {\"type\":\"not_literal_or\",\"tagk\":\"host\",\"filter\":\"web01|web02|web03\",\"groupBy\":false}`,
			Description: `Accepts one or more exact values and matches if the series does NOT contain any of them. Multiple values can be included and must be separated by the | (pipe) character. The filter is case sensitive and will not allow characters that TSDB does not allow at write time.`,
		},
		"wildcard": {
			Examples:    `host=wildcard(web*),  host=wildcard(web*.tsdb.net)  {\"type\":\"wildcard\",\"tagk\":\"host\",\"filter\":\"web*.tsdb.net\",\"groupBy\":false}`,
			Description: `Performs pre, post and in-fix glob matching of values. The globs are case sensitive and multiple wildcards can be used. The wildcard character is the * (asterisk). At least one wildcard must be present in the filter value. A wildcard by itself can be used as well to match on any value for the tag key.`,
		},
		"regexp": {
			Examples:    `host=regexp(.*)  {\"type\":\"regexp\",\"tagk\":\"host\",\"filter\":\".*\",\"groupBy\":false}`,
			Description: `Provides full, POSIX compliant regular expression using the built in Java Pattern class. Note that an expression containing curly braces {} will not parse properly in URLs. If the pattern is not a valid regular expression then an exception will be raised.`,
		},
	}
}

// GetDownsamplers - returns the implemented downsamplers
func GetDownsamplers() []string {
	return []string{
		"avg",
		"count",
		"min",
		"max",
		"sum",
	}
}

// GetDownsampleFillers - returns the downsample fillers
func GetDownsampleFillers() []string {
	return []string{
		"none",
		"nan",
		"null",
		"zero",
	}
}
