//go:build integration
// +build integration

package test

import (
	"testing"
	"time"
)

type User struct {
	ID                uint
	Email             string
	EncryptedPassword string
	Role              string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type MockUserRepo struct {
	users  map[string]*User
	nextID uint
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		users:  make(map[string]*User),
		nextID: 1,
	}
}

func (m *MockUserRepo) FindByEmail(email string) (*User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (m *MockUserRepo) Create(email, encryptedPassword string) (*User, error) {
	user := &User{
		ID:                m.nextID,
		Email:             email,
		EncryptedPassword: encryptedPassword,
		Role:              "user",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	m.users[email] = user
	m.nextID++
	return user, nil
}

type UserService struct {
	repo *MockUserRepo
}

func NewUserService(repo *MockUserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(email, password string) (*User, error) {
	existingUser, _ := s.repo.FindByEmail(email)
	if existingUser != nil {
		return nil, nil
	}

	encryptedPassword := "hashed_" + password

	return s.repo.Create(email, encryptedPassword)
}

func (s *UserService) FindUserByEmail(email string) (*User, error) {
	return s.repo.FindByEmail(email)
}

func TestUserRegistration(t *testing.T) {
	repo := NewMockUserRepo()
	service := NewUserService(repo)

	email := "test@example.com"
	password := "password123"

	t.Log("Регистрация пользователя...")
	user, err := service.Register(email, password)
	if err != nil {
		t.Fatalf("Ошибка регистрации пользователя: %v", err)
	}

	if user.Email != email {
		t.Errorf("Ожидался email %s, получен %s", email, user.Email)
	}

	if user.Role != "user" {
		t.Errorf("Ожидалась роль 'user', получена %s", user.Role)
	}

	foundUser, err := service.FindUserByEmail(email)
	if err != nil {
		t.Fatalf("Ошибка поиска пользователя: %v", err)
	}

	if foundUser.ID != user.ID {
		t.Errorf("Найден неверный пользователь. Ожидался ID %d, получен %d", user.ID, foundUser.ID)
	}

	// Тест повторной регистрации (должна вернуть nil, nil)
	duplicateUser, err := service.Register(email, password)
	if duplicateUser != nil {
		t.Errorf("Ожидался nil при повторной регистрации, получен пользователь с ID %d", duplicateUser.ID)
	}

	t.Log("Тест успешно завершен!")
}
