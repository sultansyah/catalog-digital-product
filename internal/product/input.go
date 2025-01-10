package product

type CreateProductInput struct {
	CategoryId  int     `form:"category_id"`
	Name        string  `form:"name"`
	Slug        string  `form:"slug"`
	RealPrice   float64 `form:"real_price"`
	Discount    float64 `form:"discount"`
	Stock       int     `form:"stock"`
	Description string  `form:"description"`
}

type GetProductInput struct {
	Id int `uri:"id"`
}

type GetProductImageInput struct {
	Id int `uri:"id"`
}
