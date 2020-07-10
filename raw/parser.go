package raw

import (
	"io/ioutil"
	"net/http"

	"github.com/buger/jsonparser"
)

// Parse - parses the bytes tol JSON
func (rq *DataQuery) Parse(r *http.Request) error {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ErrUnmarshalling
	}

	if rq.Type, err = jsonparser.GetString(data, rawDataQueryTypeParam); err != nil {
		return ErrUnmarshalling
	}

	if rq.Type != rawDataQueryNumberType && rq.Type != rawDataQueryTextType {
		return ErrMissingMandatoryFields
	}

	if rq.Metric, err = jsonparser.GetString(data, rawDataQueryMetricParam); err != nil {
		return ErrUnmarshalling
	}

	if rq.Since, err = jsonparser.GetString(data, rawDataQuerySinceParam); err != nil {
		return ErrUnmarshalling
	}

	if rq.Until, err = jsonparser.GetString(data, rawDataQueryUntilParam); err != nil && err != jsonparser.KeyPathNotFoundError {
		return ErrUnmarshalling
	}

	if rq.EstimateSize, err = jsonparser.GetBoolean(data, rawDataQueryEstimateSize); err != nil && err != jsonparser.KeyPathNotFoundError {
		return ErrUnmarshalling
	}

	rq.Tags = map[string]string{}
	err = jsonparser.ObjectEach(data, func(key, value []byte, dataType jsonparser.ValueType, offset int) error {

		tagKey, err := jsonparser.ParseString(key)
		if err != nil {
			return ErrUnmarshalling
		}

		if rq.Tags[tagKey], err = jsonparser.ParseString(value); err != nil {
			return ErrUnmarshalling
		}

		return nil

	}, rawDataQueryTagsParam)

	if _, ok := rq.Tags[rawDataQueryKSID]; !ok {
		return ErrMissingMandatoryFields
	}

	return nil
}
