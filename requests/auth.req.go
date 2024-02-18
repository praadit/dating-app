package requests

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignupRequest struct {
	Name            string `json:"name" validate:"required,unsafe,max=200"`
	Email           string `json:"email" validate:"required,email,max=100,lowercase"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	Gender          string `json:"gender" validate:"required,max=1"`
	Picture         string `json:"picture"`
}
