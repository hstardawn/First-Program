package postService

import (
	"FORUM/app/models"
	"FORUM/config/database"
)

func CreatPost(post models.Post) error {
	result := database.DB.Create(&post)
	return result.Error
}

func UpdatePost(post models.Post) error {
	result := database.DB.Omit("user_id & post_id ").Save(&post)
	return result.Error
}

func GetPost(post_id uint) (models.Post, error) {
	var post models.Post
	result := database.DB.Unscoped().Where("post_id= ?", post_id).First(&post)
	return post ,result.Error
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
			UserId: post.UserId,
			Id:      post.PostId,
			Time:    post.Time,
		}
	}
	return newPostList
}

func GetPostList() ([]models.PostList, error) {
	var post_list []models.Post
	result := database.DB.Unscoped().Find(&post_list)
	if result.Error != nil {
		return nil, result.Error
	}

	newPostList := TransformPostList(post_list)
	return newPostList, nil
}

