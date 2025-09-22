package types

// Response represents API response structure
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorCode int

const (
	Success       ErrorCode = 0
	InternalError ErrorCode = 500
	BadRequest    ErrorCode = 400
	NotFound      ErrorCode = 404
)

// APIResponse represents standard API response
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
