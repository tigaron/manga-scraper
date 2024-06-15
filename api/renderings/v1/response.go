package v1Response

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Detail  interface{} `json:"detail,omitempty"`
} // @name ResponseV1
