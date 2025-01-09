package category

type CreateInputCategory struct {
	Name string `json:"name" binding:"required"`
}

type GetInputCategory struct {
	Id int `uri:"id" binding:"required"`
}
