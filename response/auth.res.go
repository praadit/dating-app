package response

import "time"

type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	Expiry      time.Time `json:"expiry"`
	Type        string    `json:"type"`
}

type BaseResponse struct {
	Status bool        `json:"status"`
	Error  *string     `json:"error"`
	Data   interface{} `json:"data"`
}

func FormatRequest(data interface{}, err *string) BaseResponse {
	return BaseResponse{
		Status: err == nil,
		Error:  err,
		Data:   data,
	}
}
