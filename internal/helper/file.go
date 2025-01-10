package helper

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveUploadedFile(file multipart.File, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, file)

	return err
}
