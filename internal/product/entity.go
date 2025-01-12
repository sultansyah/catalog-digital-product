package product

import (
	"catalog-digital-product/internal/category"
	"time"
)

type Product struct {
	Id            int               `json:"id"`
	CategoryId    int               `json:"category_id"`
	Category      category.Category `json:"category"`
	ProductImages []ProductImages   `json:"product_images"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	RealPrice     float64           `json:"real_price"`
	Discount      float64           `json:"discount"`
	Stock         int               `json:"stock"`
	Description   string            `json:"description"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type ProductImages struct {
	Id        int       `json:"id"`
	ProductId int       `json:"product_id"`
	ImageUrl  string    `json:"image_url"`
	IsLogo    bool      `json:"is_logo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
