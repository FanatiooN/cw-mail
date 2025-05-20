package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	JWT struct {
		Secret     string
		Expiration time.Duration
	}
	RabbitMQ struct {
		Host     string
		Port     string
		User     string
		Password string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
		DB       string
	}
	Server struct {
		Port string
	}
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{}

	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnv("DB_PORT", "5432")
	config.Database.User = getEnv("DB_USER", "postgres")
	config.Database.Password = getEnv("DB_PASSWORD", "postgres")
	config.Database.Name = getEnv("DB_NAME", "mailservice")
	config.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	config.JWT.Secret = getEnv("JWT_SECRET", "your_jwt_secret_key")
	jwtExpiration := getEnv("JWT_EXPIRATION", "24h")
	duration, err := time.ParseDuration(jwtExpiration)
	if err != nil {
		return nil, fmt.Errorf("неверный формат JWT_EXPIRATION: %w", err)
	}
	config.JWT.Expiration = duration

	config.RabbitMQ.Host = getEnv("RABBITMQ_HOST", "localhost")
	config.RabbitMQ.Port = getEnv("RABBITMQ_PORT", "5672")
	config.RabbitMQ.User = getEnv("RABBITMQ_USER", "guest")
	config.RabbitMQ.Password = getEnv("RABBITMQ_PASSWORD", "guest")

	config.Redis.Host = getEnv("REDIS_HOST", "localhost")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.DB = getEnv("REDIS_DB", "0")

	config.Server.Port = getEnv("SERVER_PORT", "8080")

	return config, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

func (c *Config) GetRabbitMQURI() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		c.RabbitMQ.User,
		c.RabbitMQ.Password,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
	)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
