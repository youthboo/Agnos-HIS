package routes

import (
	"github.com/gin-gonic/gin"
	"HIS-api/controllers"
)

func StaffRoutes(r *gin.Engine) {
	staff := r.Group("/staff")
	{
		staff.POST("/create", controllers.RegisterStaff)
		staff.POST("/login", controllers.LoginStaff)
	}
}
