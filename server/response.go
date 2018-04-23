package server

type Response struct {
	Payload interface{} `json:"payload,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type PageResponse struct {
	Payload interface{} `json:"payload"`
	Limit   int         `json:"limit"`
	Offset  int         `json:"offset"`
	Count   int         `json:"count"`
}
