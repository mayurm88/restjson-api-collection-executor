package collection

import (
	"fmt"
	"strings"

	"endpoint-collector/rest/body"
	"endpoint-collector/rest/supplier"
)

type EndpointData struct {
	Endpoint              string
	Method                string
	Body                  *body.RequestBody
	AdditionalBodyFields  map[string]supplier.Func
	Inputs                map[string]supplier.Func // map of input url parameter to the value
	QueryParams           map[string]string
	Headers               map[string]string // map of method to headers
	OutputsByResponseKey  map[string]string
	OutputValueByInputKey map[string]string
}

func NewEndpoint() *EndpointData {
	return &EndpointData{}
}

func (e *EndpointData) WithEndpoint(path string) *EndpointData {
	e.Endpoint = path
	return e
}

func (e *EndpointData) WithMethod(m string) *EndpointData {
	e.Method = m
	return e
}

func (e *EndpointData) WithRequestBody(b *body.RequestBody) *EndpointData {
	e.Body = b
	return e
}

func (e *EndpointData) WithAdditionalBodyFields(key string, value supplier.Func) *EndpointData {
	if e.AdditionalBodyFields == nil {
		e.AdditionalBodyFields = make(map[string]supplier.Func)
	}
	e.AdditionalBodyFields[key] = value
	return e
}

func (e *EndpointData) SaveRegisteredOutputs(resp ResponseBody) error {
	if e.OutputValueByInputKey == nil {
		e.OutputValueByInputKey = make(map[string]string)
	}

	if len(e.OutputsByResponseKey) == 0 {
		return nil
	}

	for respKey, outputKey := range e.OutputsByResponseKey {
		jsonPath := strings.Split(respKey, ".")
		var curValue interface{} = resp.Fields
		for idx, curKey := range jsonPath {
			curMap, ok := curValue.(map[string]interface{})
			if !ok {
				if idx != len(jsonPath)-1 {
					return fmt.Errorf("unexpected type while walking json, check your registered output %s", respKey)
				}
				curString, ok := curValue.(string)
				if !ok {
					return fmt.Errorf("output value is an unsupported type, needs to be a string outputValue: %v", curValue)
				}
				e.OutputValueByInputKey[outputKey] = curString
				return nil
			}
			curValue, ok = curMap[curKey]
			if !ok {
				return fmt.Errorf("key %s not found while traversing the registered output %s", curKey, respKey)
			}

			// If it's the last key in the path, return the value
			if idx == len(jsonPath)-1 {
				curString, ok := curValue.(string)
				if !ok {
					return fmt.Errorf("output value is an unsupported type, needs to be a string outputValue: %v", curValue)
				}
				e.OutputValueByInputKey[outputKey] = curString
				return nil
			}
		}
	}

	return fmt.Errorf("unknown state, boom")
}

func (e *EndpointData) RegisterOutput(responseKey, outputKey string) *EndpointData {
	if e.OutputsByResponseKey == nil {
		e.OutputsByResponseKey = make(map[string]string)
	}
	e.OutputsByResponseKey[responseKey] = outputKey
	return e
}

func (e *EndpointData) PutOutput(key, value string) {
	if e.OutputValueByInputKey == nil {
		e.OutputValueByInputKey = make(map[string]string)
	}
	e.OutputValueByInputKey[key] = value
}

func (e *EndpointData) GetSupplierFunc(key string) supplier.Func {
	return func() string {
		value, ok := e.OutputValueByInputKey[key]
		if !ok {
			return ""
		}
		return value
	}
}

func (e *EndpointData) WithInput(key string, s supplier.Func) *EndpointData {
	if e.Inputs == nil {
		e.Inputs = make(map[string]supplier.Func)
	}
	e.Inputs[key] = s
	return e
}

func (e *EndpointData) WithQueryParams(key string, value string) *EndpointData {
	if e.QueryParams == nil {
		e.QueryParams = make(map[string]string)
	}
	e.QueryParams[key] = value
	return e
}

func (e *EndpointData) WithHeaders(key string, value string) *EndpointData {
	if e.Headers == nil {
		e.Headers = make(map[string]string)
	}
	e.Headers[key] = value
	return e
}
