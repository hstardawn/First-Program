package midwares

import (
	"FORUM/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleNotFound(c *gin.Context) {
	utils.JsonResponse(c, 404, 200404, nil, http.StatusText(http.StatusNotFound))
}
