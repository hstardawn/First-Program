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

func CheckReportExist(user_id uint, post_id uint) error {
	result := database.DB.Where("user_id=? AND post_id= ?", user_id, post_id).First(&models.Report{})
	return result.Error
}

func CreateReportPost(report models.Report) error {
	result := database.DB.Create(&report)
	return result.Error
}

func GetReport(user_id uint) ([]models.Report, error) {
	result := database.DB.Where("user_id = ?", user_id).Select("post_id ,content ,reason ,status").First(&models.Report{})
	if result.Error != nil {
		return nil, result.Error
	}
	var report_list []models.Report
	result = database.DB.Where("user_id = ?", user_id).Select("post_id ,content ,reason ,status").Find(&report_list)
	if result.Error != nil {
		return nil, result.Error
	}

	return report_list, nil
}

func CheckRight(user_id uint) int {
	var user models.User
	database.DB.Where("user_id=?", user_id).First(&user)
	right := user.User_type
	return right
}
func GetCheckReport(user_id uint) ([]models.Check, error) {
	result := database.DB.Where("status =?", 0).First(&models.Report{})
	if result.Error != nil {
		return nil, result.Error
	}
	var reportpost []models.Report
	result = database.DB.Where("status =?", 0).Find(&reportpost)
	if result.Error != nil {
		return nil, result.Error
	}

	report_list := make([]models.Check, len(reportpost))
	for i, post := range reportpost {
		report_list[i].Content = post.Content
		report_list[i].Reason = post.Reason
		report_list[i].Username = GetUsernameByUserid(post.User_id)
		report_list[i].Post_id = post.Post_id
	}

	return report_list, nil
}

func GetContentFromDB(post_id uint) (string, error) {
	var reportData []models.Post
	result := database.DB.Where("post_id =?", post_id).First(&reportData)
	if result.Error != nil {
		return "", result.Error
	}

	if len(reportData) > 0 {
		return reportData[0].Content, nil
	}
	return "", nil
}

func GetUsernameByUserid(user_id uint) string {
	var user models.User
	result := database.DB.Where("user_id=?", user_id).First(&user)
	if result.Error != nil {
		return ""
	}
	return user.Username
}

func UpdateReport(report models.Report) error {
	result := database.DB.Where("post_id=?", report.Post_id).Save(&report)
	return result.Error
}

func GetReportData(post_id uint) (models.Report, error) {
	var report models.Report
	var defaultReport models.Report
	result := database.DB.Where("post_id = ?", post_id).First(&report)
	if result.Error != nil {
		return defaultReport, result.Error
	}
	return report, nil
}

var validUserTypes = []int{1, 2}

func CheckApproval(userType int) bool {
	for _, validType := range validUserTypes {
		if userType == validType {
			return true
		}
	}
	return false
}
