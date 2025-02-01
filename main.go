package main

import (
	"github.com/gin-gonic/gin"
	"HIS-api/config"
	"HIS-api/routes"
	"HIS-api/database"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	config.ConnectDB()  
	database.MigrateDB() 

	routes.StaffRoutes(r)
	routes.PatientRoutes(r)

	r.Run(":8080") 
}
