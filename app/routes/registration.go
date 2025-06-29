package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Registration(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	// md := middleware.AuthMiddleware()
	// log := middleware.NewLogResponse()
	registration := router.Group("")
	{
		registration.POST("/create", ctl.RegistrationCtl.Create)
		registration.PUT("/:id", ctl.RegistrationCtl.Update)
		registration.GET("/list", ctl.RegistrationCtl.List)
		registration.GET("/:id", ctl.RegistrationCtl.Get)
		registration.DELETE("/:id", ctl.RegistrationCtl.Delete)
	}
}
