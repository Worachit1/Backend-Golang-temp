package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Activity(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	// md := middleware.AuthMiddleware()
	// log := middleware.NewLogResponse()
	activity := router.Group("")
	{
		activity.POST("/create", ctl.ActivityCtl.Create)
		activity.PATCH("/:id", ctl.ActivityCtl.Update)
		activity.GET("/list", ctl.ActivityCtl.List)
		activity.GET("/:id", ctl.ActivityCtl.Get)
		activity.DELETE("/:id", ctl.ActivityCtl.Delete)
	}
}
