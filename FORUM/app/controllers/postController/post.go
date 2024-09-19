package postController

import (
	"FORUM/app/models"
	"FORUM/app/services/postService"
	"FORUM/app/services/userService"
	"FORUM/app/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostData struct {
	Content string `json:"content" binding:"required"`
	UserId uint   `json:"user_id" binding:"required"`
}

// 发布帖子
func CreatPost(c *gin.Context) {
	var data PostData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}

	// 判断用户是否存在
	_, err = userService.GetUserByUserid(data.UserId)
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "用户不存在")
		
		return
	}

	now := time.Now()
	err = postService.CreatPost(models.Post{
		Content: data.Content,
		UserId: data.UserId,
		Time:    now,
	})

	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}
	println(err)
	utils.JsonSuccessResponse(c, nil)
}

// 修改帖子
type UpdatePostData struct {
	UserId uint   `json:"user_id" binding:"required"`
	PostId uint   `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func UpdatePost(c *gin.Context) {
	var data UpdatePostData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}
	// 判断用户是否存在
	_, err = userService.GetUserByUserid(data.UserId)
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "用户不存在")
		
		return
	}
	// 判断帖子是否存在
	post , err := postService.GetPost(data.PostId)
	if err != nil{
		utils.JsonErrorResponse(c, 200506, "查找失败")
		return
	}
	
	// 判断用户
	if  post.UserId != data.UserId{
		utils.JsonErrorResponse(c, 200506, "不是帖子主人，无权修改")
		return
	}

	// 修改帖子
	now := time.Now()
	err = postService.UpdatePost(models.Post{
		UserId: data.UserId,
		PostId: data.PostId,
		Content: data.Content,
		Time:    now,
	})
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "修改帖子失败")
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

// 删除帖子
type DeletePostData struct {
	UserId uint `form:"user_id" binding:"required"`
	PostId uint `form:"post_id" binding:"required"`
}

func DeletePost(c *gin.Context) {
	var data DeletePostData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}
	// 判断用户是否存在
	_, err = userService.GetUserByUserid(data.UserId)
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "用户不存在")
		
		return
	}
	
	// 判断帖子是否存在
	post, err := postService.GetPost(data.PostId)
	if err == gorm.ErrRecordNotFound {
		utils.JsonErrorResponse(c, 200506, "帖子不存在")
		return
	}

	// 判断用户
	if  post.UserId != data.UserId {
		utils.JsonErrorResponse(c, 200506, "不是帖子主人，无权修改")
		return
	}

	// 进行删除
	err = postService.DeletePost(data.PostId)
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "删除失败")
		return
	}
	// 返回成功
	utils.JsonSuccessResponse(c, nil)
}

// 获取所有发布的帖子
func GetPost(c *gin.Context) {
	var post_list []models.PostList
	post_list, err := postService.GetPostList()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, 200506, "未发布帖子")
			return
		} else {
			utils.JsonErrorResponse(c, 200506, "获取失败")
			return
		}
	}

	utils.JsonSuccessResponse(c, gin.H{
		"post_list": post_list,
	})

}
