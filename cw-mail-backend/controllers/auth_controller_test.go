//go:build integration
// +build integration

package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/mail-service/config"
	"github.com/mail-service/models"
)

func setupTestRouter() (*gin.Engine, *gorm.DB) {
	// Настройка тестовой БД
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{})

	// Создание тестового пользователя
	models.CreateUser(db, "test@example.com", "password123")

	// Создание тестовой конфигурации
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
		// Добавьте другие необходимые поля конфигурации
	}

	// Настройка маршрутизатора
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Создание и подключение контроллера с конфигурацией
	authController := NewAuthController(db, cfg)
	r.POST("/auth/login", authController.Login)
	r.POST("/auth/register", authController.Register)

	return r, db
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
