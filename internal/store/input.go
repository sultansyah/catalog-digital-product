package store

type UpdateStoreInput struct {
	Name         string  `form:"name" binding:"required"`
	Description  string  `form:"description" binding:"required"`
	Location     string  `form:"location" binding:"required"`
	Latitude     float64 `form:"latitude" binding:"required"`
	Longitude    float64 `form:"longitude" binding:"required"`
	PhoneNumber  string  `form:"phone_number" binding:"required"`
	Email        string  `form:"email" binding:"required"`
	WhatsappLink string  `form:"whatsapp_link" binding:"required"`
}
