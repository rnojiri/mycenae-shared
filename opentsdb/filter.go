package opentsdb

import (
	"errors"
	"fmt"
)

func parseFilter(exp string, tsdb *Expression) (string, error) {

	params := parseParams(string(exp[6:]))

	if len(params) != 2 {
		return stringsEmpty, fmt.Errorf("filter expects 2 parameters but found %d: %v", len(params), params)
	}

	tsdb.FilterValue = params[0]

	for _, oper := range tsdb.Order {
		if oper == "filterValue" {
			return stringsEmpty, errors.New("found more than one 'filterValue' function")
		}
	}

	tsdb.Order = append([]string{"filterValue"}, tsdb.Order...)

	return params[1], nil
}

func writeFilter(exp, filterValue string) string {
	if filterValue != stringsEmpty {
		return fmt.Sprintf("filter(%s,%s)", filterValue, exp)
	}
	return exp
}
