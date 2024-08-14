package userControllers

import (
	"FORUM/app/models"
	"FORUM/app/services/userService"
	"FORUM/app/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 注册
type RegisterData struct {
	Username  string `json:"username" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	User_type int    `json:"user_type" binding:"required"`
}

func Register(c *gin.Context) {
	var data RegisterData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}

	err = userService.CheckUsrExistUsername(data.Username)
	if err == nil {
		utils.JsonErrorResponse(c, 200505, "用户名已存在")
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	// 检验密码
	flag := userService.IsUsernameAllDigits(data.Username)
	if !flag {
		utils.JsonErrorResponse(c, 200502, "用户名必须为纯数字")
		return
	}

	flag2 := userService.CheckPasswordLength(data.Password, 8, 16)
	if !flag2 {
		utils.JsonErrorResponse(c, 200503, "密码长度必须在8-16位")
		return
	}

	flag3 := userService.CheckUserType(data.User_type)
	if !flag3 {
		utils.JsonErrorResponse(c, 200504, "用户类型错误")
		return
	}

	err = userService.Register(models.User{
		Username:  data.Username,
		Name:      data.Name,
		Password:  data.Password,
		User_type: data.User_type,
	})

	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	utils.JsonSuccessResponse(c, nil)

}
