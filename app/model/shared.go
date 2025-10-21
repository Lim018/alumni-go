package model

// import "time"

// MetaInfo - Pagination and filtering metadata
type MetaInfo struct {
    Page   int    `json:"page"`
    Limit  int    `json:"limit"`
    Total  int    `json:"total"`
    Pages  int    `json:"pages"`
    SortBy string `json:"sort_by"`
    Order  string `json:"order"`
    Search string `json:"search"`
}

// DatatableRequest - Generic request for datatable endpoints
type DatatableRequest struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
	SortBy string `query:"sortBy"`
	Order  string `query:"order"`
}