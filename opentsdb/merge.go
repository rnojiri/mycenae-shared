package opentsdb

import (
	"errors"
	"fmt"
)

func parseMerge(exp string, tsdb *Expression) (string, error) {

	params := parseParams(string(exp[5:]))

	if len(params) != 2 {
		return stringsEmpty, fmt.Errorf("merge expects 2 parameters but found %d: %v", len(params), params)
	}

	tsdb.Aggregator = params[0]

	for _, oper := range tsdb.Order {
		if oper == "aggregation" {
			return stringsEmpty, errors.New("found more than one 'aggregation' function")
		}
	}

	tsdb.Order = append([]string{"aggregation"}, tsdb.Order...)

	return params[1], nil
}

func writeMerge(exp, operator string) string {
	return fmt.Sprintf("merge(%s,%s)", operator, exp)
}
