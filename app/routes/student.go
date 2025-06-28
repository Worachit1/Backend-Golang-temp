package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Student(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	// md := middleware.AuthMiddleware()
	// log := middleware.NewLogResponse()
	student := router.Group("")
	{
		student.POST("/create", ctl.StudentCtl.Create)
		student.GET("/list", ctl.StudentCtl.List)
	}
}
