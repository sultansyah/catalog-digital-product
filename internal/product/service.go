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
	Update(ctx context.Context, input CreateProductInput) (Product, error)
	Delete(ctx context.Context, input GetProductInput) error
	Get(ctx context.Context, input GetProductInput) (Product, error)
	GetAll(ctx context.Context) ([]Product, error)
	CreateImage(ctx context.Context, input GetProductInput, productImagesFile map[string]multipart.File) error
	UpdateImageSetLogo(ctx context.Context, input CreateProductInput) ([]ProductImages, error)
	DeleteImage(ctx context.Context, input GetProductInput) error
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
	panic("unimplemented")
}

func (p *ProductServiceImpl) DeleteImage(ctx context.Context, input GetProductInput) error {
	panic("unimplemented")
}

func (p *ProductServiceImpl) Get(ctx context.Context, input GetProductInput) (Product, error) {
	panic("unimplemented")
}

func (p *ProductServiceImpl) GetAll(ctx context.Context) ([]Product, error) {
	panic("unimplemented")
}

func (p *ProductServiceImpl) SetLogo(ctx context.Context, inputProductId GetProductInput, inputProductImageId GetProductImageInput) error {
	panic("unimplemented")
}

func (p *ProductServiceImpl) Update(ctx context.Context, input CreateProductInput) (Product, error) {
	panic("unimplemented")
}

func (p *ProductServiceImpl) UpdateImageSetLogo(ctx context.Context, input CreateProductInput) (ProductImages, error) {
	panic("unimplemented")
}
