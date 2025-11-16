package storage

import (
	"my_go_learning/models"
	"sync"
)

type MemoryStorage struct {
	data    map[string]*models.ShortURL
	counter int
	mutex   sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data:    make(map[string]*models.ShortURL),
		counter: 1,
	}
}

func (s *MemoryStorage) FindByShortCode(shortCode string) (*models.ShortURL, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if shortURL, exists := s.data[shortCode]; exists {
		return shortURL, nil
	}
	return nil, nil
}

func (s *MemoryStorage) Save(shortURL *models.ShortURL) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	shortURL.ID = s.counter
	s.counter++
	s.data[shortURL.ShortCode] = shortURL
	return nil
}

func (s *MemoryStorage) GetAll() []*models.ShortURL {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make([]*models.ShortURL, 0, len(s.data))
	for _, shortURL := range s.data {
		result = append(result, shortURL)
	}
	return result
}
