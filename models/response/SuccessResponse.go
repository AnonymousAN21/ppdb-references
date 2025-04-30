package response

type SuccessResponse struct {
	Success bool `json:"success"`
	Message any  `json:"message"`
	Data    any  `json:"data"`
}

type SuccessResponseGet struct {
	Success    bool       `json:"success"`
	Message    any        `json:"message"`
	Data       any        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
	Has_next     bool   `json:"has_next"`
	Has_previous bool   `json:"has_previous"`
	Next_url     string `json:"next_url"`
	Prev_url     string `json:"prev_url"`
}
