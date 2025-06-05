package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

var ValidRoles = []string{RoleUser, RoleAdmin}

type User struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Email             string    `json:"email" gorm:"unique;not null"`
	EncryptedPassword string    `json:"-" gorm:"not null"`
	Role              string    `json:"role" gorm:"default:user;check:role IN ('user','admin')"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func IsValidRole(role string) bool {
	for _, validRole := range ValidRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

func (u *User) SetRole(role string) error {
	if !IsValidRole(role) {
		return errors.New("недопустимая роль")
	}
	u.Role = role
	return nil
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) HasRole(role string) bool {
	return u.Role == role
}

func (u *User) CanModifyRoles() bool {
	return u.Role == RoleAdmin
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.EncryptedPassword = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password))
	return err == nil
}

func FindUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(db *gorm.DB, email, password string) (*User, error) {
	user := &User{
		Email: email,
		Role:  RoleUser,
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func CreateUserWithRole(db *gorm.DB, email, password, role string) (*User, error) {
	if !IsValidRole(role) {
		return nil, errors.New("недопустимая роль")
	}

	user := &User{
		Email: email,
		Role:  role,
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUserRole(db *gorm.DB, userID uint, newRole string, modifierUser *User) error {
	if !modifierUser.CanModifyRoles() {
		return errors.New("недостаточно прав для изменения ролей")
	}

	if !IsValidRole(newRole) {
		return errors.New("недопустимая роль")
	}

	return db.Model(&User{}).Where("id = ?", userID).Update("role", newRole).Error
}
