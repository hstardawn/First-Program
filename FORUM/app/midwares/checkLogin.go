package midwares

import (
	"FORUM/app/services/sessionService"
	"FORUM/app/utils"

	"github.com/gin-gonic/gin"
)

func CheckLogin(c *gin.Context) {
	isLogin := sessionService.CheckUserSession(c)
	if !isLogin {
		utils.JsonErrorResponse(c, 401, "用户未登录，请先登录")
		c.Abort()
		return
	}
	c.Next()
}
