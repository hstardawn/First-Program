package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JsonResponse(c *gin.Context, httpStatusCode int, code int, data interface{}, msg string) {
	c.JSON(httpStatusCode, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

func JsonSuccessResponse(c *gin.Context, data interface{}) {
	JsonResponse(c, http.StatusOK, 200, data, "sucesses")
}

func JsonErrorResponse(c *gin.Context, code int, msg string) {
	JsonResponse(c, http.StatusInternalServerError, code, nil, msg)
}

func JsonInternalServerErrorResponse(c *gin.Context) {
	JsonResponse(c, http.StatusInternalServerError, 200500, nil, "Internal server error")
}
