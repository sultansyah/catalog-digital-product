package product

import (
	"catalog-digital-product/internal/custom"
	"catalog-digital-product/internal/helper"
	"context"
	"database/sql"
	"mime/multipart"
	"sync"

	"github.com/gosimple/slug"
)

type ProductService interface {
	Create(ctx context.Context, input CreateProductInput, productImagesFile map[string]multipart.File) (Product, error)
	Update(ctx context.Context, inputData UpdateProductInput, inputId GetProductInput) (Product, error)
	Delete(ctx context.Context, input GetProductInput) error
	Get(ctx context.Context, input GetProductInput) (Product, error)
	GetAll(ctx context.Context) ([]Product, error)
	CreateImage(ctx context.Context, input GetProductInput, productImagesFile map[string]multipart.File) error
	DeleteImage(ctx context.Context, input GetProductImageInput) error
	SetLogo(ctx context.Context, inputProductId GetProductInput, inputProductImageId GetProductImageInput) error
}

type ProductServiceImpl struct {
	ProductRepository ProductRepository
	DB                *sql.DB
}

func NewProductService(productRepository ProductRepository, DB *sql.DB) ProductService {
	return &ProductServiceImpl{ProductRepository: productRepository, DB: DB}
}

func (p *ProductServiceImpl) Create(ctx context.Context, input CreateProductInput, productImagesFile map[string]multipart.File) (Product, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return Product{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	isProductExist, err := p.ProductRepository.FindBySlug(ctx, tx, slug.Make(input.Name))
	if err != nil && err != custom.ErrNotFound {
		return Product{}, err
	}
	if isProductExist.Id > 0 {
		return Product{}, custom.ErrAlreadyExists
	}

	product := Product{
		CategoryId:  input.CategoryId,
		Name:        input.Name,
		Slug:        slug.Make(input.Name),
		RealPrice:   input.RealPrice,
		Discount:    input.Discount,
		Stock:       input.Stock,
		Description: input.Description,
	}
	product, err = p.ProductRepository.Insert(ctx, tx, product)
	if err != nil {
		return Product{}, err
	}

	productImage := ProductImages{
		ProductId: product.Id,
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	isFirstFile := true

	errChan := make(chan error, len(productImagesFile))

	for fileName, file := range productImagesFile {
		wg.Add(1)
		go func() {
			defer wg.Done()

			mu.Lock()
			if isFirstFile {
				productImage.IsLogo = true
				isFirstFile = false
			}
			mu.Unlock()

			insertProductImages(*p, "product", fileName, file, errChan, &productImage, ctx, tx)
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return Product{}, err
		}
	}

	return product, nil
}

func (p *ProductServiceImpl) CreateImage(ctx context.Context, input GetProductInput, productImagesFile map[string]multipart.File) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.HandleTransaction(tx, &err)

	product, err := p.ProductRepository.FindById(ctx, tx, input.Id)
	if err != nil {
		return err
	}
	if product.Id <= 0 {
		return custom.ErrNotFound
	}

	var wg sync.WaitGroup
	errChan := make(chan error)
	productImage := ProductImages{
		ProductId: input.Id,
		IsLogo:    false,
	}

	for fileName, file := range productImagesFile {
		wg.Add(1)
		go func(fileName string, file multipart.File) {
			defer wg.Done()
			insertProductImages(*p, "product", fileName, file, errChan, &productImage, ctx, tx)
		}(fileName, file)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ProductServiceImpl) Delete(ctx context.Context, input GetProductInput) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.HandleTransaction(tx, &err)

	product, err := p.ProductRepository.FindById(ctx, tx, input.Id)
	if err != nil {
		return err
	}
	if product.Id <= 0 {
		return custom.ErrNotFound
	}

	err = p.ProductRepository.Delete(ctx, tx, input.Id)
	if err != nil {
		return err
	}

	chanErr := make(chan error, len(product.ProductImages))
	var wg sync.WaitGroup

	for _, image := range product.ProductImages {
		wg.Add(1)
		go func(image ProductImages) {
			defer wg.Done()
			err := deleteImage("product", image.ImageUrl)
			chanErr <- err
		}(image)
	}

	go func() {
		wg.Wait()
		close(chanErr)
	}()

	for err := range chanErr {
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ProductServiceImpl) DeleteImage(ctx context.Context, input GetProductImageInput) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.HandleTransaction(tx, &err)

	productImage, err := p.ProductRepository.FindImageById(ctx, tx, input.Id)
	if err != nil {
		return err
	}
	if productImage.Id <= 0 {
		return custom.ErrNotFound
	}

	err = p.ProductRepository.Delete(ctx, tx, input.Id)
	if err != nil {
		return err
	}

	err = deleteImage("product", productImage.ImageUrl)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductServiceImpl) Get(ctx context.Context, input GetProductInput) (Product, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return Product{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	product, err := p.ProductRepository.FindById(ctx, tx, input.Id)
	if err != nil {
		return Product{}, err
	}
	if product.Id <= 0 {
		return Product{}, custom.ErrNotFound
	}

	return product, nil
}

func (p *ProductServiceImpl) GetAll(ctx context.Context) ([]Product, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return []Product{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	products, err := p.ProductRepository.FindAll(ctx, tx)
	if err != nil {
		return []Product{}, err
	}
	if len(products) <= 0 {
		return []Product{}, nil
	}

	return products, nil
}

func (p *ProductServiceImpl) SetLogo(ctx context.Context, inputProductId GetProductInput, inputProductImageId GetProductImageInput) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.HandleTransaction(tx, &err)

	product, err := p.ProductRepository.FindById(ctx, tx, inputProductId.Id)
	if err != nil {
		return err
	}
	if product.Id <= 0 {
		return custom.ErrNotFound
	}

	productImages, err := p.ProductRepository.FindImageById(ctx, tx, inputProductImageId.Id)
	if err != nil {
		return err
	}
	if productImages.Id <= 0 {
		return custom.ErrNotFound
	}

	err = p.ProductRepository.UpdateImagesLogoFalse(ctx, tx, inputProductId.Id)
	if err != nil {
		return err
	}

	err = p.ProductRepository.UpdateImage(ctx, tx, inputProductImageId.Id)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductServiceImpl) Update(ctx context.Context, inputData UpdateProductInput, inputId GetProductInput) (Product, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return Product{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	product, err := p.ProductRepository.FindById(ctx, tx, inputId.Id)
	if err != nil {
		return Product{}, err
	}
	if product.Id <= 0 {
		return Product{}, custom.ErrNotFound
	}

	isSlugExist, err := p.ProductRepository.FindBySlug(ctx, tx, slug.Make(inputData.Name))
	if err != nil && err != custom.ErrNotFound {
		return Product{}, err
	}
	if isSlugExist.Id > 0 && isSlugExist.Id != product.Id {
		return Product{}, custom.ErrAlreadyExists
	}

	product.CategoryId = inputData.CategoryId
	product.Name = inputData.Name
	product.Slug = slug.Make(inputData.Name)
	product.RealPrice = inputData.RealPrice
	product.Discount = inputData.Discount
	product.Stock = inputData.Stock
	product.Description = inputData.Description

	product, err = p.ProductRepository.Update(ctx, tx, product)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}
