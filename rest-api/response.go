package restapi

type Response struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  *Error      `json:"error,omitempty"`
}
