package entities

type User struct {
	Id                   string `json:"id" binding:"required"`
	Email                string `json:"email" binding:"required"`
	FirstName            string `json:"first_name" binding:"required"`
	LastName             string `json:"last_name" binding:"required"`
	PhotoUrl             string `json:"photo_url"`
	HexBytesPhotoPreview string `json:"photo_preview_hex"`
}
