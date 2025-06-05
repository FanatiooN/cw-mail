package models

import (
	"testing"

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

func TestIsValidRole(t *testing.T) {
	// Валидные роли
	if !IsValidRole(RoleUser) {
		t.Error("RoleUser должна быть валидной")
	}
	if !IsValidRole(RoleAdmin) {
		t.Error("RoleAdmin должна быть валидной")
	}

	// Невалидные роли
	if IsValidRole("invalid") {
		t.Error("'invalid' не должна быть валидной ролью")
	}
	if IsValidRole("") {
		t.Error("Пустая строка не должна быть валидной ролью")
	}
}

func TestUserRoles(t *testing.T) {
	// Тест администратора
	admin := &User{Role: RoleAdmin}
	if !admin.IsAdmin() {
		t.Error("IsAdmin() должен возвращать true для администратора")
	}
	if !admin.CanModifyRoles() {
		t.Error("CanModifyRoles() должен возвращать true для администратора")
	}

	// Тест обычного пользователя
	user := &User{Role: RoleUser}
	if user.IsAdmin() {
		t.Error("IsAdmin() должен возвращать false для обычного пользователя")
	}
	if user.CanModifyRoles() {
		t.Error("CanModifyRoles() должен возвращать false для обычного пользователя")
	}
}
