package opentsdb

import (
	"errors"
	"fmt"
	"strings"
)

func parseDownsample(exp string, tsdb *Expression) (string, error) {

	params := parseParams(string(exp[10:]))

	if len(params) != 4 {
		return stringsEmpty, fmt.Errorf("downsample expects 4 parameters but found %d: %v", len(params), params)
	}

	tsdb.Downsample = fmt.Sprintf("%s-%s-%s", params[0], params[1], params[2])

	for _, oper := range tsdb.Order {
		if oper == "downsample" {
			return stringsEmpty, errors.New("found more than one 'downsample' function")
		}
	}

	tsdb.Order = append([]string{"downsample"}, tsdb.Order...)

	return params[3], nil
}

func writeDownsample(exp, dsInfo string) string {
	if dsInfo != stringsEmpty {
		info := strings.Split(dsInfo, "-")
		if len(info) == 2 {
			info = append(info, "none")
		}
		exp = fmt.Sprintf("downsample(%s,%s,%s,%s)", info[0], info[1], info[2], exp)
	}
	return exp
}
