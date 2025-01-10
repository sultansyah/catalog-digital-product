package store

import (
	"catalog-digital-product/internal/custom"
	"catalog-digital-product/internal/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StoreHandler interface {
	Update(c *gin.Context)
	GetStore(c *gin.Context)
}

type StoreHandlerImpl struct {
	StoreService StoreService
}

func NewStoreHandler(storeService StoreService) StoreHandler {
	return &StoreHandlerImpl{StoreService: storeService}
}

func (s *StoreHandlerImpl) GetStore(c *gin.Context) {
	store, err := s.StoreService.GetStore(c.Request.Context())
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success get store data",
		Data:    store,
	})
}

func (s *StoreHandlerImpl) Update(c *gin.Context) {
	var input UpdateStoreInput
	if !helper.BindAndValidate(c, &input, "form") {
		return
	}

	imageFile, imageHeader, err := c.Request.FormFile("image")
	if err != nil {
		helper.HandleErrorResponde(c, custom.ErrImageRequired)
		return
	}

	imageFileName := imageHeader.Filename

	store, err := s.StoreService.Update(c.Request.Context(), input, imageFile, imageFileName)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success update store data",
		Data:    store,
	})
}
