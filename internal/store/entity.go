package store

import "time"

type Store struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Location     string    `json:"location"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	ImageURL     string    `json:"image_url"`
	WhatsappLink string    `json:"whatsapp_link"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
