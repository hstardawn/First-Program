package postService

import (
	"FORUM/app/models"
	"FORUM/config/database"

	"gorm.io/gorm"
)

func CreatPost(post models.Post) error {
	result := database.DB.Create(&post)
	return result.Error
}

func CheckUserByUserid(userid uint) error {
	result := database.DB.Where("user_id=?", userid).First(&models.User{})
	return result.Error
}

func UpdatePost(post models.Post) error {
	result := database.DB.Omit("user_id & post_id ").Save(&post)
	return result.Error
}

func GetUseridByPostid(post_id uint) uint {
	var post models.Post
	result := database.DB.Where("post_id=?", post_id).First(&post)
	if result.Error == gorm.ErrRecordNotFound {
		return 0
	}
	return post.User_id
}

func CheckPostExist(post_id uint) error {
	result := database.DB.Where("post_id= ?", post_id).First(&models.Post{})
	return result.Error
}

func DeletePost(post_id uint) error {
	result := database.DB.Where("post_id= ?", post_id).Delete(&models.Post{})
	return result.Error
}

func TransformPostList(postList []models.Post) []models.PostList {
	newPostList := make([]models.PostList, len(postList))
	for i, post := range postList {
		newPostList[i] = models.PostList{
			Content: post.Content,
			User_id: post.User_id,
			Id:      post.Post_id,
			Time:    post.Time,
		}
	}
	return newPostList
}

func GetPostList() ([]models.PostList, error) {
	var post_list []models.Post
	result := database.DB.Find(&post_list)
	if result.Error != nil {
		return nil, result.Error
	}

	newPostList := TransformPostList(post_list)
	return newPostList, nil
}

