package tests

import (
	"HIS-api/config"
	"HIS-api/controllers"
	"HIS-api/middlewares"
	"HIS-api/models"
	"HIS-api/routes"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func ptr(s string) *string {
	return &s
}

// ตั้งค่าข้อมูลก่อนการทดสอบ (สร้างทั้ง Staff และ Patient)
func setupTestDB() {
	config.ConnectDB()
	config.DB.Exec("DELETE FROM staffs")
	config.DB.Exec("DELETE FROM patients")

	// พิ่ม Staff ทดสอบ
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	staff := models.Staff{
		Username: "admin",
		Password: string(hashedPassword),
		Hospital: "Hospital",
	}
	config.DB.Create(&staff)

	// เพิ่ม Staff ที่อยู่คนละโรงพยาบาล
	staffOther := models.Staff{
		Username: "admin_other",
		Password: string(hashedPassword),
		Hospital: "OtherHospital",
	}
	config.DB.Create(&staffOther)

	// เพิ่มข้อมูลคนไข้ตัวอย่าง
	dob, _ := time.Parse("2006-01-02", "1990-05-12") //แปลงวันเกิดให้ถูกต้อง
	patient := models.Patient{
		FirstNameTH: "สมชาย",
		LastNameTH:  "สุขดี",
		DateOfBirth: dob,
		PatientHN:   ptr("HN001"),
		NationalID:  ptr("1234567890123"),
		PassportID:  ptr("A12345678"),
		PhoneNumber: "0812345678",
		Email:       "somchai@example.com",
		Gender:      "M",
		Hospital:    "Hospital",
	}
	config.DB.Create(&patient)
}

// ฟังก์ชันสร้าง Token จริงจากการล็อกอิน
func getValidToken(username string, hospital string) string {
	loginPayload := map[string]string{
		"username": username,
		"password": "password",
		"hospital": hospital,
	}
	payloadBytes, _ := json.Marshal(loginPayload)

	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/staff/login", controllers.LoginStaff)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil || w.Code != http.StatusOK {
		log.Fatalf("Login failed: %v", err)
	}
	return response["token"].(string)
}

// ตั้งค่า Router
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	config.ConnectDB()
	r := gin.Default()
	r.Use(middlewares.AuthMiddleware())
	routes.PatientRoutes(r)
	return r
}

// ฟังก์ชันช่วยทำ Request
func performRequest(r http.Handler, method, path string, body []byte, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ทดสอบค้นหาผู้ป่วยด้วยพารามิเตอร์ต่างๆ
func TestSearchPatient_ByVariousFields(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()
	token := getValidToken("admin", "Hospital")

	tests := []struct {
		query  string
		status int
	}{
		{"?national_id=1234567890123", http.StatusOK},
		{"?passport_id=A12345678", http.StatusOK},
		{"?first_name=สมชาย", http.StatusOK},
		{"?last_name=สุขดี", http.StatusOK},
		{"?email=somchai@example.com", http.StatusOK},
		{"?phone_number=0812345678", http.StatusOK},
	}

	for _, test := range tests {
		w := performRequest(router, "GET", "/patient/search"+test.query, nil, token)
		require.Equal(t, test.status, w.Code)
	}
}

// ทดสอบกรณีที่ไม่พบผู้ป่วย (`404 Not Found`)
func TestSearchPatient_NotFound(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()
	token := getValidToken("admin", "Hospital")

	w := performRequest(router, "GET", "/patient/search?national_id=9999999999999", nil, token)

	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	require.Empty(t, response["patients"])
}

// ทดสอบกรณี Staff อยู่คนละโรงพยาบาล (`401 Unauthorized`)
func TestSearchPatient_WrongHospital(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()
	token := getValidToken("admin_other", "OtherHospital") // Staff จากโรงพยาบาลอื่น

	w := performRequest(router, "GET", "/patient/search?national_id=1234567890123", nil, token)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

// ทดสอบกรณี `date_of_birth` ผิดรูปแบบ (`400 Bad Request`)
func TestSearchPatient_InvalidDateFormat(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()
	token := getValidToken("admin", "Hospital")

	w := performRequest(router, "GET", "/patient/search?date_of_birth=invalid_date", nil, token)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

// ทดสอบกรณี Token หมดอายุ (`401 Unauthorized`)
func TestSearchPatient_ExpiredToken(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()

	// Token ที่หมดอายุ
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDAwMDAwMDB9.-GXBv2pYYm3qE2dNjU2x2ti_mSR2e1G03FVY97OfwYs"

	w := performRequest(router, "GET", "/patient/search?national_id=1234567890123", nil, expiredToken)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

// ทดสอบเมื่อไม่มี Token (`401 Unauthorized`)
func TestSearchPatient_Unauthorized(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()

	w := performRequest(router, "GET", "/patient/search?national_id=9876543210987", nil, "")

	require.Equal(t, http.StatusUnauthorized, w.Code)
}
