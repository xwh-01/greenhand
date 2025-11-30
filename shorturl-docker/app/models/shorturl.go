package models

import (
	"time"
)

type ShortURL struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	ShortCode  string    `gorm:"size:10;uniqueIndex;not null" json:"short_code"`
	LongURL    string    `gorm:"type:text;not null" json:"long_url"`
	CreatedAt  time.Time `json:"created_at"`
	ClickCount int       `gorm:"default:0" json:"click_count"`
}

func (ShortURL) TableName() string {
	return "short_urls"
}

func NewShortURL(shortCode, longURL string) *ShortURL {
	return &ShortURL{
		ShortCode:  shortCode,
		LongURL:    longURL,
		CreatedAt:  time.Now(),
		ClickCount: 0,
	}
}

func (s *ShortURL) IncrementClick() {
	s.ClickCount++
}
