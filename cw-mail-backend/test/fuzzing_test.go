package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func mockRegisterHandler(c *gin.Context) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if request.Email == "" {
		c.JSON(400, gin.H{"error": "Email is required"})
		return
	}

	if !strings.Contains(request.Email, "@") {
		c.JSON(400, gin.H{"error": "Invalid email format"})
		return
	}

	if request.Password == "" {
		c.JSON(400, gin.H{"error": "Password is required"})
		return
	}

	if len(request.Password) < 6 {
		c.JSON(400, gin.H{"error": "Password too short"})
		return
	}

	// SQL инъекции
	suspiciousPatterns := []string{"drop", "delete", "update", "insert", "select", "union", "--", "/*", "*/"}
	emailLower := strings.ToLower(request.Email)
	passwordLower := strings.ToLower(request.Password)

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(emailLower, pattern) || strings.Contains(passwordLower, pattern) {
			c.JSON(400, gin.H{"error": "Suspicious input detected"})
			return
		}
	}

	// XSS
	if strings.Contains(request.Email, "<") || strings.Contains(request.Password, "<") ||
		strings.Contains(request.Email, "script") || strings.Contains(request.Password, "script") {
		c.JSON(400, gin.H{"error": "Invalid characters detected"})
		return
	}

	// Имитация успешной регистрации
	c.JSON(201, gin.H{
		"message": "User registered successfully",
		"user_id": 123,
	})
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/auth/register", mockRegisterHandler)
	return router
}

func FuzzUserRegistration(f *testing.F) {
	router := setupTestRouter()

	f.Add("test@example.com", "password123")
	f.Add("", "")
	f.Add("admin@test.com", "admin")
	f.Add("user@domain.com", "12345")
	f.Add("test'; DROP TABLE users; --@test.com", "password")
	f.Add("test@test.com", "<script>alert('xss')</script>")

	f.Fuzz(func(t *testing.T, email, password string) {
		requestBody := map[string]string{
			"email":    email,
			"password": password,
		}

		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == 0 {
			t.Errorf("Server crashed on input: email=%s, password=%s", email, password)
		}

		validCodes := []int{200, 201, 400, 401, 409, 422, 500}
		isValid := false
		for _, code := range validCodes {
			if w.Code == code {
				isValid = true
				break
			}
		}

		if !isValid {
			t.Errorf("Unexpected status code %d for email=%s, password=%s", w.Code, email, password)
		}

		if strings.Contains(strings.ToLower(email), "script") ||
			strings.Contains(strings.ToLower(password), "drop") ||
			strings.Contains(email, "<") || strings.Contains(password, "<") {
			t.Logf("Potentially malicious input detected: email=%s, password=%s", email, password)
		}

		if len(email) > 1000 || len(password) > 1000 {
			start := time.Now()
			router.ServeHTTP(httptest.NewRecorder(), req)
			duration := time.Since(start)

			if duration > time.Second*2 {
				t.Logf("Slow response (%v) for large input: email_len=%d, password_len=%d",
					duration, len(email), len(password))
			}
		}

		if strings.Contains(strings.ToLower(email), "drop") && w.Code != 400 {
			t.Logf("Potential SQL injection not blocked: email=%s", email)
		}

		if strings.Contains(email, "<script") && w.Code != 400 {
			t.Logf("Potential XSS attack not blocked: email=%s", email)
		}
	})
}
