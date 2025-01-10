package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) []string {
	var errors []string

	for _, v := range err.(validator.ValidationErrors) {
		errors = append(errors, v.Error())
	}

	return errors
}

func BindAndValidate(c *gin.Context, input any, bind string) bool {
	webResponse := WebResponse{
		Code:    http.StatusUnprocessableEntity,
		Status:  "error",
		Message: "invalid request payload",
	}

	switch bind {
	case "json":
		if err := c.ShouldBindJSON(input); err != nil {
			webResponse.Data = err.Error()
			APIResponse(c, webResponse)
			return false
		}
		return true
	case "uri":
		if err := c.ShouldBindUri(input); err != nil {
			errors := FormatValidationErrors(err)
			webResponse.Data = errors
			APIResponse(c, webResponse)
			return false
		}
		return true
	case "form":
		if err := c.ShouldBind(input); err != nil {
			webResponse.Data = err.Error()
			APIResponse(c, webResponse)
			return false
		}
		return true
	default:
		return false
	}
}
