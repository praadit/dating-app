package request

type SwipeRequest struct {
	UserId int   `json:"user_id" validate:"required"`
	Liked  *bool `json:"liked" validate:"required"`
}
