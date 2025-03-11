package storage

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"time"
	"url-shortener/internal/handlers/models"
)

// DatabaseStorage реализует хранение в PostgreSQL
type DatabaseStorage struct {
	db             *sql.DB
	logger         *zap.SugaredLogger
	defaultTimeout time.Duration
}

// NewDatabaseStorage создаёт новое хранилище
func NewDatabaseStorage(db *sql.DB, logger *zap.SugaredLogger) *DatabaseStorage {
	ds := &DatabaseStorage{
		db:             db,
		logger:         logger,
		defaultTimeout: 3 * time.Second, // Таймаут по умолчанию
	}

	// Создаём таблицу, если её нет
	if err := ds.init(); err != nil {
		logger.Fatalw("Failed to initialize database storage", "error", err)
	}

	return ds
}

// withTimeout создаёт контекст с таймаутом по умолчанию
func (ds *DatabaseStorage) withTimeout(seconds ...int) (context.Context, context.CancelFunc) {
	timeout := ds.defaultTimeout
	if len(seconds) > 0 {
		timeout = time.Duration(seconds[0]) * time.Second
	}
	return context.WithTimeout(context.Background(), timeout)
}

// init инициализирует БД
func (ds *DatabaseStorage) init() error {
	ctx, cancel := ds.withTimeout(5)
	defer cancel()

	query := `
	CREATE TABLE IF NOT EXISTS short_urls (
		id TEXT PRIMARY KEY,
		full_url TEXT NOT NULL,
		created BIGINT NOT NULL
	);`
	_, err := ds.db.ExecContext(ctx, query)
	return err
}

func (ds *DatabaseStorage) Ping(timeoutSeconds ...int) error {
	ctx, cancel := ds.withTimeout(timeoutSeconds...)
	defer cancel()

	if err := ds.db.PingContext(ctx); err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}
	return nil
}

// SetShortURL сохраняет сокращённый URL в БД
func (ds *DatabaseStorage) SetShortURL(shortURL *models.ShortURL) error {
	ctx, cancel := ds.withTimeout()
	defer cancel()

	query := `INSERT INTO short_urls (id, full_url, created) VALUES ($1, $2, $3) 
			  ON CONFLICT (id) DO UPDATE SET full_url = EXCLUDED.full_url, created = EXCLUDED.created;`
	_, err := ds.db.ExecContext(ctx, query, shortURL.ID, shortURL.FullURL, shortURL.Created)
	if err != nil {
		ds.logger.Errorf("Failed to save short URL: %v", err)
		return err
	}
	return nil
}

// GetShortURL возвращает сокращённый URL по ID
func (ds *DatabaseStorage) GetShortURL(id string) (*models.ShortURL, error) {
	ctx, cancel := ds.withTimeout()
	defer cancel()

	query := `SELECT id, full_url, created FROM short_urls WHERE id = $1;`
	row := ds.db.QueryRowContext(ctx, query, id)

	var shortURL models.ShortURL
	if err := row.Scan(&shortURL.ID, &shortURL.FullURL, &shortURL.Created); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("short URL not found")
		}
		ds.logger.Errorf("Failed to retrieve short URL: %v", err)
		return nil, err
	}

	return &shortURL, nil
}
