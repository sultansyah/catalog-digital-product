package helper

import (
	"catalog-digital-product/internal/custom"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type WebResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func APIResponse(c *gin.Context, response WebResponse) {
	c.JSON(response.Code, response)
}

func HandleErrorResponde(c *gin.Context, err error) {
	webResponse := WebResponse{
		Data: nil,
	}

	switch err {
	case custom.ErrAlreadyExists:
		webResponse.Code = http.StatusConflict
		webResponse.Status = "error"
		webResponse.Message = err.Error()
	case custom.ErrNotFound:
		webResponse.Code = http.StatusNotFound
		webResponse.Status = "error"
		webResponse.Message = err.Error()
	case custom.ErrInternal:
		webResponse.Code = http.StatusInternalServerError
		webResponse.Status = "error"
		webResponse.Message = err.Error()
	case bcrypt.ErrMismatchedHashAndPassword:
		webResponse.Code = http.StatusUnauthorized
		webResponse.Status = "error"
		webResponse.Message = "username or password is incorrect"
	case custom.ErrUnauthorized:
		webResponse.Code = http.StatusUnauthorized
		webResponse.Status = "error"
		webResponse.Message = "unauthorized"
	default:
		webResponse.Code = http.StatusInternalServerError
		webResponse.Status = "error"
		webResponse.Message = err.Error()
	}

	APIResponse(c, webResponse)
}
