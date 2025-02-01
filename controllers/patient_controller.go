package controllers

import (
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"HIS-api/models"
	"HIS-api/config"
)

func SearchPatient(c *gin.Context) {
	// ดึงข้อมูล Staff จาก Context
	staff, exists := c.Get("staff")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, ok := staff.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	hospital, ok := claims["hospital"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token data"})
		return
	}

	log.Println("Searching for patients in hospital:", hospital)

	// กำหนดเงื่อนไขการค้นหา
	conditions := []string{}
	args := []interface{}{}

	if nationalID := c.Query("national_id"); nationalID != "" {
		conditions = append(conditions, "national_id = ?")
		args = append(args, nationalID)
	}
	if passportID := c.Query("passport_id"); passportID != "" {
		conditions = append(conditions, "passport_id = ?")
		args = append(args, passportID)
	}
	if firstName := c.Query("first_name"); firstName != "" {
		conditions = append(conditions, "(first_name_th ILIKE ? OR first_name_en ILIKE ?)")
		args = append(args, "%"+firstName+"%", "%"+firstName+"%")
	}
	if middleName := c.Query("middle_name"); middleName != "" {
		conditions = append(conditions, "(middle_name_th ILIKE ? OR middle_name_en ILIKE ?)")
		args = append(args, "%"+middleName+"%", "%"+middleName+"%")
	}
	if lastName := c.Query("last_name"); lastName != "" {
		conditions = append(conditions, "(last_name_th ILIKE ? OR last_name_en ILIKE ?)")
		args = append(args, "%"+lastName+"%", "%"+lastName+"%")
	}

	// ตรวจสอบและแปลง date_of_birth
	if dobStr := c.Query("date_of_birth"); dobStr != "" {
		dob, err := time.Parse("2006-01-02", dobStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_of_birth format. Use YYYY-MM-DD"})
			return
		}
		conditions = append(conditions, "date_of_birth = ?")
		args = append(args, dob)
	}

	if phone := c.Query("phone_number"); phone != "" {
		conditions = append(conditions, "phone_number = ?")
		args = append(args, phone)
	}
	if email := c.Query("email"); email != "" {
		conditions = append(conditions, "email = ?")
		args = append(args, email)
	}

	// ป้องกันการ Query ข้อมูลทั้งหมดถ้าไม่มีเงื่อนไขใดเลย
	if len(conditions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one search criteria is required"})
		return
	}

	// ใช้ strings.Join() เพื่อสร้าง Query String ที่ปลอดภัย
	queryStr := strings.Join(conditions, " AND ")

	// Query ข้อมูลจาก DB
	var patients []models.Patient
	if err := config.DB.Where(queryStr, args...).Find(&patients).Error; err != nil {
		log.Println("Database Query Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching patients"})
		return
	}

	// ตรวจสอบว่ามีผู้ป่วยที่พบหรือไม่
	if len(patients) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Patients not found",
			"patients": []models.Patient{},
		})
		return
	}

	// ตรวจสอบว่า Staff กำลังค้นหาผู้ป่วยจากโรงพยาบาลตัวเอง
	for _, patient := range patients {
		if patient.Hospital != hospital {
			log.Println("Unauthorized access attempt! Staff from", hospital, "tried to access patient in", patient.Hospital)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: You can only search for patients in your own hospital"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"patients": patients})
}
