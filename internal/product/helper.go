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

func insertProductImages(p ProductServiceImpl, fileName string, file multipart.File, errChan chan error, productImage *ProductImages, ctx context.Context, tx *sql.Tx) {
	cwd, err := os.Getwd()
	if err != nil {
		errChan <- err
	}

	productImageFileName := fmt.Sprintf("%d-store-%s", time.Now().Unix(), fileName)
	productImagePath := fmt.Sprintf("%s/public/images/store/%s", cwd, productImageFileName)

	productImage.ImageUrl = productImageFileName
	err = p.ProductRepository.InsertImage(ctx, tx, *productImage)
	if err != nil {
		errChan <- err
		return
	}

	if err := helper.SaveUploadedFile(file, productImagePath); err != nil {
		defer os.Remove(productImagePath)
		errChan <- err
		return
	}

	errChan <- nil
}
