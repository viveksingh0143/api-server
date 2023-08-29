package valueobjects

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	TotalItems int         `json:"total_items"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}
