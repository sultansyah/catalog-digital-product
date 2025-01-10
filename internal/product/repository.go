package product

import (
	"catalog-digital-product/internal/custom"
	"context"
	"database/sql"
)

type ProductRepository interface {
	Insert(ctx context.Context, tx *sql.Tx, product Product) (Product, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (Product, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]Product, error)
	Update(ctx context.Context, tx *sql.Tx, product Product) (Product, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	InsertImage(ctx context.Context, tx *sql.Tx, productImages ProductImages) error
	DeleteImage(ctx context.Context, tx *sql.Tx, id int) error
	UpdateImage(ctx context.Context, tx *sql.Tx, id int) error
	UpdateImagesLogoFalse(ctx context.Context, tx *sql.Tx, productId int) error
}

type ProductRepositoryImpl struct {
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}

func (p *ProductRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	sql := "delete from product where id = ?"
	_, err := tx.ExecContext(ctx, sql, id)
	return err
}

func (p *ProductRepositoryImpl) DeleteImage(ctx context.Context, tx *sql.Tx, id int) error {
	sql := "delete from product_images where id = ?"
	_, err := tx.ExecContext(ctx, sql, id)
	return err
}

func (p *ProductRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]Product, error) {
	sql := `
			SELECT 
				p.id,
				p.category_id, 
				p.name, 
				p.slug, 
				p.real_price, 
				p.discount, 
				p.stock, 
				p.description, 
				p.created_at, 
				p.updated_at,
				c.id AS category_id, 
				c.name AS category_name, 
				pi.id AS image_id, 
				pi.image_url, 
				pi.is_logo
			FROM products AS p
			INNER JOIN categories AS c ON c.id = p.category_id
			INNER JOIN product_images AS pi ON pi.product_id = p.id
		`
	rows, err := tx.QueryContext(ctx, sql)
	if err != nil {
		return []Product{}, err
	}

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.Id, &product.CategoryId, &product.Name, &product.Slug, &product.RealPrice, &product.Discount, &product.Stock, &product.Description, &product.CreatedAt, &product.UpdatedAt, &product.Category.Id, &product.Category.Name, &product.ProductImages.Id, &product.ProductImages.ImageUrl, &product.ProductImages.IsLogo); err != nil {
			return []Product{}, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (p *ProductRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (Product, error) {
	sql := `
			SELECT 
				p.id,
				p.category_id, 
				p.name, 
				p.slug, 
				p.real_price, 
				p.discount, 
				p.stock, 
				p.description, 
				p.created_at, 
				p.updated_at,
				c.id AS category_id, 
				c.name AS category_name, 
				pi.id AS image_id, 
				pi.image_url, 
				pi.is_logo
			FROM products AS p
			INNER JOIN categories AS c ON c.id = p.category_id
			INNER JOIN product_images AS pi ON pi.product_id = p.id
			WHERE p.id = ?
		`
	row, err := tx.QueryContext(ctx, sql, id)
	if err != nil {
		return Product{}, err
	}

	var product Product
	if row.Next() {
		if err := row.Scan(&product.Id, &product.CategoryId, &product.Name, &product.Slug, &product.RealPrice, &product.Discount, &product.Stock, &product.Description, &product.CreatedAt, &product.UpdatedAt, &product.Category.Id, &product.Category.Name, &product.ProductImages.Id, &product.ProductImages.ImageUrl, &product.ProductImages.IsLogo); err != nil {
			return Product{}, err
		}

		return product, nil
	}

	return product, custom.ErrNotFound
}

func (p *ProductRepositoryImpl) Insert(ctx context.Context, tx *sql.Tx, product Product) (Product, error) {
	sql := `
	INSERT INTO products
	(category_id, name, slug, real_price, discount, stock, description) 
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := tx.ExecContext(ctx, sql, product.CategoryId, product.Name, product.Slug, product.RealPrice, product.Discount, product.Stock, product.Description)
	if err != nil {
		return Product{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Product{}, err
	}

	product.Id = int(id)
	return product, nil
}

func (p *ProductRepositoryImpl) InsertImage(ctx context.Context, tx *sql.Tx, productImages ProductImages) error {
	sql := `
	INSERT INTO product_images
	(product_id, image_url, is_logo) 
	VALUES (?, ?, ?)
	`
	result, err := tx.ExecContext(ctx, sql, productImages.ProductId, productImages.ImageUrl, productImages.IsLogo)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	productImages.Id = int(id)
	return nil
}

func (p *ProductRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, product Product) (Product, error) {
	sql := `
	UPDATE products 
	SET category_id=?,name=?,slug=?,real_price=?,discount=?,stock=?,description=? WHERE id = ?
	`
	_, err := tx.ExecContext(ctx, sql, product.CategoryId, product.Name, product.Slug, product.RealPrice, product.Discount, product.Stock, product.Description, product.Id)
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func (p *ProductRepositoryImpl) UpdateImage(ctx context.Context, tx *sql.Tx, id int) error {
	sql := `
	UPDATE product_images 
	SET is_logo = ? WHERE id = ?
	`
	_, err := tx.ExecContext(ctx, sql, id)
	return err
}

func (p *ProductRepositoryImpl) UpdateImagesLogoFalse(ctx context.Context, tx *sql.Tx, productId int) error {
	sql := `
	UPDATE product_images 
	SET is_logo = FALSE WHERE product_id = ?
	`
	_, err := tx.ExecContext(ctx, sql, productId)
	return err
}
