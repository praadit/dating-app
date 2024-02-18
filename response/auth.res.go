package response

import "time"

type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	Expiry      time.Time `json:"expiry"`
	Type        string    `json:"type"`
}

type BaseResponse struct {
	Status bool `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
