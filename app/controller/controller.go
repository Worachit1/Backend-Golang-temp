package controller

import (
	"app/app/controller/activity"
	"app/app/controller/login"
	"app/app/controller/logout"
	"app/app/controller/registration"
	"app/app/controller/student"
	"app/config"
)

type Controller struct {
	StudentCtl      *student.Controller
	ActivityCtl     *activity.Controller
	RegistrationCtl *registration.Controller
	LoginCtl        *login.Controller  // Assuming LoginController is in the student package
	LogoutCtl       *logout.Controller // Uncomment if you have a logout controller

	// Other controllers...
}

func New() *Controller {
	db := config.GetDB()
	return &Controller{

		StudentCtl:      student.NewController(db),
		ActivityCtl:     activity.NewController(db),
		RegistrationCtl: registration.NewController(db),
		LoginCtl:        login.NewController(db),  // Assuming LoginController is in the student package
		LogoutCtl:       logout.NewController(db), // Uncomment if you have a logout controller
		// Other controllers...
	}
}
