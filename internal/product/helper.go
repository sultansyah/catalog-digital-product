package product

import (
	"catalog-digital-product/internal/helper"
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"os"
	"time"
)

func insertProductImages(p ProductServiceImpl, folder string, fileName string, file multipart.File, errChan chan error, productImage *ProductImages, ctx context.Context, tx *sql.Tx) {
	cwd, err := os.Getwd()
	if err != nil {
		errChan <- err
	}

	productImageFileName := fmt.Sprintf("%d-product-%s", time.Now().Unix(), fileName)
	productImagePath := fmt.Sprintf("%s/public/images/%s/%s", cwd, folder, productImageFileName)

	productImage.ImageUrl = productImageFileName
	err = p.ProductRepository.InsertImage(ctx, tx, *productImage)
	if err != nil {
		errChan <- err
		return
	}

	if err := helper.SaveUploadedFile(file, productImagePath); err != nil {
		defer func() {
			if err := os.Remove(productImagePath); err != nil {
				errChan <- err
			}
		}()
		errChan <- err
		return
	}

	errChan <- nil
}

func deleteImage(folder string, imageName string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	productImagePath := fmt.Sprintf("%s/public/images/%s/%s", cwd, folder, imageName)
	err = os.Remove(productImagePath)
	if err != nil {
		return err
	}
	return nil
}
