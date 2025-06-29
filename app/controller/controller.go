package controller

import (
	"app/app/controller/activity"
	"app/app/controller/registration"
	"app/app/controller/student"
	"app/config"
)

type Controller struct {
	StudentCtl      *student.Controller
	ActivityCtl     *activity.Controller
	RegistrationCtl *registration.Controller

	// Other controllers...
}

func New() *Controller {
	db := config.GetDB()
	return &Controller{

		StudentCtl:      student.NewController(db),
		ActivityCtl:     activity.NewController(db),
		RegistrationCtl: registration.NewController(db),

		// Other controllers...
	}
}
