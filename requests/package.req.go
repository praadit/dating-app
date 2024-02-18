package requests

type Pagination struct {
	Page    int    `form:"page"`
	PerPage int    `form:"perPage"`
	Order   string `form:"order"`
}

type BuyPackage struct {
	PackageID int `json:"package_id" validate:"required"`
}
