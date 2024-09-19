package userControllers

import (
	"FORUM/app/midwares"
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
	UserType int    `json:"user_type" binding:"required"`
}

func Register(c *gin.Context) {
	var data RegisterData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}

	user, err := userService.GetUserByUsername(data.Username)
	if user != nil && len(user.Password)!= 0{
		utils.JsonErrorResponse(c, 200505, "用户名已存在")
		print(user)
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.JsonErrorResponse(c, 200506 ,"查找失败")
		return
	}

	// 检验密码
	flag := userService.IsUsernameAllDigits(data.Password)
	if !flag {
		utils.JsonErrorResponse(c, 200502, "用户名必须为纯数字")
		return
	}

	flag2 := userService.CheckPasswordLength(data.Password)
	if !flag2 {
		utils.JsonErrorResponse(c, 200503, "密码长度必须在8-16位")
		return
	}

	flag3 := userService.CheckUserType(data.UserType)
	if !flag3 {
		utils.JsonErrorResponse(c, 200504, "用户类型错误")
		return
	}
	// 加密
	hashpassword , err:= midwares.HashPassword(data.Password)
	if err != nil{
		utils.JsonErrorResponse(c, 200506, "加密失败")
		return
	}
	err = userService.Register(models.User{
		Username:  data.Username,
		Name:      data.Name,
		Password:  hashpassword,
		UserType: data.UserType,
	})

	if err != nil {
		utils.JsonErrorResponse(c, 200506, "注册失败")
		return
	}

	utils.JsonSuccessResponse(c, nil)

}
