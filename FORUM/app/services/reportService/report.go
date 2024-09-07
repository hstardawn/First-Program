package reportService

import (
	"FORUM/app/models"
	"FORUM/app/services/userService"
	"FORUM/config/database"

	"gorm.io/gorm"
)

func CheckReportExist(post_id uint) error {
	result := database.DB.Where("post_id= ?" , post_id).First(&models.Report{})
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

func GetCheckReport(user_id uint) ([]models.Check, error) {
	var reportpost []models.Report
	result := database.DB.Where("status =?", 0).Find(&reportpost)
	if result.RowsAffected == 0 {
		result.Error =gorm.ErrRecordNotFound
	}
	if result.Error != nil{
		return nil, result.Error
	}
	report_list := make([]models.Check, len(reportpost))
	for i, post := range reportpost {
		user ,err:=userService.GetUserByUserid(post.User_id)
		if err != nil {
			return report_list, err
		}

		report_list[i].Content = post.Content
		report_list[i].Reason = post.Reason
		report_list[i].Username = user.Username
		report_list[i].Post_id = post.Post_id
	}

	return report_list, result.Error
}

func UpdateReport(report models.Report) error {
	result := database.DB.Where("post_id=?", report.Post_id).Save(&report)
	return result.Error
}

func GetReportData(post_id uint) (models.Report, error) {
	var report models.Report
	result := database.DB.Where("post_id = ?", post_id).First(&report)
	return report, result.Error
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
