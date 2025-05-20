package models

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserPassword(t *testing.T) {
	user := &User{
		Email: "test@example.com",
		Role:  "user",
	}

	// Тест установки пароля
	err := user.SetPassword("password123")
	if err != nil {
		t.Fatalf("Ошибка установки пароля: %v", err)
	}

	// Проверка что пароль был хеширован
	if user.EncryptedPassword == "" {
		t.Error("Пароль не был хеширован")
	}

	// Проверка верного пароля
	if !user.CheckPassword("password123") {
		t.Error("Не удалось проверить корректный пароль")
	}

	// Проверка неверного пароля
	if user.CheckPassword("wrong_password") {
		t.Error("Проверка неверного пароля должна вернуть false")
	}
}

func TestCreateUser(t *testing.T) {
	// Настройка тестовой БД в памяти
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Ошибка подключения к тестовой БД: %v", err)
	}

	// Миграция таблиц
	err = db.AutoMigrate(&User{})
	if err != nil {
		t.Fatalf("Ошибка миграции: %v", err)
	}

	// Тест создания пользователя
	user, err := CreateUser(db, "new@example.com", "password123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Проверки
	if user.Email != "new@example.com" {
		t.Errorf("Ожидался email 'new@example.com', получен %s", user.Email)
	}

	if user.Role != "user" {
		t.Errorf("Ожидалась роль 'user', получена %s", user.Role)
	}

	// Проверка поиска пользователя
	found, err := FindUserByEmail(db, "new@example.com")
	if err != nil {
		t.Fatalf("Ошибка поиска пользователя: %v", err)
	}

	if found.ID != user.ID {
		t.Errorf("Неверный ID пользователя: ожидался %d, получен %d", user.ID, found.ID)
	}
}
