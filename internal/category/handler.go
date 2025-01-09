package category

import (
	"catalog-digital-product/internal/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type CategoryHandlerImpl struct {
	CategoryService CategoryService
}

func NewCategoryHandler(categoryService CategoryService) CategoryHandler {
	return &CategoryHandlerImpl{CategoryService: categoryService}
}

func (h *CategoryHandlerImpl) GetAll(c *gin.Context) {
	categories, err := h.CategoryService.GetAll(c.Request.Context())
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success get categories",
		Data:    categories,
	})
}

func (h *CategoryHandlerImpl) Update(c *gin.Context) {
	var inputId GetInputCategory
	if !helper.BindAndValidate(c, &inputId, "uri") {
		return
	}

	var inputData CreateInputCategory
	if !helper.BindAndValidate(c, &inputData, "json") {
		return
	}

	category, err := h.CategoryService.Update(c.Request.Context(), inputData, inputId)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success get category",
		Data:    category,
	})
}

func (h *CategoryHandlerImpl) Create(c *gin.Context) {
	var input CreateInputCategory
	if !helper.BindAndValidate(c, &input, "json") {
		return
	}

	category, err := h.CategoryService.Create(c.Request.Context(), input)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success create category",
		Data:    category,
	})
}

func (h *CategoryHandlerImpl) Delete(c *gin.Context) {
	var input GetInputCategory
	if !helper.BindAndValidate(c, &input, "uri") {
		return
	}

	err := h.CategoryService.Delete(c.Request.Context(), input)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success get category",
		Data:    "OK",
	})
}

func (h *CategoryHandlerImpl) Get(c *gin.Context) {
	var input GetInputCategory
	if !helper.BindAndValidate(c, &input, "uri") {
		return
	}

	category, err := h.CategoryService.Get(c.Request.Context(), input)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success get category",
		Data:    category,
	})
}
