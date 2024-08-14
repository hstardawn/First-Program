package userControllers

import (
	"FORUM/app/models"
	"FORUM/app/services/userService"
	"FORUM/app/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录
func Login(c *gin.Context) {
	var data LoginData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return

	}

	err = userService.CheckUsrExistUsername(data.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, 200506, "用户不存在")
		} else {
			utils.JsonInternalServerErrorResponse(c)
		}
		return
	}

	var user *models.User
	user, err = userService.GetUserByUsername(data.Username)
	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	flag := userService.ComparePwd(data.Password, user.Password)
	if !flag {
		utils.JsonErrorResponse(c, 200507, "密码错误")
		return
	}

	utils.JsonSuccessResponse(c, gin.H{
		"user_id":   user.User_id,
		"user_type": user.User_type,
	})
}
