package entities

type BaseResponse struct {
	Message   *string `json:"message"`
	Exception *string `json:"exception"`
}
