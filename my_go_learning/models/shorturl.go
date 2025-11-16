package models

import "time"

type ShortURL struct {
	ID         int       `json:"id"`
	ShortCode  string    `json:"short_code"`
	LongURL    string    `json:"long_url"`
	CreatedAt  time.Time `json:"created_at"`
	ClickCount int       `json:"click_count"`
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
