package collection

type EndpointResult struct {
	ResponseBody    []byte
	ResponseHeaders map[string]string
	Status          int
}

type Result struct {
	resultByEndpoint map[*EndpointData]EndpointResult
}
