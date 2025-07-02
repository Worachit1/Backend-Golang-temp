package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Emergency(router *gin.RouterGroup) {
	ctl := controller.New()
	// md := middleware.AuthMiddleware()
	emergency := router.Group("")
	{
		// User routes - for reporting and managing own emergencies
		emergency.POST("/create", ctl.EmergencyCtl.Create)                // Create emergency report
		emergency.PATCH("/officer/:id", ctl.EmergencyCtl.UpdateByOfficer) // Officer update emergency
		emergency.GET("/list", ctl.EmergencyCtl.List)
		emergency.GET("/users/:id", ctl.EmergencyCtl.GetByUserIDEmergency)                     // List all emergencies for debugging
	}
}
