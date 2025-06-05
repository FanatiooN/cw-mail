package models

import (
	"testing"
)

func TestUserPassword(t *testing.T) {
	user := &User{
		Email: "test@example.com",
		Role:  "user",
	}

	err := user.SetPassword("password123")
	if err != nil {
		t.Fatalf("Ошибка установки пароля: %v", err)
	}

	if user.EncryptedPassword == "" {
		t.Error("Пароль не был хеширован")
	}

	if !user.CheckPassword("password123") {
		t.Error("Не удалось проверить корректный пароль")
	}

	if user.CheckPassword("wrong_password") {
		t.Error("Проверка неверного пароля должна вернуть false")
	}
}

func TestIsValidRole(t *testing.T) {
	if !IsValidRole(RoleUser) {
		t.Error("RoleUser должна быть валидной")
	}
	if !IsValidRole(RoleAdmin) {
		t.Error("RoleAdmin должна быть валидной")
	}

	if IsValidRole("invalid") {
		t.Error("'invalid' не должна быть валидной ролью")
	}
	if IsValidRole("") {
		t.Error("Пустая строка не должна быть валидной ролью")
	}
}

func TestUserRoles(t *testing.T) {
	admin := &User{Role: RoleAdmin}
	if !admin.IsAdmin() {
		t.Error("IsAdmin() должен возвращать true для администратора")
	}
	if !admin.CanModifyRoles() {
		t.Error("CanModifyRoles() должен возвращать true для администратора")
	}

	user := &User{Role: RoleUser}
	if user.IsAdmin() {
		t.Error("IsAdmin() должен возвращать false для обычного пользователя")
	}
	if user.CanModifyRoles() {
		t.Error("CanModifyRoles() должен возвращать false для обычного пользователя")
	}
}
