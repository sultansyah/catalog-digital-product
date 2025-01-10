package store

import (
	"catalog-digital-product/internal/custom"
	"context"
	"database/sql"
)

type StoreRepository interface {
	Update(ctx context.Context, tx *sql.Tx, store Store) (Store, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (Store, error)
}

type StoreRepositoryImpl struct {
}

func NewStoreRepository() StoreRepository {
	return &StoreRepositoryImpl{}
}

func (s *StoreRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, store Store) (Store, error) {
	sql := "UPDATE store SET name=?,description=?,location=?,latitude=?,longtitude=?,phone_number=?,email=?,image_url=?,whatsapp_link=? WHERE id = ?"
	result, err := tx.ExecContext(ctx, sql, store.Name, store.Description, store.Location, store.Latitude, store.Longitude, store.PhoneNumber, store.Email, store.ImageURL, store.WhatsappLink, store.Id)
	if err != nil {
		return Store{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Store{}, err
	}

	store.Id = int(id)
	return store, nil
}

func (s *StoreRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (Store, error) {
	sql := "SELECT id, name, description, location, latitude, longtitude, phone_number, email, image_url, whatsapp_link, created_at, updated_at FROM store WHERE id = ?"
	row, err := tx.QueryContext(ctx, sql, id)
	if err != nil {
		return Store{}, err
	}
	defer row.Close()

	var store Store
	if row.Next() {
		if err := row.Scan(&store.Id, &store.Name, &store.Description, &store.Location, &store.Latitude, &store.Longitude, &store.PhoneNumber, &store.Email, &store.ImageURL, &store.WhatsappLink, &store.CreatedAt, &store.UpdatedAt); err != nil {
			return Store{}, err
		}

		return store, nil
	}

	return Store{}, custom.ErrNotFound
}
