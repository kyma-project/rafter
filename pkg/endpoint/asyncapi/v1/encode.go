package v1

import (
	"errors"
	"io"

	"github.com/asyncapi/converter-go/pkg/encode"
)

var errInvalidType = errors.New("invalid type")

// removes null entries from JSON array
func sanitizeSlice(data []interface{}) interface{} {
	result := []interface{}{}
	for _, item := range data {
		if item == nil {
			continue
		}

		switch v := item.(type) {
		case []interface{}:
			if fixedSlice := sanitizeSlice(v); fixedSlice != nil {
				result = append(result, fixedSlice)
			}
		case map[string]interface{}:
			if fixedMap := sanitizeMap(v); fixedMap != nil {
				result = append(result, fixedMap)
			}
		default:
			result = append(result, item)
		}
	}

	return result
}

// removes null entries from JSON object
func sanitizeMap(data map[string]interface{}) interface{} {
	for k, v := range data {
		if v == nil {
			delete(data, k)
		}

		switch v := v.(type) {
		case []interface{}:
			if fixedSlice := sanitizeSlice(v); fixedSlice != nil {
				data[k] = fixedSlice
			}
		case map[string]interface{}:
			if fixedMap := sanitizeMap(v); fixedMap != nil {
				data[k] = fixedMap
			}
		}
	}

	return data
}

var errInvlidDataType = errors.New("invalid data type")

func defaultJSONEncoder(i interface{}, w io.Writer) error {
	data, ok := i.(*map[string]interface{})
	if !ok {
		return errInvalidType
	}

	sanitizedData := sanitizeMap(*data)
	return encode.ToJSON(&sanitizedData, w)
}
