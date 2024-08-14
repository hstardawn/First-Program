package router

import (
	"FORUM/app/controllers/postController"
	"FORUM/app/controllers/userControllers"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	const pre = "/api"

	api := r.Group(pre)
	{
		api.POST("/user/login", userControllers.Login)
		api.POST("/user/reg", userControllers.Register)

		post := api.Group("/student/post")
		{
			post.POST("", postController.CreatPost)
			post.GET("", postController.GetPost)
			post.DELETE("", postController.DeletePost)
			post.PUT("", postController.UpdatePost)
		}

		reportpost := api.Group("/student/report-post")
		{
			reportpost.POST("", postController.CreateReport)
			reportpost.GET("", postController.CheckReport)
		}

		report := api.Group("/admin/report")
		{
			report.GET("", postController.GetCheckReport)
			report.POST("", postController.TrailPost)
		}
	}
}
