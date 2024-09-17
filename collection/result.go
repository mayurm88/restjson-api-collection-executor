package collection

type EndpointResult struct {
	ResponseBody    []byte
	ResponseHeaders map[string]string
	Status          int
}

type Result struct {
	resultByEndpoint map[*EndpointData]EndpointResult
}

func NewResult() Result {
	return Result{
		resultByEndpoint: map[*EndpointData]EndpointResult{},
	}
}
