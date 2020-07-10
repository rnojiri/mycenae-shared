package opentsdb

import (
	"errors"
	"fmt"
	"strconv"
)

func parseRate(exp string, tsdb *Expression) (string, error) {

	params := parseParams(string(exp[4:]))

	if len(params) != 4 {
		return stringsEmpty, fmt.Errorf("rate expects 4 parameters but found %d: %v", len(params), params)
	}

	b, err := strconv.ParseBool(params[0])
	if err != nil {
		return stringsEmpty, err
	}

	tsdb.RateOptions.Counter = b

	if params[1] != "null" {
		counterMax, err := strconv.ParseInt(params[1], 10, 64)
		if err != nil {
			return stringsEmpty, err
		}
		tsdb.RateOptions.CounterMax = &counterMax
	}

	tsdb.RateOptions.ResetValue, err = strconv.ParseInt(params[2], 10, 64)
	if err != nil {
		return stringsEmpty, err
	}

	tsdb.Rate = true

	for _, oper := range tsdb.Order {
		if oper == "rate" {
			return stringsEmpty, errors.New("found more than one 'rate' function")
		}
	}

	tsdb.Order = append([]string{"rate"}, tsdb.Order...)

	return params[3], nil
}

func writeRate(exp string, rate bool, rateOptions Rate) string {
	if rate {
		cm := "null"

		if rateOptions.CounterMax != nil {
			cm = fmt.Sprintf("%d", *rateOptions.CounterMax)
		}

		exp = fmt.Sprintf("rate(%t,%s,%d,%s)", rateOptions.Counter, cm, rateOptions.ResetValue, exp)
	}
	return exp
}
