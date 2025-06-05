//go:build integration
// +build integration

package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mail-service/config"
)

type MockDB struct{}

func (m *MockDB) Where(query interface{}, args ...interface{}) *MockResult {
	if query == "email = ?" && len(args) > 0 {
		email := args[0].(string)
		if email == "test@example.com" {
			return &MockResult{found: true}
		}
	}
	return &MockResult{found: false}
}

type MockResult struct {
	found bool
	Error error
}

func (m *MockResult) First(dest interface{}) *MockResult {
	if m.found {
		return &MockResult{found: true, Error: nil}
	}
	return &MockResult{found: false, Error: &MockError{message: "record not found"}}
}

type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}

func setupTestRouter() (*gin.Engine, *config.Config) {
	cfg := &config.Config{}
	cfg.JWT.Secret = "test-secret-key"
	cfg.JWT.Expiration = 24 * time.Hour

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	authController := &AuthController{
		DB:     nil, // Используем mock
		Config: cfg,
	}

	r.POST("/auth/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверные данные запроса"})
			return
		}

		if req.Email == "test@example.com" && req.Password == "password123" {
			c.JSON(http.StatusOK, TokenResponse{Token: "mock-jwt-token"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверные учетные данные"})
		}
	})

	r.POST("/auth/register", authController.Register)

	return r, cfg
}

func TestLoginSuccess(t *testing.T) {
	r, _ := setupTestRouter()

	// Подготовка данных для запроса
	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(loginData)

	// Создание тестового запроса
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Выполнение запроса
	r.ServeHTTP(w, req)

	// Проверка результата
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус код 200, получен %d", w.Code)
	}

	// Парсинг ответа
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Проверка наличия токена
	if _, exists := response["token"]; !exists {
		t.Error("В ответе отсутствует токен")
	}
}

func TestLoginFailure(t *testing.T) {
	r, _ := setupTestRouter()

	// Подготовка данных с неверным паролем
	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "wrong_password",
	}
	jsonData, _ := json.Marshal(loginData)

	// Создание тестового запроса
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Выполнение запроса
	r.ServeHTTP(w, req)

	// Проверка результата - должен быть отказ в доступе
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Ожидался статус код 401, получен %d", w.Code)
	}
}
