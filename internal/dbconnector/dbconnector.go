package dbconnector

import (
	"fmt"
	"log"

	"github.com/cepmap/party-finder-bot/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnector struct {
	db *gorm.DB
}

// NewDBConnector создает новый коннектор с существующим подключением
func NewDBConnector(db *gorm.DB) *DBConnector {
	return &DBConnector{db: db}
}

// NewDBConnectorFromConfig создает новый коннектор с параметрами из конфигурации
func NewDBConnectorFromConfig(cfg *config.Config) (*DBConnector, error) {
	// Формируем строку подключения к PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Moscow",
		cfg.DatabaseURL, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error), // Логируем только ошибки
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL database")

	// Автоматическая миграция
	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return &DBConnector{db: db}, nil
}

// EnsureAdminsExists проверяет наличие админов в БД и создает их при необходимости
func (dc *DBConnector) EnsureAdminsExists(adminTelegramIDs []int64) error {
	for _, adminID := range adminTelegramIDs {
		var user User
		result := dc.db.Where("telegram_id = ?", adminID).First(&user)

		if result.Error != nil {
			// Создаем админа
			admin := User{
				TelegramID: adminID,
				IsAdmin:    true,
			}

			if err := dc.db.Create(&admin).Error; err != nil {
				return fmt.Errorf("failed to create admin with ID %d: %v", adminID, err)
			}
			log.Printf("Created admin with Telegram ID: %d", adminID)
		} else if !user.IsAdmin {
			// Обновляем существующего пользователя до админа
			if err := dc.db.Model(&user).Update("is_admin", true).Error; err != nil {
				return fmt.Errorf("failed to update user to admin with ID %d: %v", adminID, err)
			}
			log.Printf("Updated user to admin with Telegram ID: %d", adminID)
		}
	}

	return nil
}

// GetDB возвращает подключение к базе данных
func (dc *DBConnector) GetDB() *gorm.DB {
	return dc.db
}

// IsUserAdmin проверяет, является ли пользователь админом
func (dc *DBConnector) IsUserAdmin(telegramID int64, adminTelegramIDs []int64) bool {
	// Проверяем в списке админов из конфига
	for _, adminTelegramID := range adminTelegramIDs {
		if adminTelegramID == telegramID {
			return true
		}
	}

	// Проверяем в базе данных
	var user User
	if err := dc.db.Where("telegram_id = ?", telegramID).First(&user).Error; err == nil {
		return user.IsAdmin
	}

	return false
}

// GetOrCreateUser получает пользователя из БД или создает нового
func (dc *DBConnector) GetOrCreateUser(telegramID int64, adminTelegramIDs []int64) (*User, error) {
	var user User
	result := dc.db.Where("telegram_id = ?", telegramID).First(&user)

	if result.Error != nil {
		// Создаем нового пользователя
		isAdmin := dc.IsUserAdmin(telegramID, adminTelegramIDs)
		user = User{
			TelegramID: telegramID,
			IsAdmin:    isAdmin,
		}

		if err := dc.db.Create(&user).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

type User struct {
	gorm.Model
	TelegramID int64 `gorm:"uniqueIndex"`
	IsAdmin    bool  `gorm:"default:false"`
}

type Games struct {
	gorm.Model
	GameName      string
	AssignedUsers []User
}

type Events struct {
	gorm.Model
	EventName     string
	AssignedUsers []User
}
