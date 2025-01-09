package category

import (
	"catalog-digital-product/internal/custom"
	"context"
	"database/sql"
)

type CategoryRepository interface {
	Insert(ctx context.Context, tx *sql.Tx, category Category) (Category, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]Category, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (Category, error)
	Update(ctx context.Context, tx *sql.Tx, category Category) (Category, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) error
}

type CategoryRepositoryImpl struct {
}

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{}
}

func (c *CategoryRepositoryImpl) Insert(ctx context.Context, tx *sql.Tx, category Category) (Category, error) {
	sql := "insert into categories(name) VALUES (?)"
	result, err := tx.ExecContext(ctx, sql, category.Name)
	if err != nil {
		return Category{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Category{}, err
	}

	category.Id = int(id)

	return category, nil
}

func (c *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]Category, error) {
	sql := "select id, name, created_at, updated_at from categories"
	rows, err := tx.QueryContext(ctx, sql)
	if err != nil {
		return []Category{}, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category

		if err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return []Category{}, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (c *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (Category, error) {
	sql := "select id, name, created_at, updated_at from categories where id = ?"
	row, err := tx.QueryContext(ctx, sql, id)
	if err != nil {
		return Category{}, err
	}
	defer row.Close()

	var category Category
	if row.Next() {
		if err := row.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return Category{}, err
		}

		return category, nil
	}

	return category, custom.ErrNotFound
}

func (c *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	sql := "delete from categories where id = ?"
	_, err := tx.ExecContext(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category Category) (Category, error) {
	sql := "update categories set name = ? where id = ?"
	_, err := tx.ExecContext(ctx, sql, category.Name, category.Id)
	if err != nil {
		return Category{}, err
	}

	return category, nil
}
