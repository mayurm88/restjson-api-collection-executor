package collection

import (
	"net/http"

	"github.com/mayurm88/restjson-api-collection-executor/supplier"
)

func GetRestCreateEndpointWithURL(url string) *EndpointData {
	return NewEndpoint().
		WithEndpoint(url).
		WithHeaders("Content-Type", "application/json").
		WithMethod(http.MethodPost)
}

func GetRestListEndpointWithURL(url string) *EndpointData {
	return NewEndpoint().
		WithEndpoint(url).
		WithMethod(http.MethodGet)
}

func GetRestReadEndpointWithURL(url, inputKey string, supplierFunc supplier.Func) *EndpointData {
	return NewEndpoint().
		WithEndpoint(url).
		WithMethod(http.MethodGet).
		WithInput(inputKey, supplierFunc)
}

func GetRestUpdateEndpointWithURL(url, inputKey string, supplier supplier.Func) *EndpointData {
	return NewEndpoint().
		WithEndpoint(url).
		WithHeaders("Content-Type", "application/json").
		WithMethod(http.MethodPut).
		WithInput(inputKey, supplier)
}

func GetRestDeleteEndpointWithURL(url, inputKey string, supplier supplier.Func) *EndpointData {
	return NewEndpoint().
		WithEndpoint(url).
		WithMethod(http.MethodDelete).
		WithInput(inputKey, supplier)
}
