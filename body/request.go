package body

import (
	"encoding/json"

	"endpoint-collector/rest/supplier"
)

type RequestBody struct {
	Fields map[string]interface{}
}

func GetEmptyBody() RequestBody {
	return RequestBody{
		Fields: map[string]interface{}{},
	}
}

func (rb *RequestBody) GenerateJSON() ([]byte, error) {
	body := make(map[string]interface{})
	for key, value := range rb.Fields {
		switch v := value.(type) {
		case supplier.Func:
			body[key] = v()
		case supplier.ArrayFunc:
			body[key] = v()
		case *RequestBody:
			nestedJSON, err := v.GenerateMap()
			if err != nil {
				return nil, err
			}
			body[key] = nestedJSON
		default:
			body[key] = value
		}
	}
	return json.Marshal(body)
}

func (rb *RequestBody) GenerateMap() (map[string]interface{}, error) {
	body := make(map[string]interface{})
	for key, value := range rb.Fields {
		switch v := value.(type) {
		case supplier.Func:
			body[key] = v()
		case supplier.ArrayFunc:
			body[key] = v()
		case *RequestBody:
			nestedMap, err := v.GenerateMap()
			if err != nil {
				return nil, err
			}
			body[key] = nestedMap
		default:
			body[key] = value
		}
	}
	return body, nil
}
