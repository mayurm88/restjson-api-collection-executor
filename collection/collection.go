package collection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/mayurm88/restjson-api-collection-executor/supplier"
)

type Collection struct {
	Endpoints    []*EndpointData
	CommonInputs map[string]supplier.Func
	BaseURL      string
	BearerToken  string
}

type ResponseBody struct {
	Fields map[string]interface{} `json:"fields"`
}

var bodyProcessMethods = map[string]bool{
	http.MethodPut:   true,
	http.MethodPost:  true,
	http.MethodPatch: true,
}

func NewCollection() *Collection {
	return &Collection{}
}

func (c *Collection) WithEndpoint(data *EndpointData) *Collection {
	if c.Endpoints == nil {
		c.Endpoints = make([]*EndpointData, 0)
	}
	c.Endpoints = append(c.Endpoints, data)
	return c
}

func (c *Collection) WithCommonInput(key string, s supplier.Func) *Collection {
	if c.CommonInputs == nil {
		c.CommonInputs = make(map[string]supplier.Func)
	}
	c.CommonInputs[key] = s
	return c
}

func (c *Collection) WithBaseURL(b string) *Collection {
	c.BaseURL = b
	return c
}

func (c *Collection) WithBearerToken(b string) *Collection {
	c.BearerToken = b
	return c
}

func (c *Collection) Execute() (*Result, error) {

	client := http.Client{
		Transport: http.DefaultTransport,
		Timeout:   time.Second * 10,
	}

	result := NewResult()

	for _, e := range c.Endpoints {
		url, err := c.constructURL(e)
		if err != nil {
			return nil, err
		}

		var reqBodyBytes []byte
		if e.Body != nil {
			reqBodyBytes, err = e.Body.GenerateJSON()
			if err != nil {
				return nil, err
			}
		}
		var bodyReader *bytes.Reader
		var req *http.Request
		if e.Method == http.MethodPost || e.Method == http.MethodPut {
			bodyReader = bytes.NewReader(reqBodyBytes)
			req, err = http.NewRequest(e.Method, url, bodyReader)
		} else {
			req, err = http.NewRequest(e.Method, url, nil)
		}

		for k, v := range e.Headers {
			req.Header.Set(k, v)
		}

		if len(c.BearerToken) > 0 {
			req.Header.Set("Authorization", "Bearer "+c.BearerToken)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("", err)
			return nil, err
		}

		defer resp.Body.Close()

		log.Printf("response for endpoint %s and Method %s received status %s\n", url, e.Method, resp.Status)
		if !bodyProcessMethods[e.Method] {
			err = addResponseToResult(resp, &result, e)
			if err != nil {
				return nil, err
			}
			continue
		}

		statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
		if !statusOK {
			err = addResponseToResult(resp, &result, e)
			if err != nil {
				return nil, err
			}
			continue
		}

		err = addResponseToResult(resp, &result, e)
		if err != nil {
			return nil, err
		}

		respBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		var respBody ResponseBody
		err = json.Unmarshal(respBodyBytes, &respBody.Fields)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		err = e.SaveRegisteredOutputs(respBody)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func addResponseToResult(resp *http.Response, result *Result, endpoint *EndpointData) error {
	headers := resp.Header
	var resultHeaders map[string]string
	for k, _ := range headers {
		value := headers.Get(k)
		resultHeaders[k] = value
	}
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		log.Printf("request failed with status %d", resp.StatusCode)
		result.resultByEndpoint[endpoint] = EndpointResult{
			ResponseHeaders: resultHeaders,
			Status:          resp.StatusCode,
		}
		return nil
	}

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}

	result.resultByEndpoint[endpoint] = EndpointResult{
		ResponseHeaders: resultHeaders,
		Status:          resp.StatusCode,
		ResponseBody:    respBodyBytes,
	}
	return nil
}

func (c *Collection) constructURL(e *EndpointData) (string, error) {
	keys := c.extractURLParameters(e.Endpoint)
	url := c.BaseURL + "/" + c.normalizePath(e.Endpoint)
	for _, key := range keys {
		replaceValue, ok := c.CommonInputs[key]
		if !ok {
			replaceValue, ok = e.Inputs[key]
			if !ok {
				return "", fmt.Errorf("couldn't find key %s in inputs", key)
			}
		}
		url = strings.Replace(url, "{"+key+"}", replaceValue(), -1)
	}
	return url, nil
}

func (c *Collection) extractURLParameters(url string) []string {
	// Define a regular expression pattern to find parameters in the format {paramName}
	pattern := `\{([^}]+)\}`
	re := regexp.MustCompile(pattern)

	// Find all matches in the URL
	matches := re.FindAllStringSubmatch(url, -1)

	// Extract the parameter names from the matches
	var params []string
	for _, match := range matches {
		if len(match) > 1 {
			params = append(params, match[1])
		}
	}

	return params
}

func (c *Collection) normalizePath(path string) string {
	if strings.HasPrefix(path, "/") {
		return path[1:]
	}
	return path
}
