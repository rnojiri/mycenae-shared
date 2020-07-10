package raw

import (
	"github.com/buger/jsonparser"
)

//
// The raw query json parser.
// author: rnojiri
//

// Parse - parses the bytes tol JSON
func (dq *Query) Parse(data []byte) error {

	var err error

	if dq.Type, err = jsonparser.GetString(data, rawDataQueryTypeParam); err != nil {
		return ErrUnmarshalling
	}

	if dq.Type != rawDataQueryNumberType && dq.Type != rawDataQueryTextType {
		return ErrMissingMandatoryFields
	}

	if dq.Metric, err = jsonparser.GetString(data, rawDataQueryMetricParam); err != nil {
		return ErrUnmarshalling
	}

	if dq.Since, err = jsonparser.GetString(data, rawDataQuerySinceParam); err != nil {
		return ErrUnmarshalling
	}

	if dq.Until, err = jsonparser.GetString(data, rawDataQueryUntilParam); err != nil && err != jsonparser.KeyPathNotFoundError {
		return ErrUnmarshalling
	}

	if dq.EstimateSize, err = jsonparser.GetBoolean(data, rawDataQueryEstimateSize); err != nil && err != jsonparser.KeyPathNotFoundError {
		return ErrUnmarshalling
	}

	dq.Tags = map[string]string{}
	err = jsonparser.ObjectEach(data, func(key, value []byte, dataType jsonparser.ValueType, offset int) error {

		tagKey, err := jsonparser.ParseString(key)
		if err != nil {
			return ErrUnmarshalling
		}

		if dq.Tags[tagKey], err = jsonparser.ParseString(value); err != nil {
			return ErrUnmarshalling
		}

		return nil

	}, rawDataQueryTagsParam)

	if _, ok := dq.Tags[rawDataQueryKSID]; !ok {
		return ErrMissingMandatoryFields
	}

	return nil
}
