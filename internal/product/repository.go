package product

import (
	"catalog-digital-product/internal/custom"
	"context"
	"database/sql"
)

type ProductRepository interface {
	Insert(ctx context.Context, tx *sql.Tx, product Product) (Product, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (Product, error)
	FindLatestImage(ctx context.Context, tx *sql.Tx, productId int) (ProductImages, error)
	FindBySlug(ctx context.Context, tx *sql.Tx, slug string) (Product, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]Product, error)
	Update(ctx context.Context, tx *sql.Tx, product Product) (Product, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	InsertImage(ctx context.Context, tx *sql.Tx, productImages ProductImages) error
	DeleteImage(ctx context.Context, tx *sql.Tx, id int) error
	UpdateImage(ctx context.Context, tx *sql.Tx, id int) error
	UpdateImagesLogoFalse(ctx context.Context, tx *sql.Tx, productId int) error
	FindImageById(ctx context.Context, tx *sql.Tx, id int) (ProductImages, error)
}

type ProductRepositoryImpl struct {
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}

func (p *ProductRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	sql := "delete from products where id = ?"
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
				c.created_at, 
				c.updated_at,
				pi.id AS image_id, 
				pi.product_id,
				pi.image_url,
				pi.is_logo,
				pi.created_at, 
				pi.updated_at
			FROM products AS p
			INNER JOIN categories AS c ON c.id = p.category_id
			INNER JOIN product_images AS pi ON pi.product_id = p.id
		`
	rows, err := tx.QueryContext(ctx, sql)
	if err != nil {
		return []Product{}, err
	}
	defer rows.Close()

	productMap := make(map[int]Product)

	for rows.Next() {
		var product Product
		var image ProductImages

		if err := rows.Scan(
			&product.Id, &product.CategoryId, &product.Name, &product.Slug, &product.RealPrice,
			&product.Discount, &product.Stock, &product.Description, &product.CreatedAt, &product.UpdatedAt,
			&product.Category.Id, &product.Category.Name, &product.Category.CreatedAt,
			&product.Category.UpdatedAt, &image.Id, &image.ProductId,
			&image.ImageUrl, &image.IsLogo, &image.CreatedAt, &image.UpdatedAt); err != nil {
			return []Product{}, err
		}

		if existingProduct, exist := productMap[product.Id]; exist {
			existingProduct.ProductImages = append(existingProduct.ProductImages, image)
			productMap[product.Id] = existingProduct
		} else {
			product.ProductImages = append(product.ProductImages, image)
			productMap[product.Id] = product
		}
	}

	var products []Product
	for _, product := range productMap {
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
				c.created_at, 
				c.updated_at,
				pi.id AS image_id, 
				pi.product_id,
				pi.image_url,
				pi.is_logo,
				pi.created_at, 
				pi.updated_at
			FROM products AS p
			INNER JOIN categories AS c ON c.id = p.category_id
			INNER JOIN product_images AS pi ON pi.product_id = p.id
			WHERE p.id = ?
		`
	rows, err := tx.QueryContext(ctx, sql, id)
	if err != nil {
		return Product{}, err
	}
	defer rows.Close()

	var product Product
	var images []ProductImages
	for rows.Next() {
		var image ProductImages
		if err := rows.Scan(
			&product.Id, &product.CategoryId, &product.Name, &product.Slug, &product.RealPrice,
			&product.Discount, &product.Stock, &product.Description, &product.CreatedAt, &product.UpdatedAt,
			&product.Category.Id, &product.Category.Name, &product.Category.CreatedAt,
			&product.Category.UpdatedAt, &image.Id, &image.ProductId,
			&image.ImageUrl, &image.IsLogo, &image.CreatedAt, &image.UpdatedAt); err != nil {
			return Product{}, err
		}

		images = append(images, image)
	}
	product.ProductImages = images

	return product, nil
}

func (p *ProductRepositoryImpl) FindLatestImage(ctx context.Context, tx *sql.Tx, productId int) (ProductImages, error) {
	sql := `
			SELECT 
				ps.id, 
				ps.product_id,
				ps.image_url,
				ps.is_logo,
				ps.created_at, 
				ps.updated_at
			FROM product_images AS ps
			INNER JOIN products AS pi ON pi.id = ps.product_id
			WHERE pi.id = ?
            ORDER BY ps.updated_at DESC
            LIMIT 1
		`
	row, err := tx.QueryContext(ctx, sql, productId)
	if err != nil {
		return ProductImages{}, err
	}
	defer row.Close()

	var productImage ProductImages
	if row.Next() {
		if err := row.Scan(
			&productImage.Id, &productImage.ProductId,
			&productImage.ImageUrl, &productImage.IsLogo, &productImage.CreatedAt, &productImage.UpdatedAt); err != nil {
			return ProductImages{}, err
		}
	}

	return productImage, nil
}

func (p *ProductRepositoryImpl) FindBySlug(ctx context.Context, tx *sql.Tx, slug string) (Product, error) {
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
				c.created_at, 
				c.updated_at,
				pi.id AS image_id, 
				pi.product_id,
				pi.image_url,
				pi.is_logo,
				pi.created_at, 
				pi.updated_at
			FROM products AS p
			INNER JOIN categories AS c ON c.id = p.category_id
			INNER JOIN product_images AS pi ON pi.product_id = p.id
			WHERE p.slug = ?
		`
	rows, err := tx.QueryContext(ctx, sql, slug)
	if err != nil {
		return Product{}, err
	}
	defer rows.Close()

	var product Product
	var images []ProductImages
	for rows.Next() {
		var image ProductImages
		if err := rows.Scan(
			&product.Id, &product.CategoryId, &product.Name, &product.Slug, &product.RealPrice,
			&product.Discount, &product.Stock, &product.Description, &product.CreatedAt, &product.UpdatedAt,
			&product.Category.Id, &product.Category.Name, &product.Category.CreatedAt,
			&product.Category.UpdatedAt, &image.Id, &image.ProductId,
			&image.ImageUrl, &image.IsLogo, &image.CreatedAt, &image.UpdatedAt); err != nil {
			return Product{}, err
		}

		images = append(images, image)
	}
	product.ProductImages = images

	return product, nil
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
	_, err := tx.ExecContext(ctx, sql, productImages.ProductId, productImages.ImageUrl, productImages.IsLogo)
	if err != nil {
		return err
	}

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
	SET is_logo = TRUE WHERE id = ?
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

func (p *ProductRepositoryImpl) FindImageById(ctx context.Context, tx *sql.Tx, id int) (ProductImages, error) {
	sql := "select id, product_id, image_url, is_logo, created_at, updated_at from product_images where id = ?"
	row, err := tx.QueryContext(ctx, sql, id)
	if err != nil {
		return ProductImages{}, err
	}
	defer row.Close()

	var productImages ProductImages
	if row.Next() {
		if err := row.Scan(&productImages.Id, &productImages.ProductId, &productImages.ImageUrl, &productImages.IsLogo, &productImages.CreatedAt, &productImages.UpdatedAt); err != nil {
			return ProductImages{}, err
		}

		return productImages, nil
	}

	return ProductImages{}, custom.ErrNotFound
}
