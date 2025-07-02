package controller

import (
	"app/app/controller/emergency"
	"app/app/controller/login"
	"app/app/controller/logout"
	"app/app/controller/officers"
	"app/app/controller/user"

	"app/config"
)

type Controller struct {
	LoginCtl     *login.Controller     // Assuming LoginController is in the student package
	LogoutCtl    *logout.Controller    // Uncomment if you have a logout controller
	UserCtl      *user.Controller      // Uncomment if you have a user controller
	OfficerCtl   *officers.Controller  // Assuming OfficerController is in the user package
	EmergencyCtl *emergency.Controller // Emergency controller
	// Other controllers...
}

func New() *Controller {
	db := config.GetDB()
	return &Controller{

		LoginCtl:     login.NewController(db),     // Assuming LoginController is in the student package
		LogoutCtl:    logout.NewController(db),    // Uncomment if you have a logout controller
		UserCtl:      user.NewController(db),      // Uncomment if you have a user controller
		OfficerCtl:   officers.NewController(db),  // Assuming OfficerController is in the user package
		EmergencyCtl: emergency.NewController(db), // Emergency controller
		// Other controllers...
	}
}
