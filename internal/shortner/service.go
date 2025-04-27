package shortner

import (
	"fmt"
	"linkshortener/internal/database"
	"linkshortener/internal/utils"
	"log/slog"
	"sync"
)

type Service struct {
	db    database.Database
	cache database.Cache
	mx    sync.Mutex
	id    int
}

func NewService(db database.Database, cache database.Cache) *Service {
	Service := &Service{
		db:    db,
		cache: cache,
	}
	Service.id = Service.getLatestId()
	return Service
}

func (s *Service) getLatestId() int {
	rows, err := s.db.Query("SELECT id from links ORDER BY id DESC LIMIT 1")
	if err != nil {
		slog.Error("Error in getting starting ID", "Error", err)
		return 0
	}

	var id int
	rows.Next()
	err = rows.Scan(&id)

	if err != nil {
		slog.Error("Error in getting starting ID", "Error", err)
		return 0
	}
	slog.Info("Starting with", "id", id)

	return id
}

func (s *Service) createShortUrl(originalUrl string) (string, error) {

	if !utils.IsValidURL(originalUrl) {
		return "", fmt.Errorf("Invalid URL! check HTTP/HTTPs")
	}

	shortCode, id := s.getShortCodeFromId()

	query := "INSERT INTO links (id, short_code, original_url) VALUES ($1, $2, $3)"
	slog.Debug("Running", "Query", query)

	stmt, _ := s.db.Prepare(query)
	_, err := stmt.Exec(id, shortCode, originalUrl)

	if err != nil {
		return "", err
	}

	s.cache.Upsert(shortCode, originalUrl)

	return shortCode, nil

}

func (s *Service) getShortCodeFromId() (string, int) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.id++
	shortCode := utils.EncodeBase62(int64(s.id))
	return shortCode, s.id
}

func (s *Service) getOriginalUrl(shortCode string) (string, error) {

	originalUrl, cacheErr := s.cache.Get(shortCode)

	if cacheErr != nil {
		slog.Debug("Cache", "Error", cacheErr)
	} else if len(originalUrl) != 0 {
		return originalUrl, nil
	}

	query := "SELECT original_url from links where short_code='" + shortCode + "'"
	slog.Debug("Running", "Query", query)

	rows, err := s.db.Query(query)
	if err != nil {
		return "", err
	}

	for rows.Next() {
		rows.Scan(&originalUrl)
	}
	return originalUrl, nil
}
