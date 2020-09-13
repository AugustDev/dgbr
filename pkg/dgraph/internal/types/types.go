package types

// ExportResponse - response object of export request endpoint
type ExportResponse struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Errors  []errorResponse `json:"errors,omitempty"`
}

type errorResponse struct {
	Message    string `json:"message"`
	Extensions struct {
		Code string `json:"code"`
	} `json:"extensions"`
}
