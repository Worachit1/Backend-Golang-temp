package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Officer(router *gin.RouterGroup) {
	ctl := controller.New()
	md := middleware.AuthMiddleware()
	officer := router.Group("")
	{
		officer.POST("/register", ctl.OfficerCtl.Create)
		officer.PATCH("/:id", md, ctl.OfficerCtl.Update)
		officer.GET("/list", md, ctl.OfficerCtl.List)
		officer.GET("/:id", md, ctl.OfficerCtl.Get)
		officer.DELETE("/:id", md, ctl.OfficerCtl.Delete)
	}
}