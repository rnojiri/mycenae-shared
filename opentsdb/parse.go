package opentsdb

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// GetRelativeStart - returns a start time based on an end time and a duration string
func GetRelativeStart(end time.Time, s string) (time.Time, error) {

	if string(s[len(s)-2:]) == "ms" {
		d, err := time.ParseDuration(s)
		return end.Add(-d), err
	}

	switch s[len(s)-1:] {
	case "s", "m", "h":
		d, err := time.ParseDuration(s)
		return end.Add(-d), err
	case "d":
		i, err := strconv.Atoi(string(s[:len(s)-1]))
		return end.AddDate(0, 0, -i), err
	case "w":
		i, err := strconv.Atoi(string(s[:len(s)-1]))
		return end.AddDate(0, 0, -i*7), err
	case "n":
		i, err := strconv.Atoi(string(s[:len(s)-1]))
		return end.AddDate(0, -i, 0), err
	case "y":
		i, err := strconv.Atoi(string(s[:len(s)-1]))
		return end.AddDate(-i, 0, 0), err
	}

	return time.Time{}, fmt.Errorf("unknown time unit: %s", s[len(s)-1:])
}

func parseParams(exp string) []string {

	var param []byte

	params := []string{}

	for i := 1; i < len(exp); i++ {

		if string(exp[i]) == "(" {
			param = append(param, exp[i])
			f := 1
			for j := i + 1; j < len(exp); j++ {

				if string(exp[j]) == "(" {
					f++
				}

				if string(exp[j]) == ")" {
					f--
				}

				param = append(param, exp[j])

				if f == 0 {
					i = j + 1
					if i == len(exp) {
						return params
					}
					break
				}
			}
		}

		if string(exp[i]) == "{" {
			param = append(param, exp[i])
			for j := i + 1; j < len(exp); j++ {

				param = append(param, exp[j])

				if string(exp[j]) == "}" {
					i = j + 1
					break
				}
			}
		}

		if string(exp[i]) == "," {
			params = append(params, string(param))
			param = []byte{}
			continue
		}

		if string(exp[i]) == ")" {
			if i+1 == len(exp) {
				params = append(params, string(param))
				break
			}
			return params
		}

		param = append(param, exp[i])
	}

	return params
}

func parseMap(exp string) (map[string][]string, error) {

	if len(exp) == 0 {
		return nil, errors.New(`empty map`)
	}

	if string(exp[0]) != "{" {
		return nil, errors.New(`missing '{' at the beginning of map`)
	}

	var key, value []byte

	m := map[string][]string{}

	for i := 1; i < len(exp); i++ {

		if string(exp[i]) == "=" {

			if len(key) == 0 {
				return nil, errors.New(`map key cannot be empty`)
			}

			if _, ok := m[string(key)]; !ok {
				m[string(key)] = []string{}
			}

			for j := i + 1; j < len(exp); j++ {

				if string(exp[j]) == "," || string(exp[j]) == "}" {
					if len(value) == 0 {
						return nil, errors.New(`map value cannot be empty`)
					}
					m[string(key)] = append(m[string(key)], string(value))
					key = []byte{}
					value = []byte{}
					i = j
					break
				}
				value = append(value, exp[j])
			}
			continue
		}

		if string(exp[i]) == "," || string(exp[i]) == "}" {
			return nil, errors.New(`bad map format`)
		}

		key = append(key, exp[i])
	}

	return m, nil
}
