package tests

import (
	"HIS-api/config"
	"HIS-api/controllers"
	"HIS-api/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupStaffTestDB() {
	config.ConnectDB()
	config.DB.Exec("DELETE FROM staffs")

	// เข้ารหัสรหัสผ่านก่อนบันทึก
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	staff := models.Staff{
		Username: "admin",
		Password: string(hashedPassword),
		Hospital: "Hospital",
	}
	config.DB.Create(&staff)
}

func teardownTestDB() {
	config.DB.Exec("DELETE FROM staffs") // ล้างข้อมูลทดสอบออก
}

// ทดสอบสร้างบัญชี Staff สำเร็จ
func TestRegisterStaff_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupStaffTestDB()
	defer teardownTestDB()

	registerPayload := map[string]string{
		"username": "newuser",
		"password": "password123",
		"hospital": "Hospital",
	}
	payloadBytes, _ := json.Marshal(registerPayload)

	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/staff/create", controllers.RegisterStaff)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// ทดสอบสร้างบัญชี Staff ที่ username ซ้ำกัน (ควรได้ 400)
func TestRegisterStaff_DuplicateUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	defer teardownTestDB()

	registerPayload := map[string]string{
		"username": "admin",
		"password": "password123",
		"hospital": "Hospital",
	}
	payloadBytes, _ := json.Marshal(registerPayload)

	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/staff/create", controllers.RegisterStaff)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ทดสอบกรณีไม่ส่ง password (ควรได้ 400)
func TestRegisterStaff_MissingPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	defer teardownTestDB()

	registerPayload := map[string]string{
		"username": "newuser",
		"hospital": "Hospital",
	}
	payloadBytes, _ := json.Marshal(registerPayload)

	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/staff/create", controllers.RegisterStaff)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ทดสอบกรณีไม่ส่ง hospital (ควรได้ 400)
func TestRegisterStaff_MissingHospital(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	defer teardownTestDB()

	registerPayload := map[string]string{
		"username": "newuser",
		"password": "password123",
	}
	payloadBytes, _ := json.Marshal(registerPayload)

	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/staff/create", controllers.RegisterStaff)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ทดสอบกรณีไม่ส่ง username (ควรได้ 400)
func TestRegisterStaff_MissingUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	defer teardownTestDB()

	registerPayload := map[string]string{
		"password": "password123",
		"hospital": "Hospital",
	}
	payloadBytes, _ := json.Marshal(registerPayload)

	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/staff/create", controllers.RegisterStaff)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ทดสอบกรณีส่ง Request Body ไม่ใช่ JSON (ควรได้ 400)
func TestRegisterStaff_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	defer teardownTestDB()

	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer([]byte("Invalid Body")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/staff/create", controllers.RegisterStaff)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ทดสอบเข้าสู่ระบบสำเร็จ
func TestLoginStaff_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	defer teardownTestDB()

	loginPayload := map[string]string{
		"username": "admin",
		"password": "password",
		"hospital": "Hospital",
	}
	payloadBytes, _ := json.Marshal(loginPayload)

	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/staff/login", controllers.LoginStaff)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	_, exists := response["token"]
	assert.True(t, exists, "Token should be returned")
}

// ทดสอบเข้าสู่ระบบด้วยรหัสผ่านผิด (ควรได้ 401)
func TestLoginStaff_InvalidPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	defer teardownTestDB()

	loginPayload := map[string]string{
		"username": "admin",
		"password": "wrongpassword",
		"hospital": "Hospital",
	}
	payloadBytes, _ := json.Marshal(loginPayload)

	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/staff/login", controllers.LoginStaff)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ทดสอบเข้าสู่ระบบโดยไม่มี `hospital` (ควรได้ 400)
func TestLoginStaff_MissingHospital(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	defer teardownTestDB()

	loginPayload := map[string]string{
		"username": "admin",
		"password": "password",
	}
	payloadBytes, _ := json.Marshal(loginPayload)

	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/staff/login", controllers.LoginStaff)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ทดสอบเข้าสู่ระบบโดยใช้ username ที่ไม่มีอยู่จริง (ควรได้ 401)
func TestLoginStaff_NonExistentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	defer teardownTestDB()

	loginPayload := map[string]string{
		"username": "unknown_user",
		"password": "password",
		"hospital": "Hospital",
	}
	payloadBytes, _ := json.Marshal(loginPayload)

	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/staff/login", controllers.LoginStaff)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ทดสอบเข้าสู่ระบบโดยไม่ได้ส่งข้อมูลใดๆ (ควรได้ 400)
func TestLoginStaff_EmptyRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	defer teardownTestDB()

	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/staff/login", controllers.LoginStaff)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
