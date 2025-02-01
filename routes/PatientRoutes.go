package routes

import (
	"github.com/gin-gonic/gin"
	"HIS-api/controllers"
	"HIS-api/middlewares"
)

func PatientRoutes(r *gin.Engine) {
	patient := r.Group("/patient")
	patient.Use(middlewares.AuthMiddleware()) 
	{
		patient.GET("/search", controllers.SearchPatient)
	}
}
