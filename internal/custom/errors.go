package custom

import (
	"catalog-digital-product/internal/helper"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAlreadyExists = errors.New("resource already exists")
	ErrNotFound      = errors.New("resource not found")
	ErrInternal      = errors.New("internal server error")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("you are not authorized to access this resource")
)

func HandleErrorResponde(c *gin.Context, err error) {
	webResponse := helper.WebResponse{
		Data: nil,
	}

	switch err {
	case ErrAlreadyExists:
		webResponse.Code = http.StatusConflict
		webResponse.Status = "error"
		webResponse.Message = err.Error()
	case ErrNotFound:
		webResponse.Code = http.StatusNotFound
		webResponse.Status = "error"
		webResponse.Message = err.Error()
	case ErrInternal:
		webResponse.Code = http.StatusInternalServerError
		webResponse.Status = "error"
		webResponse.Message = err.Error()
	case bcrypt.ErrMismatchedHashAndPassword:
		webResponse.Code = http.StatusUnauthorized
		webResponse.Status = "error"
		webResponse.Message = "email or password is incorrect"
	case ErrUnauthorized:
		webResponse.Code = http.StatusUnauthorized
		webResponse.Status = "error"
		webResponse.Message = "unauthorized"
	default:
		webResponse.Code = http.StatusInternalServerError
		webResponse.Status = "error"
		webResponse.Message = err.Error()
	}

	helper.APIResponse(c, webResponse)
}
