package postController

import (
	"FORUM/app/models"
	"FORUM/app/services/postService"
	"FORUM/app/services/reportService"
	"FORUM/app/services/sessionService"
	"FORUM/app/services/userService"
	"FORUM/app/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 举报帖子
type ReportData struct {
	UserId uint   `json:"user_id" binding:"required"`
	PostId uint   `json:"post_id" binding:"required"`
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
	_, err = sessionService.GetUserSession(c)
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "用户不存在")
		return
	}

	// 获取帖子，并判断存在情况
	_,err = postService.GetPost(data.PostId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, 200506, "帖子不存在")
			return
		} else {
			utils.JsonErrorResponse(c, 200506, "获取帖子失败")
			return
		}
	}

	// 判断重复举报
	err = reportService.CheckReportExist(data.PostId, data.UserId)
	if err == nil {
		utils.JsonErrorResponse(c, 200506, "请勿重复举报")
		return
	}

	// 创建举报
	err = reportService.CreateReportPost(models.Report{
		UserId: data.UserId,
		PostId: data.PostId,
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
	UserId uint `form:"user_id" binding:"required"`
}

func CheckReport(c *gin.Context) {
	var data GetReport
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}

	// 判断用户存在
	_,err = userService.GetUserByUserid(data.UserId)
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "用户不存在")
		return
	}

	// 获取举报列表反馈
	var report_list []models.Report
	report_list, err = reportService.GetReport(data.UserId)
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "获取失败")
		return
	}
	if len(report_list)== 0{
		utils.JsonErrorResponse(c, 200506, "举报列表为空")
		return
	}

	// 变更格式，去userid
	newReportList := make([]models.ReportList, len(report_list))
	for i, report := range report_list {
		post, err := postService.GetPost(report.PostId)
		if err != nil {
			utils.JsonErrorResponse(c, 200506, "获取帖子失败")
			return
		}
		newReportList[i] = models.ReportList{
			PostId: report.PostId,
			Content: post.Content,
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
	UserId uint `form:"user_id" binding:"required"`
}

func GetCheckReport(c *gin.Context) {
	var data GetReport
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
	}

	// 判断用户存在
	user,err := userService.GetUserByUserid(data.UserId)
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "用户不存在")
		return
	}

	// 判断权限
	if user.UserType == 1 {
		utils.JsonErrorResponse(c, 200506, "权限不足")
		return
	}

	// 获取列表
	reportpost, err := reportService.GetCheckReport()
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "获取失败")
		return
	}
	if len(reportpost)== 0{
		utils.JsonErrorResponse(c, 200506, "举报列表为空")
		return
	}
	report_list := make([]models.Check, len(reportpost))
	for i, report := range reportpost {
		user ,err:=userService.GetUserByUserid(report.UserId)
		if err != nil {
			utils.JsonErrorResponse(c, 200506, "获取用户失败")
			return
		}

		post, err := postService.GetPost(report.PostId)
		if err != nil {
			utils.JsonErrorResponse(c, 200506, "获取帖子失败")
			return
		}
		report_list[i].Content = post.Content
		report_list[i].Reason = report.Reason
		report_list[i].Username = user.Username
		report_list[i].PostId = report.PostId
	}
	utils.JsonSuccessResponse(c, report_list)
}

// 审核被举报的帖子
type CheckPost struct {
	UserId  uint `json:"user_id" binding:"required"`
	PostId  uint `json:"post_id" binding:"required"`
	Approval int  `json:"approval" binding:"required"`
}

func TrailPost(c *gin.Context) {
	var data CheckPost
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, 200501, "参数错误")
		return
	}
	// 判断用户存在

	user,err := userService.GetUserByUserid(data.UserId)
	if err != nil {
		utils.JsonErrorResponse(c, 200506, "用户不存在")
		return
	}

	// 判断权限
	if user.UserType == 1{
		utils.JsonErrorResponse(c, 200506, "权限不足")
		return
	}
	// 判断帖子是否存在，存在便获取
	report, err := reportService.GetReportData(data.PostId)
	if err != nil {
		utils.JsonErrorResponse(c, 200504, "帖子不存在")
		return
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
			UserId: report.UserId,
			PostId: report.PostId,
			Reason:  report.Reason,
			Status:  data.Approval,
		})
		if err != nil {
			utils.JsonErrorResponse(c, 200506, "更新失败")
			return
		}
		err = postService.DeletePost(data.PostId)
		if err != nil {
			utils.JsonErrorResponse(c, 200506, "删除失败")
			return
		}
		// 不同意，更新状态
	} else {
		err = reportService.UpdateReport(models.Report{
			UserId: report.UserId,
			PostId: report.PostId,
			Reason:  report.Reason,
			Status:  data.Approval,
		})
		if err != nil {
			utils.JsonErrorResponse(c, 200506, "更新失败")
			return
		}

	}

	utils.JsonSuccessResponse(c, nil)
}
