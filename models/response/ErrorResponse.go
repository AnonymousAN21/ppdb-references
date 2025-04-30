package response

type ErrorResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Error   map[string]string `json:"errors,omitempty"`
}
