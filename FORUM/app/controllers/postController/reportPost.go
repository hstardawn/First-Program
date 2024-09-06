package postController

import (
	"FORUM/app/models"
	"FORUM/app/services/postService"
	"FORUM/app/services/reportService"
	"FORUM/app/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 举报帖子
type ReportData struct {
	User_id uint   `json:"user_id" binding:"required"`
	Post_id uint   `json:"post_id" binding:"required"`
	Reason  string `json:"reason" binding:"required"`
}

func CreateReport(c *gin.Context) {
	var data ReportData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}

	// 判断用户存在
	err = postService.CheckUserByUserid(data.User_id)
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "用户不存在")
		return
	}

	// 获取帖子，并判断存在情况
	Content, err := reportService.GetContentFromDB(data.Post_id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, 200506, "帖子不存在")
			return
		} else {
			utils.JsonInternalServerErrorResponse(c)
			return
		}
	}

	// 判断重复举报
	err = reportService.CheckReportExist(data.User_id, data.Post_id)
	if err == nil {
		utils.JsonErrorResponse(c, 200506, "请勿重复举报")
		return
	}

	// 创建举报
	err = reportService.CreateReportPost(models.Report{
		User_id: data.User_id,
		Post_id: data.Post_id,
		Content: Content,
		Reason:  data.Reason,
	})
	if err != nil {
		utils.JsonInternalServerErrorResponse(c)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

// 查看举报审批
type GetReport struct {
	User_id uint `form:"user_id" binding:"required"`
}

func CheckReport(c *gin.Context) {
	var data GetReport
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}

	// 判断用户存在
	err = postService.CheckUserByUserid(data.User_id)
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "用户不存在")
		return
	}

	// 获取举报列表反馈
	var report_list []models.Report
	report_list, err = reportService.GetReport(data.User_id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, 200506, "举报列表为空")
			return
		} else {
			utils.JsonInternalServerErrorResponse(c)
			return
		}
	}

	// 变更格式，去userid
	newReportList := make([]models.ReportList, len(report_list))
	for i, report := range report_list {
		newReportList[i] = models.ReportList{
			Post_id: report.Post_id,
			Content: report.Content,
			Reason:  report.Reason,
			Status:  report.Status,
		}
	}

	utils.JsonSuccessResponse(c, gin.H{
		"report_list": newReportList,
	})
}

// 获取所有未审批的被举报帖子
type GetCheckData struct {
	User_id uint `form:"user_id" binding:"required"`
}

func GetCheckReport(c *gin.Context) {
	var data GetReport
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
	}

	// 判断用户存在
	err = postService.CheckUserByUserid(data.User_id)
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "用户不存在")
		return
	}

	// 判断权限
	right := reportService.CheckRight(data.User_id)
	if right == 1 {
		utils.JsonErrorResponse(c, 200506, "权限不足")
		return
	}

	// 获取列表
	var report_list []models.Check
	report_list, err = reportService.GetCheckReport(data.User_id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, 200506, "举报列表为空")
			return
		} else {
			utils.JsonInternalServerErrorResponse(c)
			return
		}
	}

	utils.JsonSuccessResponse(c, report_list)
}

// 审核被举报的帖子
type CheckPost struct {
	User_id  uint `json:"user_id" binding:"required"`
	Post_id  uint `json:"post_id" binding:"required"`
	Approval int  `json:"approval" binding:"required"`
}

func TrailPost(c *gin.Context) {
	var data CheckPost
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}
	fmt.Println(data.Approval)
	// 判断用户存在
	err = postService.CheckUserByUserid(data.User_id)
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "用户不存在")
		return
	}

	// 判断权限
	right := reportService.CheckRight(data.User_id)
	if right == 1 {
		utils.JsonErrorResponse(c, 200506, "权限不足")
		return
	}

	// 判断帖子是否存在，存在便获取
	report, err := reportService.GetReportData(data.Post_id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, 200506, "帖子不存在")
			return
		} else {
			utils.JsonInternalServerErrorResponse(c)
			return
		}
	}

	// 判断帖子状态
	if report.Status != 0 {
		utils.JsonErrorResponse(c, 200506, "帖子已被审批")
		return
	}

	// 检验aproval数值
	flag := reportService.CheckApproval(data.Approval)
	if !flag {
		utils.JsonErrorResponse(c, 200504, "审核类型错误")
		return
	}

	// 同意，更新数据，进行删除
	if data.Approval == 1 {
		err = reportService.UpdateReport(models.Report{
			User_id: report.User_id,
			Post_id: report.Post_id,
			Content: report.Content,
			Reason:  report.Reason,
			Status:  data.Approval,
		})
		if err != nil {
			utils.JsonInternalServerErrorResponse(c)
			return
		}
		err = postService.DeletePost(data.Post_id)
		if err != nil {
			utils.JsonInternalServerErrorResponse(c)
			return
		}
		// 不同意，更新状态
	} else {
		err = reportService.UpdateReport(models.Report{
			User_id: report.User_id,
			Post_id: report.Post_id,
			Content: report.Content,
			Reason:  report.Reason,
			Status:  data.Approval,
		})
		if err != nil {
			utils.JsonInternalServerErrorResponse(c)
			return
		}

	}

	utils.JsonSuccessResponse(c, nil)
}
