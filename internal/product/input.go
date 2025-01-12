package product

type CreateProductInput struct {
	CategoryId  int     `form:"category_id" binding:"required"`
	Name        string  `form:"name" binding:"required"`
	RealPrice   float64 `form:"real_price" binding:"required"`
	Discount    float64 `form:"discount" binding:"required"`
	Stock       int     `form:"stock" binding:"required"`
	Description string  `form:"description" binding:"required"`
}

type UpdateProductInput struct {
	CategoryId  int     `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	RealPrice   float64 `json:"real_price" binding:"required"`
	Discount    float64 `json:"discount" binding:"required"`
	Stock       int     `json:"stock" binding:"required"`
	Description string  `json:"description" binding:"required"`
}

type GetProductInput struct {
	Id int `uri:"id" binding:"required"`
}

type GetProductImageInput struct {
	Id int `uri:"imageId" binding:"required"`
}
