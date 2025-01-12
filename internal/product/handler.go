package product

import (
	"catalog-digital-product/internal/custom"
	"catalog-digital-product/internal/helper"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler interface {
	Insert(c *gin.Context)
	Update(c *gin.Context)
	Get(c *gin.Context)
	GetAll(c *gin.Context)
	Delete(c *gin.Context)
	InsertImage(c *gin.Context)
	DeleteImage(c *gin.Context)
	SetLogoImage(c *gin.Context)
}

type ProductHandlerImpl struct {
	ProductService ProductService
}

func NewProductHandler(productService ProductService) ProductHandler {
	return &ProductHandlerImpl{ProductService: productService}
}

func (p *ProductHandlerImpl) Delete(c *gin.Context) {
	var input GetProductInput
	if !helper.BindAndValidate(c, &input, "uri") {
		return
	}

	err := p.ProductService.Delete(c.Request.Context(), input)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success delete product",
		Data:    "OK",
	})
}

func (p *ProductHandlerImpl) DeleteImage(c *gin.Context) {
	var input GetProductImageInput
	if !helper.BindAndValidate(c, &input, "uri") {
		return
	}

	err := p.ProductService.DeleteImage(c.Request.Context(), input)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success delete image",
		Data:    "OK",
	})
}

func (p *ProductHandlerImpl) Get(c *gin.Context) {
	var input GetProductInput
	if !helper.BindAndValidate(c, &input, "uri") {
		return
	}

	product, err := p.ProductService.Get(c.Request.Context(), input)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success get data",
		Data:    product,
	})
}

func (p *ProductHandlerImpl) GetAll(c *gin.Context) {
	products, err := p.ProductService.GetAll(c.Request.Context())
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success get all data product",
		Data:    products,
	})
}

func (p *ProductHandlerImpl) Insert(c *gin.Context) {
	var input CreateProductInput
	if !helper.BindAndValidate(c, &input, "form") {
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	imageFiles := form.File["images"]
	if len(imageFiles) <= 0 {
		helper.HandleErrorResponde(c, custom.ErrImageRequired)
		return
	}

	productImagesFile := make(map[string]multipart.File)
	for _, file := range imageFiles {
		tempFile, err := file.Open()
		if err != nil {
			helper.HandleErrorResponde(c, err)
			return
		}

		productImagesFile[file.Filename] = tempFile
		defer tempFile.Close()
	}

	product, err := p.ProductService.Create(c.Request.Context(), input, productImagesFile)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success insert new product",
		Data:    product,
	})
}

func (p *ProductHandlerImpl) InsertImage(c *gin.Context) {
	var input GetProductInput
	if !helper.BindAndValidate(c, &input, "uri") {
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	imageFiles := form.File["images"]

	productImagesFile := make(map[string]multipart.File)
	for _, file := range imageFiles {
		tempFile, err := file.Open()
		if err != nil {
			helper.HandleErrorResponde(c, err)
			return
		}

		productImagesFile[file.Filename] = tempFile
		defer tempFile.Close()
	}

	err = p.ProductService.CreateImage(c.Request.Context(), input, productImagesFile)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success insert new product",
		Data:    "OK",
	})
}

func (p *ProductHandlerImpl) SetLogoImage(c *gin.Context) {
	var inputProductId GetProductInput
	if !helper.BindAndValidate(c, &inputProductId, "uri") {
		return
	}

	var inputProductImageId GetProductImageInput
	if !helper.BindAndValidate(c, &inputProductImageId, "uri") {
		return
	}

	err := p.ProductService.SetLogo(c.Request.Context(), inputProductId, inputProductImageId)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success set logo",
		Data:    "OK",
	})
}

func (p *ProductHandlerImpl) Update(c *gin.Context) {
	var inputData UpdateProductInput
	if !helper.BindAndValidate(c, &inputData, "json") {
		return
	}

	var inputId GetProductInput
	if !helper.BindAndValidate(c, &inputId, "uri") {
		return
	}

	product, err := p.ProductService.Update(c.Request.Context(), inputData, inputId)
	if err != nil {
		helper.HandleErrorResponde(c, err)
		return
	}

	helper.APIResponse(c, helper.WebResponse{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "success update product",
		Data:    product,
	})
}
