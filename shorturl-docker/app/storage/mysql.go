package storage

import (
	"shortlink-engine-sql/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLStorage struct {
	db *gorm.DB
}

func NewMySQLStorage(dsn string) (*MySQLStorage, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移创建表
	err = db.AutoMigrate(&models.ShortURL{})
	if err != nil {
		return nil, err
	}

	return &MySQLStorage{db: db}, nil
}

func (s *MySQLStorage) Save(shortURL *models.ShortURL) error {
	return s.db.Create(shortURL).Error
}

func (s *MySQLStorage) FindByShortCode(shortCode string) (*models.ShortURL, error) {
	var shortURL models.ShortURL
	result := s.db.Where("short_code = ?", shortCode).First(&shortURL)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &shortURL, nil
}

func (s *MySQLStorage) GetAll() ([]*models.ShortURL, error) {
	var shortURLs []*models.ShortURL
	result := s.db.Order("created_at DESC").Find(&shortURLs)
	return shortURLs, result.Error
}

func (s *MySQLStorage) IncrementClick(shortCode string) error {
	return s.db.Model(&models.ShortURL{}).
		Where("short_code = ?", shortCode).
		Update("click_count", gorm.Expr("click_count + ?", 1)).Error
}
