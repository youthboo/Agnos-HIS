package controllers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"HIS-api/models"
	"HIS-api/config"
	"os"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func generateToken(username string, hospital string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"hospital": hospital, 
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func RegisterStaff(c *gin.Context) {
	var staff models.Staff
	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ตรวจสอบค่า `username`, `password`, `hospital` ต้องไม่ว่าง
	if staff.Username == "" || staff.Password == "" || staff.Hospital == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	// เช็คว่า `username` ต้องไม่ซ้ำ
	var existingStaff models.Staff
	if err := config.DB.Where("username = ?", staff.Username).First(&existingStaff).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(staff.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	staff.Password = string(hashedPassword)

	// บันทึกลงฐานข้อมูล
	if err := config.DB.Create(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving staff to database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Staff registered successfully!"})
}

func LoginStaff(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Hospital string `json:"hospital" binding:"required"`
	}

	// ตรวจสอบว่ามีค่าครบ
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields (username, password, hospital) are required"})
		return
	}

	var storedStaff models.Staff

	// ตรวจสอบ username และ hospital พร้อมกัน
	if err := config.DB.Where("username = ? AND hospital = ?", input.Username, input.Hospital).First(&storedStaff).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// ตรวจสอบ password
	if err := bcrypt.CompareHashAndPassword([]byte(storedStaff.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// สร้าง JWT Token
	token, err := generateToken(storedStaff.Username, storedStaff.Hospital)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// ซ่อน password ก่อนส่ง response
	storedStaff.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"staff":   storedStaff,
	})
}
