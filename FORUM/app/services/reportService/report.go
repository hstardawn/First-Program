package reportService

import (
	"FORUM/app/models"
	"FORUM/config/database"
)

func CheckReportExist(postid uint, userid uint) error {
	result := database.DB.Where("post_id= ?" , postid).First(&models.Report{})
	return result.Error
}

func CreateReportPost(report models.Report) error {
	result := database.DB.Create(&report)
	return result.Error
}

func GetReport(userid uint) ([]models.Report, error) {
	var report_list []models.Report
	result := database.DB.Where("user_id = ?", userid).Select("post_id,reason ,status").Find(&report_list)
	return report_list, result.Error
}

func GetCheckReport() ([]models.Report, error) {
	var reportpost []models.Report
	result := database.DB.Where("status =?", 0).Find(&reportpost)
	return reportpost , result.Error
}

func UpdateReport(report models.Report) error {
	result := database.DB.Where("post_id=?", report.PostId).Save(&report)
	return result.Error
}

func GetReportData(postid uint) (models.Report, error) {
	var report models.Report
	result := database.DB.Where("post_id = ?", postid).First(&report)
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
