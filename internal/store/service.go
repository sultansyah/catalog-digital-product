package store

import (
	"catalog-digital-product/internal/custom"
	"catalog-digital-product/internal/helper"
	"context"
	"database/sql"
)

type StoreService interface {
	Update(ctx context.Context, input UpdateStoreInput) (Store, error)
	GetStore(ctx context.Context) (Store, error)
}

type StoreServiceImpl struct {
	StoreRepository StoreRepository
	DB              *sql.DB
}

func NewStoreService(storeRepository StoreRepository, DB *sql.DB) StoreService {
	return &StoreServiceImpl{StoreRepository: storeRepository}
}

func (s *StoreServiceImpl) GetStore(ctx context.Context) (Store, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return Store{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	store, err := s.StoreRepository.FindById(ctx, tx, 1)
	if err != nil {
		return Store{}, err
	}
	if store.Id <= 0 {
		return Store{}, custom.ErrNotFound
	}

	return store, nil
}

func (s *StoreServiceImpl) Update(ctx context.Context, input UpdateStoreInput) (Store, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return Store{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	store, err := s.StoreRepository.FindById(ctx, tx, 1)
	if err != nil {
		return Store{}, err
	}
	if store.Id <= 0 {
		return Store{}, custom.ErrNotFound
	}

	store, err = s.StoreRepository.Update(ctx, tx, store)
	if err != nil {
		return Store{}, err
	}

	return store, nil
}
