// internal/config/config.go
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config содержит все конфигурационные настройки приложения
type Config struct {
	DB        DBConfig
	JWT       JWTConfig
	SMTP      SMTPConfig
	CentralCB CentralCBConfig
	App       AppConfig
}

// DBConfig содержит параметры подключения к PostgreSQL
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig содержит настройки JWT-аутентификации
type JWTConfig struct {
	Secret   string
	Lifetime time.Duration
}

// SMTPConfig содержит параметры SMTP-сервера
type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

// CentralCBConfig содержит настройки интеграции с ЦБ РФ
type CentralCBConfig struct {
	WSDLURL      string
	Timeout      time.Duration
	RetryCount   int
	RetryDelay   time.Duration
}

// AppConfig содержит общие настройки приложения
type AppConfig struct {
	Env          string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type EncryptionConfig struct {
    Key    string
}

type HMACConfig struct {
    Secret string
}

// LoadConfig загружает конфигурацию из файла .env и переменных окружения
func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")

	// Чтение конфигурационного файла
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Загрузка значений в структуру Config
	cfg := &Config{
		DB: DBConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		JWT: JWTConfig{
			Secret:   viper.GetString("JWT_SECRET"),
			Lifetime: viper.GetDuration("JWT_LIFETIME"),
		},
		SMTP: SMTPConfig{
			Host:     viper.GetString("SMTP_HOST"),
			Port:     viper.GetInt("SMTP_PORT"),
			User:     viper.GetString("SMTP_USER"),
			Password: viper.GetString("SMTP_PASSWORD"),
			From:     viper.GetString("SMTP_FROM"),
		},
		CentralCB: CentralCBConfig{
			WSDLURL:      viper.GetString("CENTRAL_CB_WSDL_URL"),
			Timeout:      viper.GetDuration("CENTRAL_CB_TIMEOUT"),
			RetryCount:   viper.GetInt("CENTRAL_CB_RETRY_COUNT"),
			RetryDelay:   viper.GetDuration("CENTRAL_CB_RETRY_DELAY"),
		},
		App: AppConfig{
			Env:          viper.GetString("APP_ENV"),
			HTTPPort:     viper.GetInt("HTTP_PORT"),
			ReadTimeout:  viper.GetDuration("READ_TIMEOUT") * time.Second,
			WriteTimeout: viper.GetDuration("WRITE_TIMEOUT") * time.Second,
		},
	}

	// Валидация обязательных полей
	if cfg.DB.Host == "" {
		return nil, fmt.Errorf("DB_HOST is required")
	}
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}
