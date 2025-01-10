package product

import (
	"catalog-digital-product/internal/custom"
	"catalog-digital-product/internal/helper"
	"context"
	"database/sql"
	"mime/multipart"
	"sync"
)

type ProductService interface {
	Create(ctx context.Context, input CreateProductInput, productImagesFile map[string]multipart.File) (Product, error)
	Update(ctx context.Context, inputData CreateProductInput, inputId GetProductInput) (Product, error)
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

	product := Product{
		CategoryId:  input.CategoryId,
		Name:        input.Name,
		Slug:        input.Slug,
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

			insertProductImages(*p, fileName, file, errChan, &productImage, ctx, tx)
		}()
	}

	wg.Wait()

	close(errChan)
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
		go insertProductImages(*p, fileName, file, errChan, &productImage, ctx, tx)
	}
	wg.Wait()

	close(errChan)
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

	return nil
}

func (p *ProductServiceImpl) DeleteImage(ctx context.Context, input GetProductImageInput) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.HandleTransaction(tx, &err)

	product, err := p.ProductRepository.FindImageById(ctx, tx, input.Id)
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
		return []Product{}, custom.ErrNotFound
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

func (p *ProductServiceImpl) Update(ctx context.Context, inputData CreateProductInput, inputId GetProductInput) (Product, error) {
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

	product.CategoryId = inputData.CategoryId
	product.Name = inputData.Name
	product.Slug = inputData.Slug
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
