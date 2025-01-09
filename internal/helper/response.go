package helper

import "github.com/gin-gonic/gin"

type WebResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func APIResponse(c *gin.Context, response WebResponse) {
	c.JSON(response.Code, response)
}
