package store

import (
	"catalog-digital-product/internal/custom"
	"catalog-digital-product/internal/helper"
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"os"
	"time"
)

type StoreService interface {
	Update(ctx context.Context, input UpdateStoreInput, imageFile multipart.File, imageFileName string) (Store, error)
	GetStore(ctx context.Context) (Store, error)
}

type StoreServiceImpl struct {
	StoreRepository StoreRepository
	DB              *sql.DB
}

func NewStoreService(storeRepository StoreRepository, DB *sql.DB) StoreService {
	return &StoreServiceImpl{StoreRepository: storeRepository, DB: DB}
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

func (s *StoreServiceImpl) Update(ctx context.Context, input UpdateStoreInput, imageFile multipart.File, imageFileName string) (Store, error) {
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

	cwd, err := os.Getwd()
	if err != nil {
		return Store{}, err
	}

	imageFileName = fmt.Sprintf("%d-store-%s", time.Now().Unix(), imageFileName)
	imagePath := fmt.Sprintf("%s/public/images/store/%s", cwd, imageFileName)

	store.ImageURL = imageFileName
	store, err = s.StoreRepository.Update(ctx, tx, store)
	if err != nil {
		return Store{}, err
	}

	if err := helper.SaveUploadedFile(imageFile, imagePath); err != nil {
		defer os.Remove(imagePath)
		return Store{}, err
	}

	return store, nil
}
