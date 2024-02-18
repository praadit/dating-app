package response

import "github.com/praadit/dating-apps/models"

type PaginationResponse struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}
type PackagesResponse struct {
	Pagination PaginationResponse       `json:"pagination"`
	Packages   []models.PackageResponse `json:"packages"`
}
