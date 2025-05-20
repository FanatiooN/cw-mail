package database

import (
	"fmt"
	"log"

	"github.com/mail-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	log.Println("База данных успешно инициализирована")
	return db, nil
}


func Migrate(db *gorm.DB) error {

	err := db.AutoMigrate(
		&models.User{},
		&models.Message{},
	)
	if err != nil {
		return fmt.Errorf("ошибка миграции базы данных: %w", err)
	}

	log.Println("Миграция базы данных успешно выполнена")
	return nil
}
