package opentsdb

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

func parseGroup(exp string, tsdb *Expression) (string, error) {

	var queryExp string

	for i := 1; i < len(exp); i++ {

		if string(exp[i]) == "(" {
			f := 1
			for j := i + 1; j < len(exp); j++ {

				if string(exp[j]) == "(" {
					f++
				}

				if string(exp[j]) == ")" {
					f--
				}

				if f == 0 {
					i = j + 1
					if i == len(exp) {
						return stringsEmpty, errors.New("groupBy cannot be used by itself")
					}

					if string(exp[i]) != "|" {
						return stringsEmpty, errors.New("groupBy should be followed by a |")
					}

					if i+1 == len(exp) {
						return stringsEmpty, errors.New("groupBy should be followed by a | and a query expression")
					}

					queryExp = exp[i+1:]
					exp = exp[:i]
					break
				}
			}
		}
	}

	params := parseParams(string(exp[7:]))

	if len(params) != 1 {
		return stringsEmpty, fmt.Errorf("groupBy expects 1 parameter but found %d: %v", len(params), params)
	}

	tags, err := parseMap(params[0])
	if err != nil {
		return stringsEmpty, err
	}

	for k, vs := range tags {
		for _, v := range vs {

			var ft, cv string

			if strings.HasPrefix(v, "regexp(") && strings.HasSuffix(v, ")") {
				ft = "regexp"
				cv = v[7 : len(v)-1]
			} else if strings.HasPrefix(v, "wildcard(") && strings.HasSuffix(v, ")") {
				ft = "wildcard"
				cv = v[9 : len(v)-1]
			} else if strings.HasPrefix(v, "or(") && strings.HasSuffix(v, ")") {
				ft = "literal_or"
				cv = v[3 : len(v)-1]
			} else if strings.HasPrefix(v, "notor(") && strings.HasSuffix(v, ")") {
				ft = "not_literal_or"
				cv = v[6 : len(v)-1]
			} else {
				ft = "wildcard"
				cv = v
			}

			filter := Filter{
				Ftype:   ft,
				Tagk:    k,
				Filter:  cv,
				GroupBy: true,
			}
			tsdb.Filters = append(tsdb.Filters, filter)
		}
	}

	for _, oper := range tsdb.Order {
		if oper == "groupBy" {
			return stringsEmpty, errors.New("found more than one 'groupBy' function")
		}
	}

	tsdb.Order = append([]string{"groupBy"}, tsdb.Order...)

	return queryExp, nil
}

func writeGroup(exp string, filters []Filter) string {

	gExp := stringsEmpty

	orderedTags := []string{}

	joinFilters := map[string][]string{}

	for _, filter := range filters {
		if filter.GroupBy {
			if _, ok := joinFilters[filter.Tagk]; !ok {
				switch filter.Ftype {
				case "wildcard":
					joinFilters[filter.Tagk] = []string{
						filter.Filter,
					}
				case "regexp":
					joinFilters[filter.Tagk] = []string{
						fmt.Sprintf("%s(%s)", filter.Ftype, filter.Filter),
					}
				case "literal_or":
					joinFilters[filter.Tagk] = []string{
						fmt.Sprintf("or(%s)", filter.Filter),
					}
				case "not_literal_or":
					joinFilters[filter.Tagk] = []string{
						fmt.Sprintf("notor(%s)", filter.Filter),
					}
				}
				orderedTags = append(orderedTags, filter.Tagk)
			} else {
				switch filter.Ftype {
				case "wildcard":
					joinFilters[filter.Tagk] = append(
						joinFilters[filter.Tagk],
						filter.Filter,
					)
				case "regexp":
					joinFilters[filter.Tagk] = append(
						joinFilters[filter.Tagk],
						fmt.Sprintf("%s(%s)", filter.Ftype, filter.Filter),
					)
				case "literal_or":
					joinFilters[filter.Tagk] = append(
						joinFilters[filter.Tagk],
						fmt.Sprintf("or(%s)", filter.Ftype),
					)
					joinFilters[filter.Tagk] = []string{
						fmt.Sprintf("or(%s)", filter.Filter),
					}
				case "not_literal_or":
					joinFilters[filter.Tagk] = append(
						joinFilters[filter.Tagk],
						fmt.Sprintf("notor(%s)", filter.Ftype),
					)
				}
			}
		}
	}

	if len(orderedTags) > 0 {

		gExp = "groupBy({"

		sort.Strings(orderedTags)

		for _, tk := range orderedTags {

			sort.Strings(joinFilters[tk])

			for _, fv := range joinFilters[tk] {
				gExp = fmt.Sprintf("%s%s=%s,", gExp, tk, fv)
			}
		}

		gExp = gExp[:len(gExp)-1]

		exp = fmt.Sprintf("%s})|%s", gExp, exp)
	}

	return exp
}
